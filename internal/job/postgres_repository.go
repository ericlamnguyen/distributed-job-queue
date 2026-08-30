package job

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

// Ensure PostgresRepository implements the Repository interface
var _ Repository = (*PostgresRepository)(nil)

var ErrNoPendingJobs = errors.New("no pending jobs available")

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, job Job) error {
	_, err := r.pool.Exec(
		ctx,
		`
		INSERT INTO jobs (
			id,
			type,
			payload,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.CreatedAt,
		job.UpdatedAt,
	)

	return err
}

func (r *PostgresRepository) Get(ctx context.Context, id uuid.UUID) (Job, error) {
	row := r.pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			type,
			payload,
			status,
			created_at,
			updated_at
		FROM jobs
		WHERE id = $1
		`,
		id,
	)

	var job Job

	err := row.Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return Job{}, err
	}

	return job, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Job, error) {
	rows, err := r.pool.Query(
		ctx,
		`
		SELECT
			id,
			type,
			payload,
			status,
			created_at,
			updated_at
		FROM jobs
		ORDER BY created_at
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job

		if err := rows.Scan(
			&job.ID,
			&job.Type,
			&job.Payload,
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
		); err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *PostgresRepository) ClaimNextPendingJob(ctx context.Context) (Job, error) {
	var job Job

	err := r.pool.QueryRow(ctx, `
		UPDATE jobs
		SET
			status = 'processing',
			updated_at = NOW()
		WHERE id = (
			SELECT id
			FROM jobs
			WHERE status = 'pending'
			ORDER BY created_at
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, type, payload, status, created_at, updated_at
	`,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Job{}, ErrNoPendingJobs
	}

	if err != nil {
		return Job{}, err
	}

	return job, nil
}
