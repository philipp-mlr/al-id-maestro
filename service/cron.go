package service

import (
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/model"
)

var (
	mutex     sync.Mutex
	isRunning bool
)

func StartCronJob(interval time.Duration, db *sqlx.DB, config *model.Config) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		// Check if the function is already running
		if !isRunning {
			go syncJob(db, config) // Only run if not already running
		} else {
			log.Println("Previous sync instance still running, skipping...")
		}
	}
}

func syncJob(db *sqlx.DB, config *model.Config) {
	// Lock the mutex to ensure only one instance runs
	mutex.Lock()
	isRunning = true
	defer func() {
		isRunning = false
		mutex.Unlock()
	}()
	log.Print("Starting sync job...\n\n")

	err := Sync(db, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Expiring claimed objects...\n")
	err = UpdateClaimed(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Sync job finished\n\n")
}
