package main

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philipp-mlr/al-id-maestro/handler"
	"github.com/philipp-mlr/al-id-maestro/service"
)

func main() {
	log.Println("Initializing database...")
	db, err := service.InitDB("al-id-maestro")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loading configuration...")
	config, err := service.NewConfig(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Calculating allowed ID ranges...")
	allowedList, err := service.NewAllowList(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing repositories in memory...")
	err = service.GetRepositories(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting cron jobs...")
	go service.StartCronJob(30*time.Second, db, config)

	log.Print("Starting server...\n\n")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &handler.CustomValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}
	e.Static("/static", "./public")

	indexHandler := handler.IndexHandler{
		DB: db,
	}
	e.GET("/", indexHandler.HandleIndexShow)

	claimHandler := handler.ClaimHandler{
		DB:          db,
		AllowedList: allowedList,
	}
	e.GET("/claim", claimHandler.HandlePageShow)

	historyHandler := handler.HistoryHandler{
		DB: db,
	}
	e.GET("/history", historyHandler.HandleHistoryShow)
	e.POST("/history", historyHandler.HandlePostQuery)

	duplicatesHandler := handler.DuplicatesHandler{
		DB: db,
	}
	e.GET("/duplicates", duplicatesHandler.HandleDuplicatesShow)

	remoteHandler := handler.RemoteHandler{
		DB: db,
	}
	e.GET("/remote", remoteHandler.HandleRemoteShow)

	usedHandler := handler.UsedHandler{
		DB: db,
	}
	e.GET("/used", usedHandler.HandleUsedShow)

	aboutHandler := handler.AboutHandler{}
	e.GET("/about", aboutHandler.HandleAboutShow)

	e.POST("/claim/query-type", claimHandler.HandleObjectTypeQuery)

	e.POST("/claim/request-id", claimHandler.HandleRequestID)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
