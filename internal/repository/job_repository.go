package repository

import (
	"context"
	"database/sql"
	"errors"

	"job-tracker-api/internal/domain"

	"github.com/jmoiron/sqlx"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	FindAll(ctx context.Context, userID int64) ([]domain.Job, error)
	FindByID(ctx context.Context, id, userID int64) (*domain.Job, error)
	UpdateStatus(ctx context.Context, id, userID int64, status domain.JobStatus) error
}

type jobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `
		INSERT INTO jobs (user_id, company_id, title, description, status, applied_at)
		VALUES (:user_id, :company_id, :title, :description, :status, :applied_at)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, job)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
	}
	return nil
}

func (r *jobRepository) FindAll(ctx context.Context, userID int64) ([]domain.Job, error) {
	var jobs []domain.Job
	err := r.db.SelectContext(ctx, &jobs,
		"SELECT * FROM jobs WHERE user_id = $1 ORDER BY created_at DESC", userID)
	return jobs, err
}

func (r *jobRepository) FindByID(ctx context.Context, id, userID int64) (*domain.Job, error) {
	var job domain.Job
	err := r.db.GetContext(ctx, &job,
		"SELECT * FROM jobs WHERE id = $1 AND user_id = $2", id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &job, err
}

func (r *jobRepository) UpdateStatus(ctx context.Context, id, userID int64, status domain.JobStatus) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE jobs SET status = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3",
		status, id, userID)
	return err
}
