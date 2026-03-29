package services

import (
	"cafeteria-delivery/internal/base/core/ports"
	"context"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
)

type ItemCategoryService struct {
	repo ports.ItemCategoryRepository
}

var _ ports.ItemCategoryService = (*ItemCategoryService)(nil)

func NewItemCategoryService(repository ports.ItemCategoryRepository) *ItemCategoryService {
	return &ItemCategoryService{repo: repository}
}

// Store - store new item category
func (s *ItemCategoryService) Store(ctx context.Context, category *domain.ItemCategory) error {
	return s.repo.Store(ctx, category)
}

// Find - find item category
func (s *ItemCategoryService) Find(ctx context.Context, id uint) (*domain.ItemCategory, error) {
	return s.repo.Find(ctx, id)
}

// List - list item categories
func (s *ItemCategoryService) List(ctx context.Context, filter dto.ListFilter) ([]domain.ItemCategory, error) {
	return s.repo.List(ctx, filter)
}

// Count - count item categories
func (s *ItemCategoryService) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}

// Update - update item category
func (s *ItemCategoryService) Update(ctx context.Context, category *domain.ItemCategory) error {
	return s.repo.Update(ctx, category)
}

// Delete - delete item category
func (s *ItemCategoryService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
