package job

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	jobs map[uuid.UUID]Job
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		jobs: make(map[uuid.UUID]Job),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID] = job
	return nil
}

func (r *MemoryRepository) Get(ctx context.Context, id uuid.UUID) (Job, error) {
	select {
	case <-ctx.Done():
		return Job{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		return Job{}, errors.New("job not found")
	}

	return job, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Job, error) {
	select {
	case <-ctx.Done():
		return []Job{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *MemoryRepository) ClaimNextPendingJob(ctx context.Context) (Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, job := range r.jobs {
		if job.Status == StatusPending {
			job.Status = StatusProcessing
			job.UpdatedAt = time.Now().UTC()

			r.jobs[id] = job

			return job, nil
		}
	}

	return Job{}, ErrNoPendingJobs
}

func (r *MemoryRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exist := r.jobs[id]
	if !exist {
		return ErrJobNotFound
	}

	job.Status = status
	job.UpdatedAt = time.Now().UTC()

	r.jobs[id] = job

	return nil
}
