package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
	"github.com/ericlamnguyen/distributed-job-queue/internal/database"
	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

func main() {
	// Create a context that is canceled on SIGINT or SIGTERM
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

	var wg sync.WaitGroup

	log.Printf("Starting %d worker with poll interval %s", cfg.NumWorkers, cfg.WorkerPollForWorkInterval)
	for i := 0; i < cfg.NumWorkers; i++ {
		workerID := i + 1
		worker := job.NewWorker(workerID, repo, handler, cfg.WorkerPollForWorkInterval)
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			worker.Start(ctx)
		}(workerID)
	}

	// Wait until SIGINT or SIGTERM is received
	<-ctx.Done()
	log.Println("Shutting down workers...")

	// Wait for all workers to finish
	wg.Wait()
	log.Println("All workers stopped")
}
