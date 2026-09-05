package job

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestWorker_ProcessNext(t *testing.T) {
	ctx := context.Background()

	repo := NewMemoryRepository()

	now := time.Now().UTC()

	expected := Job{
		ID:        uuid.New(),
		Type:      "email",
		Payload:   json.RawMessage(`{"to":"user@example.com"}`),
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Create(ctx, expected); err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	handler := NewDefaultHandler()
	workerId := 1
	worker := NewWorker(workerId, repo, handler, time.Second)

	worker.processNextJob(ctx)

	actual, err := repo.Get(ctx, expected.ID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if actual.Status != StatusCompleted {
		t.Errorf(
			"Status: expected %s, got %s",
			StatusCompleted,
			actual.Status,
		)
	}
}
