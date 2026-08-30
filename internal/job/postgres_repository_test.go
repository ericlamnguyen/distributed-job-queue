package job

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
	"github.com/ericlamnguyen/distributed-job-queue/internal/database"
	"github.com/google/uuid"
)

func TestPostgresRepository_CreateAndGet(t *testing.T) {
	ctx := context.Background()

	cfg := config.Load()
	databaseURL := cfg.DatabaseURL

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	repo := NewPostgresRepository(pool)

	defer func() {
		_, err := pool.Exec(ctx, "DELETE FROM jobs")
		if err != nil {
			t.Errorf("failed to clean up jobs: %v", err)
		}
		pool.Close()
	}()

	now := time.Now().UTC().Truncate(time.Microsecond)

	expected := Job{
		ID:        uuid.New(),
		Type:      "email",
		Payload:   json.RawMessage(`{"to":"user@example.com","subject":"Welcome"}`),
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Create(ctx, expected); err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	actual, err := repo.Get(ctx, expected.ID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if actual.ID != expected.ID {
		t.Errorf("ID: expected %v, got %v", expected.ID, actual.ID)
	}

	if actual.Type != expected.Type {
		t.Errorf("Type: expected %s, got %s", expected.Type, actual.Type)
	}

	if actual.Status != expected.Status {
		t.Errorf("Status: expected %s, got %s", expected.Status, actual.Status)
	}

	var expectedPayload any
	var actualPayload any

	if err := json.Unmarshal([]byte(expected.Payload), &expectedPayload); err != nil {
		t.Fatalf("invalid expected payload: %v", err)
	}

	if err := json.Unmarshal([]byte(actual.Payload), &actualPayload); err != nil {
		t.Fatalf("invalid actual payload: %v", err)
	}

	expectedJSON, _ := json.Marshal(expectedPayload)
	actualJSON, _ := json.Marshal(actualPayload)

	if string(expectedJSON) != string(actualJSON) {
		t.Errorf(
			"Payload: expected %s, got %s",
			expected.Payload,
			actual.Payload,
		)
	}

	if !actual.CreatedAt.Equal(expected.CreatedAt) {
		t.Errorf(
			"CreatedAt: expected %v, got %v",
			expected.CreatedAt,
			actual.CreatedAt,
		)
	}

	if !actual.UpdatedAt.Equal(expected.UpdatedAt) {
		t.Errorf(
			"UpdatedAt: expected %v, got %v",
			expected.UpdatedAt,
			actual.UpdatedAt,
		)
	}
}

func TestPostgresRepository_List(t *testing.T) {
	ctx := context.Background()

	cfg := config.Load()
	databaseURL := cfg.DatabaseURL

	pool, err := database.NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}

	repo := NewPostgresRepository(pool)

	defer func() {
		_, err := pool.Exec(ctx, "DELETE FROM jobs")
		if err != nil {
			t.Errorf("failed to clean up jobs: %v", err)
		}
		pool.Close()
	}()

	now := time.Now().UTC().Truncate(time.Microsecond)

	expectedJobs := []Job{
		{
			ID:        uuid.New(),
			Type:      "email",
			Payload:   json.RawMessage(`{"to":"user1@example.com"}`),
			Status:    StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uuid.New(),
			Type:      "report",
			Payload:   json.RawMessage(`{"format":"pdf"}`),
			Status:    StatusPending,
			CreatedAt: now.Add(time.Second),
			UpdatedAt: now.Add(time.Second),
		},
	}

	for _, job := range expectedJobs {
		if err := repo.Create(ctx, job); err != nil {
			t.Fatalf("failed to create job: %v", err)
		}
	}

	actualJobs, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list jobs: %v", err)
	}

	for _, expected := range expectedJobs {
		found := false

		for _, actual := range actualJobs {
			if actual.ID == expected.ID {
				found = true

				if actual.Type != expected.Type {
					t.Errorf(
						"Type: expected %s, got %s",
						expected.Type,
						actual.Type,
					)
				}

				if !actual.CreatedAt.Equal(expected.CreatedAt) {
					t.Errorf(
						"CreatedAt: expected %v, got %v",
						expected.CreatedAt,
						actual.CreatedAt,
					)
				}

				if actual.Status != expected.Status {
					t.Errorf(
						"Status: expected %s, got %s",
						expected.Status,
						actual.Status,
					)
				}

				break
			}
		}

		if !found {
			t.Errorf("job %v was not found", expected.ID)
		}
	}
}

func TestPostgresRepository_ClaimNext(t *testing.T) {
	ctx := context.Background()

	cfg := config.Load()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	repo := NewPostgresRepository(pool)

	defer func() {
		_, err := pool.Exec(ctx, "DELETE FROM jobs")
		if err != nil {
			t.Errorf("failed to clean up jobs: %v", err)
		}
	}()

	now := time.Now().UTC().Truncate(time.Microsecond)

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

	actual, err := repo.ClaimNextPendingJob(ctx)
	if err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}

	if actual.ID != expected.ID {
		t.Errorf("ID: expected %v, got %v", expected.ID, actual.ID)
	}

	if actual.Status != StatusProcessing {
		t.Errorf(
			"Status: expected %s, got %s",
			StatusProcessing,
			actual.Status,
		)
	}

	// Verify that the job's status is updated in the database
	stored, err := repo.Get(ctx, expected.ID)
	if err != nil {
		t.Fatalf("failed to get job: %v", err)
	}

	if stored.Status != StatusProcessing {
		t.Errorf(
			"Stored Status: expected %s, got %s",
			StatusProcessing,
			stored.Status,
		)
	}
}
