package services

import (
	"context"
	"pvz-service/internal/models"
	"pvz-service/internal/pkg/errors"
	"pvz-service/internal/repositories"
	"time"
)

type ReceptionService struct {
	repos *repositories.Repos
}

func NewReceptionService(repos *repositories.Repos) *ReceptionService {
	return &ReceptionService{repos: repos}
}

func (s *ReceptionService) CreateReception(ctx context.Context, reception models.Reception) error {
	active, err := s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, reception.PvzId)
	if err != nil {
		return err
	}
	if active != nil {
		return errors.ErrInvalidInput
	}

	reception.DateTime = time.Now()
	reception.Status = "in_progress"

	return s.repos.ReceptionRepo.CreateReception(ctx, reception)
}

func (s *ReceptionService) GetActiveReceptionByPVZID(ctx context.Context, pvzID string) (*models.Reception, error) {
	return s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, pvzID)
}
