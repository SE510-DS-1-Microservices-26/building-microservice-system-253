package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

type mockItemCategoryRepository struct {
	categories map[uint]*domain.ItemCategory
	nextID     uint
}

func newMockItemCategoryRepository() *mockItemCategoryRepository {
	return &mockItemCategoryRepository{categories: make(map[uint]*domain.ItemCategory)}
}

func (m *mockItemCategoryRepository) Store(_ context.Context, category *domain.ItemCategory) error {
	m.nextID++
	category.ID = m.nextID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	m.categories[category.ID] = category
	return nil
}

func (m *mockItemCategoryRepository) Find(_ context.Context, id uint) (*domain.ItemCategory, error) {
	c, ok := m.categories[id]
	if !ok {
		return nil, fmt.Errorf("category %d not found", id)
	}
	return c, nil
}

func (m *mockItemCategoryRepository) List(_ context.Context, _ dto.ListFilter) ([]domain.ItemCategory, error) {
	list := make([]domain.ItemCategory, 0, len(m.categories))
	for _, c := range m.categories {
		list = append(list, *c)
	}
	return list, nil
}

func (m *mockItemCategoryRepository) Count(_ context.Context) (uint, error) {
	return uint(len(m.categories)), nil
}

func (m *mockItemCategoryRepository) Update(_ context.Context, category *domain.ItemCategory) error {
	if _, ok := m.categories[category.ID]; !ok {
		return fmt.Errorf("category %d not found", category.ID)
	}
	category.UpdatedAt = time.Now()
	m.categories[category.ID] = category
	return nil
}

func (m *mockItemCategoryRepository) Delete(_ context.Context, id uint) error {
	if _, ok := m.categories[id]; !ok {
		return fmt.Errorf("category %d not found", id)
	}
	delete(m.categories, id)
	return nil
}

func setupItemCategoryService() *ItemCategoryService {
	return NewItemCategoryService(newMockItemCategoryRepository())
}

func TestItemCategoryService_Store(t *testing.T) {
	svc := setupItemCategoryService()

	category := &domain.ItemCategory{Name: "Mains"}
	err := svc.Store(ctx, category)

	require.NoError(t, err)
	assert.NotZero(t, category.ID)
	assert.NotZero(t, category.CreatedAt)
}

func TestItemCategoryService_Find(t *testing.T) {
	svc := setupItemCategoryService()

	category := &domain.ItemCategory{Name: "Drinks"}
	require.NoError(t, svc.Store(ctx, category))

	t.Run("existing category", func(t *testing.T) {
		found, err := svc.Find(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, category.ID, found.ID)
		assert.Equal(t, "Drinks", found.Name)
	})

	t.Run("non existing category", func(t *testing.T) {
		_, err := svc.Find(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestItemCategoryService_Update(t *testing.T) {
	svc := setupItemCategoryService()

	category := &domain.ItemCategory{Name: "Mains"}
	require.NoError(t, svc.Store(ctx, category))

	t.Run("existing category", func(t *testing.T) {
		category.Name = "Desserts"
		err := svc.Update(ctx, category)
		require.NoError(t, err)

		updated, err := svc.Find(ctx, category.ID)
		require.NoError(t, err)
		assert.Equal(t, "Desserts", updated.Name)
	})

	t.Run("non existing category", func(t *testing.T) {
		ghost := &domain.ItemCategory{ID: 9999, Name: "Ghost"}
		err := svc.Update(ctx, ghost)
		assert.Error(t, err)
	})
}

func TestItemCategoryService_Delete(t *testing.T) {
	svc := setupItemCategoryService()

	category := &domain.ItemCategory{Name: "Mains"}
	require.NoError(t, svc.Store(ctx, category))

	t.Run("existing category", func(t *testing.T) {
		err := svc.Delete(ctx, category.ID)
		require.NoError(t, err)

		_, err = svc.Find(ctx, category.ID)
		assert.Error(t, err)
	})

	t.Run("non existing category", func(t *testing.T) {
		err := svc.Delete(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestItemCategoryService_List(t *testing.T) {
	svc := setupItemCategoryService()

	require.NoError(t, svc.Store(ctx, &domain.ItemCategory{Name: "Mains"}))
	require.NoError(t, svc.Store(ctx, &domain.ItemCategory{Name: "Drinks"}))

	categories, err := svc.List(ctx, dto.ListFilter{})
	require.NoError(t, err)
	assert.Len(t, categories, 2)
}

func TestItemCategoryService_Count(t *testing.T) {
	svc := setupItemCategoryService()

	require.NoError(t, svc.Store(ctx, &domain.ItemCategory{Name: "Mains"}))
	require.NoError(t, svc.Store(ctx, &domain.ItemCategory{Name: "Drinks"}))

	count, err := svc.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(2), count)
}
