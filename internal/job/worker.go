package job

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	repo     Repository
	interval time.Duration
}

func NewWorker(repo Repository, interval time.Duration) *Worker {
	return &Worker{
		repo:     repo,
		interval: interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processNextJob(ctx)
		}
	}
}

func (w *Worker) processNextJob(ctx context.Context) {
	job, err := w.repo.ClaimNextPendingJob(ctx)
	if err != nil {
		if err == ErrNoPendingJobs {
			return
		}
		log.Printf("failed to claim job: %v", err)
		return
	}

	log.Printf("Processing job ID: %s, Type: %s", job.ID, job.Type)

	// To be implemented: actual job processing logic based on job.Type and job.Payload
}
