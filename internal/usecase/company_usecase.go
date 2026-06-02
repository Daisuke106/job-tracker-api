package usecase

import (
	"context"
	"errors"

	"job-tracker-api/internal/domain"
	"job-tracker-api/internal/repository"
)

var ErrCompanyNotFound = errors.New("company not found")

type CreateCompanyRequest struct {
	Name       string `json:"name" binding:"required"`
	Industry   string `json:"industry"`
	WebsiteURL string `json:"website_url"`
	Memo       string `json:"memo"`
}

type UpdateCompanyRequest struct {
	Name       string `json:"name" binding:"required"`
	Industry   string `json:"industry"`
	WebsiteURL string `json:"website_url"`
	Memo       string `json:"memo"`
}

type CompanyUsecase interface {
	Create(ctx context.Context, userID int64, req CreateCompanyRequest) (*domain.Company, error)
	GetAll(ctx context.Context, userID int64) ([]domain.Company, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Company, error)
	Update(ctx context.Context, id, userID int64, req UpdateCompanyRequest) (*domain.Company, error)
	Delete(ctx context.Context, id, userID int64) error
}

type companyUsecase struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyUsecase(companyRepo repository.CompanyRepository) CompanyUsecase {
	return &companyUsecase{companyRepo: companyRepo}
}

func (u *companyUsecase) Create(ctx context.Context, userID int64, req CreateCompanyRequest) (*domain.Company, error) {
	company := &domain.Company{
		UserID:     userID,
		Name:       req.Name,
		Industry:   req.Industry,
		WebsiteURL: req.WebsiteURL,
		Memo:       req.Memo,
	}
	if err := u.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (u *companyUsecase) GetAll(ctx context.Context, userID int64) ([]domain.Company, error) {
	return u.companyRepo.FindAll(ctx, userID)
}

func (u *companyUsecase) GetByID(ctx context.Context, id, userID int64) (*domain.Company, error) {
	company, err := u.companyRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}
	return company, nil
}

func (u *companyUsecase) Update(ctx context.Context, id, userID int64, req UpdateCompanyRequest) (*domain.Company, error) {
	company, err := u.companyRepo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, ErrCompanyNotFound
	}

	company.Name = req.Name
	company.Industry = req.Industry
	company.WebsiteURL = req.WebsiteURL
	company.Memo = req.Memo

	if err := u.companyRepo.Update(ctx, company); err != nil {
		return nil, err
	}
	return company, nil
}

func (u *companyUsecase) Delete(ctx context.Context, id, userID int64) error {
	company, err := u.companyRepo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}
	if company == nil {
		return ErrCompanyNotFound
	}
	return u.companyRepo.Delete(ctx, id, userID)
}
