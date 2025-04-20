package services

import (
	"context"
	"pvz-service/internal/models"
)

// PVZServiceInterface определяет методы для работы с ПВЗ
type PVZServiceInterface interface {
	GetAll(ctx context.Context, page, limit, from, to string) ([]models.FullPVZ, error)
	CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error)
	GetPVZByID(ctx context.Context, id string) (models.PVZ, error)
	DeletePVZ(ctx context.Context, id string) error
	DeleteLastProduct(ctx context.Context, id string) error
	CloseLastReception(ctx context.Context, id string) error
}

// UserServiceInterface определяет методы для работы с пользователями
type UserServiceInterface interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

// ProductServiceInterface определяет методы для работы с товарами
type ProductServiceInterface interface {
	AddProduct(ctx context.Context, product models.Product, pvzID string) error
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

// ReceptionServiceInterface определяет методы для работы с приемками
type ReceptionServiceInterface interface {
	CreateReception(ctx context.Context, reception models.Reception) error
	GetActiveReceptionByPVZID(ctx context.Context, pvzID string) (*models.Reception, error)
}
