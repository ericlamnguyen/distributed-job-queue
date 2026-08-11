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

	http.HandleFunc("/jobs", handler.CreateJob)

	log.Println("Starting API server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
