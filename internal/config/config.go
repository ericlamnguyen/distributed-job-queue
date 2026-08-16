package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
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

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}
}
