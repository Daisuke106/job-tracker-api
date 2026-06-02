package domain

import "time"

type JobStatus string

const (
	JobStatusApplied   JobStatus = "applied"
	JobStatusScreening JobStatus = "screening"
	JobStatusInterview JobStatus = "interview"
	JobStatusOffer     JobStatus = "offer"
	JobStatusRejected  JobStatus = "rejected"
	JobStatusWithdrawn JobStatus = "withdrawn"
)

type Job struct {
	ID          int64      `json:"id" db:"id"`
	UserID      int64      `json:"user_id" db:"user_id"`
	CompanyID   int64      `json:"company_id" db:"company_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Status      JobStatus  `json:"status" db:"status"`
	AppliedAt   *time.Time `json:"applied_at" db:"applied_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
