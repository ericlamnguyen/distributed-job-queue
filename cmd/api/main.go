package main

import (
	"log"
	"net/http"

	"github.com/ericlamnguyen/distributed-job-queue/internal/api"
	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

func main() {
	repo := job.NewMemoryRepository()
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

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
