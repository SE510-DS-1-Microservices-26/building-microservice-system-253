package services

import (
	"cafeteria-delivery/internal/base/core/ports"
	"context"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
)

type ItemService struct {
	repo ports.ItemRepository
}

var _ ports.ItemService = (*ItemService)(nil)

func NewItemService(repository ports.ItemRepository) *ItemService {
	return &ItemService{repo: repository}
}

// Store - store new menu item
func (s *ItemService) Store(ctx context.Context, item *domain.Item) error {
	return s.repo.Store(ctx, item)
}

// Find - find menu item
func (s *ItemService) Find(ctx context.Context, id uint) (*domain.Item, error) {
	return s.repo.Find(ctx, id)
}

// List - list menu items
func (s *ItemService) List(ctx context.Context, filter dto.ItemFilter) ([]domain.Item, error) {
	return s.repo.List(ctx, filter)
}

// Count - count menu items
func (s *ItemService) Count(ctx context.Context, filter dto.ItemFilter) (uint, error) {
	return s.repo.Count(ctx, filter)
}

// Update - update menu item
func (s *ItemService) Update(ctx context.Context, item *domain.Item) error {
	return s.repo.Update(ctx, item)
}

// Delete - delete menu item
func (s *ItemService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
