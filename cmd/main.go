package main

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philipp-mlr/al-id-maestro/internal/config"
	"github.com/philipp-mlr/al-id-maestro/internal/cron"
	"github.com/philipp-mlr/al-id-maestro/internal/database"
	"github.com/philipp-mlr/al-id-maestro/internal/git"
	"github.com/philipp-mlr/al-id-maestro/internal/handler"
	"github.com/philipp-mlr/al-id-maestro/internal/licensedList"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
)

func main() {
	db, allowedList, config := initSetup()

	initHttpServer(db, allowedList, config)
}

func initSetup() (*sqlx.DB, *model.LicensedObjectList, map[string]string) {
	log.Println("Initializing database...")
	db, err := database.InitDB("al-id-maestro")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loading configuration...")
	config, err := config.NewConfig(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Calculating allowed ID ranges...")
	allowedList, err := licensedList.NewLicensedObjectList(config)
	if err != nil {
		log.Fatal(err)
	}

	err = git.InitRepos(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting cron jobs...")
	go cron.StartCronJob(30*time.Second, db, config)

	log.Print("Starting server...\n\n")

	repoInformation := make(map[string]string)

	for _, remoteConfig := range config.RemoteConfiguration {
		repoInformation[remoteConfig.RepositoryName] = remoteConfig.RepositoryURL
	}

	return db, allowedList, repoInformation
}

func initHttpServer(db *sqlx.DB, allowedList *model.LicensedObjectList, repoInformation map[string]string) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &handler.CustomValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}
	e.Static("/static", "./website/public")

	indexHandler := handler.IndexHandler{
		DB: db,
	}
	e.GET("/", indexHandler.HandleIndexShow)
	e.GET("/chart", indexHandler.HandleChartShow)

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
		DB:              db,
		RepoInformation: repoInformation,
	}
	e.GET("/duplicates", duplicatesHandler.HandleDuplicatesShow)

	remoteHandler := handler.RemoteHandler{
		DB: db,
	}
	e.GET("/remote", remoteHandler.HandleRemoteShow)

	usedHandler := handler.UsedHandler{
		DB:              db,
		RepoInformation: repoInformation,
	}
	e.GET("/used", usedHandler.HandleUsedShow)

	aboutHandler := handler.AboutHandler{}
	e.GET("/about", aboutHandler.HandleAboutShow)

	e.POST("/claim/query-type", claimHandler.HandleObjectTypeQuery)

	e.POST("/claim/request-id", claimHandler.HandleNewObjectClaim)

	e.POST("/api/request-id", claimHandler.HandleNewObjectClaimAPI)

	// Start server
	e.Logger.Fatal(e.Start(":5000"))
}
