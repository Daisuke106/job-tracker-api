package domain

import "time"

type Company struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	Industry   string    `json:"industry" db:"industry"`
	WebsiteURL string    `json:"website_url" db:"website_url"`
	Memo       string    `json:"memo" db:"memo"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
