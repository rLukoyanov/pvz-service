package services

import (
	"context"
	"pvz-service/internal/models"
	"pvz-service/internal/pkg/errors"
	"pvz-service/internal/repositories"
	"strings"

	"github.com/sirupsen/logrus"
)

type PVZService struct {
	repos *repositories.Repos
}

func NewPVZService(repos *repositories.Repos) *PVZService {
	return &PVZService{repos: repos}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	pvz.City = strings.ToLower(pvz.City)

	allowedCities := map[string]bool{
		"москва":          true,
		"санкт-петербург": true,
		"казань":          true,
	}

	if _, ok := allowedCities[pvz.City]; !ok {
		return models.PVZ{}, errors.ErrCityNotAllowed
	}

	return s.repos.PvzRepo.CreatePVZ(ctx, pvz)
}

func (s *PVZService) GetPVZByID(ctx context.Context, id string) (models.PVZ, error) {
	return s.repos.PvzRepo.GetPVZByID(ctx, id)
}

func (s *PVZService) DeletePVZ(ctx context.Context, id string) error {
	return s.repos.PvzRepo.DeletePVZ(ctx, id)
}

func (s *PVZService) DeleteLastProduct(ctx context.Context, id string) error {
	pvz, err := s.repos.PvzRepo.GetPVZByID(ctx, id)
	if err != nil {
		return err
	}

	reception, err := s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, pvz.ID)
	if err != nil {
		return err
	}

	if reception == nil {
		return errors.ErrNoReceprionsFound
	}

	logrus.Info(pvz, reception)
	return s.repos.ProductRepo.DeleteLastProduct(ctx, reception.ID)
}

func (s *PVZService) CloseLastReception(ctx context.Context, id string) error {
	reception, err := s.repos.ReceptionRepo.GetActiveReceptionByPVZID(ctx, id)
	if err != nil {
		return err
	}

	if reception == nil {
		return errors.ErrNoReceprionsFound
	}

	return s.repos.ReceptionRepo.CloseReception(ctx, reception.ID)
}
