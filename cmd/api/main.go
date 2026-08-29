package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ericlamnguyen/distributed-job-queue/internal/api"
	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
	"github.com/ericlamnguyen/distributed-job-queue/internal/database"
	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	repo := job.NewPostgresRepository(pool)
	handler := api.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateJob(w, r)
		case http.MethodGet:
			handler.ListJobs(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.GetJob(w, r)
	})

	addr := ":" + strconv.Itoa(cfg.Port)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
