package services

import (
	"cafeteria-delivery/internal/users/core/domain"
	"cafeteria-delivery/internal/users/core/ports"
	"context"

	"cafeteria-delivery/internal/users/dto"
)

type UserService struct {
	repo ports.UserRepository
}

var _ ports.UserService = (*UserService)(nil)

func NewUserService(repository ports.UserRepository) *UserService {
	return &UserService{repo: repository}
}

// Store - store new user
func (s *UserService) Store(ctx context.Context, user *domain.User) error {
	return s.repo.Store(ctx, user)
}

// Find - find user info
func (s *UserService) Find(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.Find(ctx, id)
}

// List - list users
func (s *UserService) List(ctx context.Context, filter dto.ListFilter) ([]domain.User, error) {
	return s.repo.List(ctx, filter)
}

// Count - count users
func (s *UserService) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}

// Update - update user
func (s *UserService) Update(ctx context.Context, user *domain.User) error {
	return s.repo.Update(ctx, user)
}

// Delete - delete user
func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
