package repository

import (
	"context"
	"database/sql"
	"errors"

	"job-tracker-api/internal/domain"

	"github.com/jmoiron/sqlx"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	FindAll(ctx context.Context, userID int64) ([]domain.Company, error)
	FindByID(ctx context.Context, id, userID int64) (*domain.Company, error)
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id, userID int64) error
}

type companyRepository struct {
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (user_id, name, industry, website_url, memo)
		VALUES (:user_id, :name, :industry, :website_url, :memo)
		RETURNING id, created_at, updated_at`
	rows, err := r.db.NamedQueryContext(ctx, query, company)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&company.ID, &company.CreatedAt, &company.UpdatedAt)
	}
	return nil
}

func (r *companyRepository) FindAll(ctx context.Context, userID int64) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.SelectContext(ctx, &companies,
		"SELECT * FROM companies WHERE user_id = $1 ORDER BY created_at DESC", userID)
	return companies, err
}

func (r *companyRepository) FindByID(ctx context.Context, id, userID int64) (*domain.Company, error) {
	var company domain.Company
	err := r.db.GetContext(ctx, &company,
		"SELECT * FROM companies WHERE id = $1 AND user_id = $2", id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &company, err
}

func (r *companyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = :name, industry = :industry, website_url = :website_url, memo = :memo, updated_at = NOW()
		WHERE id = :id AND user_id = :user_id`
	_, err := r.db.NamedExecContext(ctx, query, company)
	return err
}

func (r *companyRepository) Delete(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM companies WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
