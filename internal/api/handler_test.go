package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

func TestCreateJob(t *testing.T) {
	repo := job.NewMemoryRepository()
	handler := NewHandler(repo)

	request := httptest.NewRequest(
		http.MethodPost,
		"/jobs",
		strings.NewReader(`{"type":"email","payload":"hello"}`),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.CreateJob(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Errorf(
			"Expected status code %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	var createdJob job.Job

	if err := json.NewDecoder(recorder.Body).Decode(&createdJob); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if createdJob.ID == uuid.Nil {
		t.Errorf("Expected job ID to be set, but it was empty")
	}

	if createdJob.Type != "email" {
		t.Errorf("Unexpected job type: %s", createdJob.Type)
	}

	var payload string

	if err := json.Unmarshal(createdJob.Payload, &payload); err != nil {
		t.Fatalf("Failed to decode job payload: %v", err)
	}

	if payload != "hello" {
		t.Errorf("Unexpected job payload: %s", createdJob.Payload)
	}

	if createdJob.Status != "pending" {
		t.Errorf("Unexpected job status: %s", createdJob.Status)
	}
}
