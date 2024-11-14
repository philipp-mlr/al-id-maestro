package cron

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/config"
	"github.com/philipp-mlr/al-id-maestro/internal/claim"
	synchronize "github.com/philipp-mlr/al-id-maestro/internal/sync"
)

var (
	mutex     sync.Mutex
	isRunning atomic.Bool
)

func StartCronJob(interval time.Duration, db *sqlx.DB, config *config.Config) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		// Check if the function is already running
		if !isRunning.Load() {
			go syncJob(db, config) // Only run if not already running
		} else {
			log.Println("Previous sync instance still running, skipping...")
		}
	}
}

func syncJob(db *sqlx.DB, config *config.Config) {
	// Lock the mutex to ensure only one instance runs
	mutex.Lock()
	isRunning.Store(true)
	defer func() {
		isRunning.Store(false)
		mutex.Unlock()
	}()
	log.Print("Starting sync job...\n\n")

	err := synchronize.Sync(db, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Expiring claimed objects...\n")
	err = claim.UpdateClaimed(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Sync job finished\n\n")
}
