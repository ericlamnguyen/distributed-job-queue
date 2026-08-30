package job

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job Job) error
	Get(ctx context.Context, id uuid.UUID) (Job, error)
	List(ctx context.Context) ([]Job, error)
	ClaimNextPendingJob(ctx context.Context) (Job, error)
}
