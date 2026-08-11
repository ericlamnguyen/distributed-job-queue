package api

import (
	"encoding/json"
	"net/http"
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
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	newJob := job.Job{
		ID:        uuid.New().String(),
		Type:      request.Type,
		Payload:   request.Payload,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = h.repo.Create(newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newJob)
}
