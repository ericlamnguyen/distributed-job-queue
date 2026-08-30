package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                      int
	DatabaseURL               string
	WorkerPollForWorkInterval time.Duration
	NumWorkers                int
}

func Load() Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080 // default port
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://jobqueue:jobqueue@localhost:5432/jobqueue?sslmode=disable" // default database URL for testing
	}

	workerPollForWorkInterval, err := time.ParseDuration(os.Getenv("WORKER_POLL_FOR_WORK_INTERVAL"))
	if err != nil {
		workerPollForWorkInterval = 5 * time.Second // default poll interval
	}

	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		numWorkers = 1 // default number of workers
	}

	return Config{
		Port:                      port,
		DatabaseURL:               databaseURL,
		WorkerPollForWorkInterval: workerPollForWorkInterval,
		NumWorkers:                numWorkers,
	}
}
