package services

import (
	"context"
	"pvz-service/internal/models"
	"pvz-service/internal/pkg/errors"
	"pvz-service/internal/repositories"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type PVZService struct {
	repos *repositories.Repos
}

func NewPVZService(repos *repositories.Repos) *PVZService {
	return &PVZService{repos: repos}
}

func (s *PVZService) GetAll(ctx context.Context, pageStr, limitStr, fromStr, toStr string) ([]models.FullPVZ, error) {

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return nil, errors.ErrInvalidInput
		}
	}

	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return nil, errors.ErrInvalidInput
		}
	}

	// TODO - добавить получение приемок и их товаров
	PVZs, err := s.repos.PvzRepo.GetAll(ctx, limit, offset, from, to)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}

	var wg sync.WaitGroup

	for _, pvz := range PVZs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for key, reception := range pvz.Receptions {
				prodcuts, err := s.repos.ProductRepo.GetByReceptionID(ctx, reception.ID)
				if err != nil {
					return
				}

				reception.Products = prodcuts
				pvz.Receptions[key] = reception
			}
		}()

	}

	wg.Wait()

	return PVZs, err
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
