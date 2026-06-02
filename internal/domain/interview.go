package domain

import "time"

type Interview struct {
	ID          int64     `json:"id" db:"id"`
	JobID       int64     `json:"job_id" db:"job_id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	ScheduledAt time.Time `json:"scheduled_at" db:"scheduled_at"`
	Location    string    `json:"location" db:"location"`
	Memo        string    `json:"memo" db:"memo"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type InterviewNote struct {
	ID          int64     `json:"id" db:"id"`
	InterviewID int64     `json:"interview_id" db:"interview_id"`
	Content     string    `json:"content" db:"content"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
