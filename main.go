package main

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/philipp-mlr/al-id-maestro/handler"
	"github.com/philipp-mlr/al-id-maestro/model"
	"github.com/philipp-mlr/al-id-maestro/service"
	"gorm.io/gorm"
)

func main() {
	config := model.NewConfig()
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := model.NewDB("database.db", &gorm.Config{})
	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Save(config.Repositories)
	if err != nil {
		log.Fatal(err)
	}

	objects := new(model.Objects)

	for _, repo := range config.Repositories {
		log.Println("Scanning repository ", repo.URL)
		err = service.ScanRepository(&repo, db, objects)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("Total objects found:", objects.Len())
	}

	// log.Println("Reading licensed objects from CSV file")
	// lic := service.ReadLicensedObjectsCSV("./ids.csv")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &handler.CustomValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}
	e.Static("/static", "./public")

	indexHandler := handler.IndexHanlder{}
	e.GET("/", indexHandler.HandleIndexShow)

	claimHandler := handler.ClaimHandler{
		DB: db,
	}
	e.GET("/claim", claimHandler.HandleClaimShow)
	e.POST("/claim/query-type", claimHandler.HandleClaimTypeQuery)
	e.POST("/claim/request-claim", claimHandler.HandleIDClaim)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
