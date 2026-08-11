package job

import (
	"errors"
	"sync"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		jobs: make(map[string]Job),
	}
}

func (r *MemoryRepository) Create(job Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID] = job
	return nil
}

func (r *MemoryRepository) Get(id string) (Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		return Job{}, errors.New("job not found")
	}

	return job, nil
}

func (r *MemoryRepository) List() ([]Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}
