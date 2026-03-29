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

type mockItemRepository struct {
	items  map[uint]*domain.Item
	nextID uint
}

func newMockItemRepository() *mockItemRepository {
	return &mockItemRepository{items: make(map[uint]*domain.Item)}
}

func (m *mockItemRepository) Store(_ context.Context, item *domain.Item) error {
	m.nextID++
	item.ID = m.nextID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.items[item.ID] = item
	return nil
}

func (m *mockItemRepository) Find(_ context.Context, id uint) (*domain.Item, error) {
	i, ok := m.items[id]
	if !ok {
		return nil, fmt.Errorf("item %d not found", id)
	}
	return i, nil
}

func (m *mockItemRepository) FindByIDs(_ context.Context, ids []uint) ([]domain.Item, error) {
	result := make([]domain.Item, 0, len(ids))
	for _, id := range ids {
		if i, ok := m.items[id]; ok {
			result = append(result, *i)
		}
	}
	return result, nil
}

func (m *mockItemRepository) List(_ context.Context, _ dto.ItemFilter) ([]domain.Item, error) {
	list := make([]domain.Item, 0, len(m.items))
	for _, i := range m.items {
		list = append(list, *i)
	}
	return list, nil
}

func (m *mockItemRepository) Count(_ context.Context, _ dto.ItemFilter) (uint, error) {
	return uint(len(m.items)), nil
}

func (m *mockItemRepository) Update(_ context.Context, item *domain.Item) error {
	if _, ok := m.items[item.ID]; !ok {
		return fmt.Errorf("item %d not found", item.ID)
	}
	item.UpdatedAt = time.Now()
	m.items[item.ID] = item
	return nil
}

func (m *mockItemRepository) Delete(_ context.Context, id uint) error {
	if _, ok := m.items[id]; !ok {
		return fmt.Errorf("item %d not found", id)
	}
	delete(m.items, id)
	return nil
}

func setupItemService() *ItemService {
	return NewItemService(newMockItemRepository())
}

func TestItemService_Store(t *testing.T) {
	svc := setupItemService()

	item := &domain.Item{Name: "Burger", Price: 9.99, Quantity: 10}
	err := svc.Store(ctx, item)

	require.NoError(t, err)
	assert.NotZero(t, item.ID)
	assert.NotZero(t, item.CreatedAt)
}

func TestItemService_Find(t *testing.T) {
	svc := setupItemService()

	item := &domain.Item{Name: "Burger", Price: 9.99, Quantity: 10}
	require.NoError(t, svc.Store(ctx, item))

	t.Run("existing item", func(t *testing.T) {
		found, err := svc.Find(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, item.ID, found.ID)
		assert.Equal(t, "Burger", found.Name)
	})

	t.Run("non existing item", func(t *testing.T) {
		_, err := svc.Find(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestItemService_Update(t *testing.T) {
	svc := setupItemService()

	item := &domain.Item{Name: "Burger", Price: 9.99, Quantity: 10}
	require.NoError(t, svc.Store(ctx, item))

	t.Run("existing item", func(t *testing.T) {
		item.Price = 12.99
		err := svc.Update(ctx, item)
		require.NoError(t, err)

		updated, err := svc.Find(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, 12.99, updated.Price)
	})

	t.Run("non existing item", func(t *testing.T) {
		ghost := &domain.Item{ID: 9999, Name: "Ghost"}
		err := svc.Update(ctx, ghost)
		assert.Error(t, err)
	})
}

func TestItemService_Delete(t *testing.T) {
	svc := setupItemService()

	item := &domain.Item{Name: "Burger", Price: 9.99, Quantity: 10}
	require.NoError(t, svc.Store(ctx, item))

	t.Run("existing item", func(t *testing.T) {
		err := svc.Delete(ctx, item.ID)
		require.NoError(t, err)

		_, err = svc.Find(ctx, item.ID)
		assert.Error(t, err)
	})

	t.Run("non existing item", func(t *testing.T) {
		err := svc.Delete(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestItemService_List(t *testing.T) {
	svc := setupItemService()

	require.NoError(t, svc.Store(ctx, &domain.Item{Name: "Burger", Price: 9.99}))
	require.NoError(t, svc.Store(ctx, &domain.Item{Name: "Cola", Price: 2.99}))

	items, err := svc.List(ctx, dto.ItemFilter{})
	require.NoError(t, err)
	assert.Len(t, items, 2)
}

func TestItemService_Count(t *testing.T) {
	svc := setupItemService()

	require.NoError(t, svc.Store(ctx, &domain.Item{Name: "Burger", Price: 9.99}))
	require.NoError(t, svc.Store(ctx, &domain.Item{Name: "Cola", Price: 2.99}))

	count, err := svc.Count(ctx, dto.ItemFilter{})
	require.NoError(t, err)
	assert.Equal(t, uint(2), count)
}
