package services

import (
	"context"
	"pvz-service/internal/models"
	"pvz-service/internal/pkg/errors"
	"pvz-service/internal/repositories"
	"time"
)

type ProductService struct {
	repos *repositories.Repos
}

func NewProductService(repos *repositories.Repos) *ProductService {
	return &ProductService{repos: repos}
}

func (s *ProductService) AddProduct(ctx context.Context, product models.Product, pvzID string) error {
	reception, err := s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, pvzID)
	if err != nil {
		return err
	}
	if reception == nil {
		return errors.ErrInvalidInput
	}

	product.DateTime = time.Now()
	product.ReceptionId = reception.ID

	allowedTypes := map[string]bool{
		"электроника": true,
		"одежда":      true,
		"обувь":       true,
	}

	if _, ok := allowedTypes[product.Type]; !ok {
		return errors.ErrInvalidInput
	}

	return s.repos.ProductRepo.AddProduct(ctx, product)
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	reception, err := s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, pvzID)
	if err != nil {
		return err
	}
	if reception == nil {
		return errors.ErrInvalidInput
	}

	return s.repos.ProductRepo.DeleteLastProduct(ctx, reception.ID)
}
