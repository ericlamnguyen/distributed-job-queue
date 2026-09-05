package job

import (
	"context"
	"log"
	"time"
)

type Worker struct {
	id       int
	repo     Repository
	handler  Handler
	interval time.Duration
}

func NewWorker(id int, repo Repository, handler Handler, interval time.Duration) *Worker {
	return &Worker{
		id:       id,
		repo:     repo,
		handler:  handler,
		interval: interval,
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Printf("Worker %d started", w.id)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopped", w.id)
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

	err = w.handler.Handle(ctx, job)
	if err != nil {
		log.Printf("failed to process job ID %s: %v", job.ID, err)

		// Update the job status to failed
		failStatusUpdateErr := w.repo.UpdateStatus(ctx, job.ID, StatusFailed)
		if failStatusUpdateErr != nil {
			log.Printf("failed to update job status to failed for job ID %s: %v", job.ID, failStatusUpdateErr)
		}

		return
	}

	// Update the job status to completed
	successStatusUpdateErr := w.repo.UpdateStatus(ctx, job.ID, StatusCompleted)
	if successStatusUpdateErr != nil {
		log.Printf("failed to update job status to completed for job ID %s: %v", job.ID, successStatusUpdateErr)

		return
	}

	log.Printf("Successfully processed job ID: %s", job.ID)
}
