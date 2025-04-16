package services

import (
	"context"
	"pvz-service/internal/models"
	"pvz-service/internal/pkg/errors"
	"pvz-service/internal/repositories"
)

type UserService struct {
	repos *repositories.Repos
}

func NewUserService(repos *repositories.Repos) *UserService {
	return &UserService{repos: repos}
}

func (s *UserService) CreateUser(ctx context.Context, user models.User) error {
	if user.Role != "client" && user.Role != "moderator" {
		return errors.ErrInvalidInput
	}
	return s.repos.AuthRepo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return s.repos.AuthRepo.GetUserByEmail(ctx, email)
}
