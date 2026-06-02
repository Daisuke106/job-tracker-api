package usecase

import (
	"context"
	"errors"
	"time"

	"job-tracker-api/internal/domain"
	"job-tracker-api/internal/repository"
)

var ErrJobNotFound = errors.New("job not found")

type CreateJobRequest struct {
	CompanyID   int64      `json:"company_id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	AppliedAt   *time.Time `json:"applied_at"`
}

type UpdateJobStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type JobUsecase interface {
	Create(ctx context.Context, userID int64, req CreateJobRequest) (*domain.Job, error)
	GetAll(ctx context.Context, userID int64) ([]domain.Job, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Job, error)
	UpdateStatus(ctx context.Context, id, userID int64, req UpdateJobStatusRequest) (*domain.Job, error)
}

type jobUsecase struct {
	jobRepo repository.JobRepository
}

func NewJobUsecase(jobRepo repository.JobRepository) JobUsecase {
	return &jobUsecase{jobRepo: jobRepo}
}

func (u *jobUsecase) Create(ctx context.Context, userID int64, req CreateJobRequest) (*domain.Job, error) {
	status := domain.JobStatus(req.Status)
	if status == "" {
		status = domain.JobStatusApplied
	}

	job := &domain.Job{
		UserID:      userID,
		CompanyID:   req.CompanyID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		AppliedAt:   req.AppliedAt,
	}
	if err := u.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (u *jobUsecase) GetAll(ctx context.Context, userID int64) ([]domain.Job, error) {
	return u.jobRepo.FindAll(ctx, userID)
}

func (u *jobUsecase) GetByID(ctx context.Context, id, userID int64) (*domain.Job, error) {
	job, err := u.jobRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}
	return job, nil
}

func (u *jobUsecase) UpdateStatus(ctx context.Context, id, userID int64, req UpdateJobStatusRequest) (*domain.Job, error) {
	job, err := u.jobRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	status := domain.JobStatus(req.Status)
	if err := u.jobRepo.UpdateStatus(ctx, id, userID, status); err != nil {
		return nil, err
	}
	job.Status = status
	return job, nil
}
