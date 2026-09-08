package integration

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
	"github.com/ericlamnguyen/distributed-job-queue/internal/database"
	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
	"github.com/google/uuid"
)

type testHandler struct {
	mu        sync.Mutex
	processed []uuid.UUID
	err       error
}

func (h *testHandler) Handle(ctx context.Context, j job.Job) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.processed = append(h.processed, j.ID)

	return h.err
}

func (h *testHandler) ProcessedCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()

	return len(h.processed)
}

func (h *testHandler) ProcessedIDs() []uuid.UUID {
	h.mu.Lock()
	defer h.mu.Unlock()

	ids := make([]uuid.UUID, len(h.processed))
	copy(ids, h.processed)

	return ids
}

func setupRepository(t *testing.T) *job.PostgresRepository {
	t.Helper()

	ctx := context.Background()

	cfg := config.Load()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	repo := job.NewPostgresRepository(pool)

	_, err = pool.Exec(ctx, "DELETE FROM jobs")
	if err != nil {
		t.Fatalf("failed to clean jobs table: %v", err)
	}

	return repo
}

func createTestJob(t *testing.T, repo job.Repository) job.Job {
	t.Helper()

	now := time.Now().UTC()

	j := job.Job{
		ID:        uuid.New(),
		Type:      "email",
		Payload:   json.RawMessage(`{"to":"test@example.com"}`),
		Status:    job.StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Create(context.Background(), j); err != nil {
		t.Fatalf("failed to create test job: %v", err)
	}

	return j
}

func TestJobProcessing_Success(t *testing.T) {
	ctx := context.Background()
	repo := setupRepository(t)

	handler := &testHandler{}

	worker := job.NewWorker(
		1,
		repo,
		handler,
		time.Second,
	)

	j := createTestJob(t, repo)

	worker.ProcessNextJob(ctx)

	actual, err := repo.Get(ctx, j.ID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if actual.Status != job.StatusCompleted {
		t.Fatalf(
			"expected status %s, got %s",
			job.StatusCompleted,
			actual.Status,
		)
	}

	if handler.ProcessedCount() != 1 {
		t.Fatalf(
			"expected handler to process 1 job, got %d",
			handler.ProcessedCount(),
		)
	}
}

func TestJobProcessing_Failure(t *testing.T) {
	ctx := context.Background()
	repo := setupRepository(t)

	handler := &testHandler{
		err: errors.New("simulated handler failure"),
	}

	worker := job.NewWorker(
		1,
		repo,
		handler,
		time.Second,
	)

	j := createTestJob(t, repo)

	worker.ProcessNextJob(ctx)

	actual, err := repo.Get(ctx, j.ID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if actual.Status != job.StatusFailed {
		t.Fatalf(
			"expected status %s, got %s",
			job.StatusFailed,
			actual.Status,
		)
	}

	if handler.ProcessedCount() != 1 {
		t.Fatalf(
			"expected handler to process 1 job, got %d",
			handler.ProcessedCount(),
		)
	}
}

func TestJobProcessing_MultipleWorkers(t *testing.T) {
	ctx := context.Background()
	repo := setupRepository(t)

	const (
		numJobs    = 10
		numWorkers = 4
	)

	handler := &testHandler{}

	workers := make([]*job.Worker, numWorkers)

	for i := 0; i < numWorkers; i++ {
		workers[i] = job.NewWorker(
			i+1,
			repo,
			handler,
			time.Second,
		)
	}

	for i := 0; i < numJobs; i++ {
		createTestJob(t, repo)
	}

	var wg sync.WaitGroup

	for _, worker := range workers {
		wg.Add(1)

		go func(w *job.Worker) {
			defer wg.Done()

			for {
				before := handler.ProcessedCount()

				w.ProcessNextJob(ctx)

				after := handler.ProcessedCount()

				if after == before {
					return
				}
			}
		}(worker)
	}

	wg.Wait()

	if handler.ProcessedCount() != numJobs {
		t.Fatalf(
			"expected %d processed jobs, got %d",
			numJobs,
			handler.ProcessedCount(),
		)
	}

	ids := handler.ProcessedIDs()

	uniqueIDs := make(map[uuid.UUID]struct{})

	for _, id := range ids {
		uniqueIDs[id] = struct{}{}
	}

	if len(uniqueIDs) != numJobs {
		t.Fatalf(
			"expected %d unique processed jobs, got %d",
			numJobs,
			len(uniqueIDs),
		)
	}

	jobs, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list jobs: %v", err)
	}

	for _, j := range jobs {
		if j.Status != job.StatusCompleted {
			t.Errorf(
				"job %s: expected status %s, got %s",
				j.ID,
				job.StatusCompleted,
				j.Status,
			)
		}
	}
}
