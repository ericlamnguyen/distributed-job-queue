package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ericlamnguyen/distributed-job-queue/internal/job"
)

type Handler struct {
	repo job.Repository
}

func NewHandler(repo job.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Type) == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	now := time.Now()
	newJob := job.Job{
		ID:        uuid.New(),
		Type:      request.Type,
		Payload:   request.Payload,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = h.repo.Create(r.Context(), newJob)
	if err != nil {
		log.Printf("failed to create job: %v", err)
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, newJob)
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/jobs/")

	if id == "" {
		http.Error(w, "job ID is required", http.StatusBadRequest)
		return
	}

	job, err := h.repo.Get(r.Context(), uuid.MustParse(id))
	if err != nil {
		log.Printf("failed to get job: %v", err)
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.repo.List(r.Context())
	if err != nil {
		log.Printf("failed to list jobs: %v", err)
		http.Error(w, "Failed to list jobs", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, jobs)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
