package services

import (
	"context"
	"pvz-service/internal/models"
)

type PVZServiceInterface interface {
	GetAll(ctx context.Context, page, limit, from, to string) ([]models.FullPVZ, error)
	CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error)
	GetPVZByID(ctx context.Context, id string) (models.PVZ, error)
	DeletePVZ(ctx context.Context, id string) error
	DeleteLastProduct(ctx context.Context, id string) error
	CloseLastReception(ctx context.Context, id string) error
}
