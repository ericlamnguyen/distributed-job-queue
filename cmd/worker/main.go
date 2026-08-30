package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
	"github.com/ericlamnguyen/distributed-job-queue/internal/database"
	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	repo := job.NewPostgresRepository(pool)
	handler := job.NewDefaultHandler()

	worker := job.NewWorker(repo, handler, cfg.WorkerPollForWorkInterval)
	log.Println("worker started")
	worker.Start(ctx)
	log.Println("worker stopped")
}
