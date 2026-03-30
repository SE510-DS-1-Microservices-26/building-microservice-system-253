package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOrderRepository struct {
	orders map[uint]*domain.Order
	nextID uint
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{orders: make(map[uint]*domain.Order)}
}

func (m *mockOrderRepository) Store(_ context.Context, order *domain.Order) error {
	m.nextID++
	order.ID = m.nextID
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	snapshot := *order
	m.orders[order.ID] = &snapshot
	return nil
}

func (m *mockOrderRepository) Find(_ context.Context, id uint) (*domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %d not found", id)
	}
	return o, nil
}

func (m *mockOrderRepository) List(_ context.Context, _ dto.OrderFilter) ([]domain.Order, error) {
	list := make([]domain.Order, 0, len(m.orders))
	for _, o := range m.orders {
		list = append(list, *o)
	}
	return list, nil
}

func (m *mockOrderRepository) Count(_ context.Context, _ dto.OrderFilter) (uint, error) {
	return uint(len(m.orders)), nil
}

func (m *mockOrderRepository) Update(_ context.Context, order *domain.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return fmt.Errorf("order %d not found", order.ID)
	}
	order.UpdatedAt = time.Now()
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(_ context.Context, id uint) error {
	if _, ok := m.orders[id]; !ok {
		return fmt.Errorf("order %d not found", id)
	}
	delete(m.orders, id)
	return nil
}

func setupOrderService() *OrderService {
	itemRepo := newMockItemRepository()
	_ = itemRepo.Store(ctx, &domain.Item{Name: "Burger", Price: 9.99, Quantity: 10})
	_ = itemRepo.Store(ctx, &domain.Item{Name: "Cola", Price: 2.99, Quantity: 5})
	return NewOrderService(newMockOrderRepository(), itemRepo, nil, nil)
}

func TestOrderService_Store(t *testing.T) {
	svc := setupOrderService()

	t.Run("sets status to pending", func(t *testing.T) {
		order := &domain.Order{
			Items: []domain.OrderItem{
				{ItemID: 1, Quantity: 2},
			},
		}
		err := svc.Store(ctx, order)
		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusPending, order.Status)
		assert.NotZero(t, order.ID)
	})

	t.Run("snapshots item price", func(t *testing.T) {
		order := &domain.Order{
			Items: []domain.OrderItem{
				{ItemID: 1, Quantity: 1},
			},
		}
		err := svc.Store(ctx, order)
		require.NoError(t, err)
		assert.Equal(t, 9.99, order.Items[0].UnitPrice)
	})

	t.Run("item not found", func(t *testing.T) {
		order := &domain.Order{
			Items: []domain.OrderItem{
				{ItemID: 9999, Quantity: 1},
			},
		}
		err := svc.Store(ctx, order)
		require.Error(t, err)
		assert.True(t, errors.Is(err, customErrors.ErrNotFound))
	})

	t.Run("insufficient quantity", func(t *testing.T) {
		order := &domain.Order{
			Items: []domain.OrderItem{
				{ItemID: 2, Quantity: 100},
			},
		}
		err := svc.Store(ctx, order)
		require.Error(t, err)
		assert.True(t, errors.Is(err, customErrors.ErrConflict))
	})
}

func TestOrderService_Find(t *testing.T) {
	svc := setupOrderService()

	order := &domain.Order{
		Items: []domain.OrderItem{
			{ItemID: 1, Quantity: 2},
			{ItemID: 2, Quantity: 1},
		},
	}
	require.NoError(t, svc.Store(ctx, order))

	t.Run("existing order computes total price", func(t *testing.T) {
		found, err := svc.Find(ctx, order.ID)
		require.NoError(t, err)
		assert.Equal(t, order.ID, found.ID)
		// 2 * 9.99 + 1 * 2.99 = 22.97
		assert.InDelta(t, 22.97, found.TotalPrice, 0.001)
	})

	t.Run("non existing order", func(t *testing.T) {
		_, err := svc.Find(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestOrderService_Update(t *testing.T) {
	svc := setupOrderService()

	order := &domain.Order{
		Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}},
	}
	require.NoError(t, svc.Store(ctx, order))

	t.Run("existing order", func(t *testing.T) {
		order.Status = domain.OrderStatusConfirmed
		err := svc.Update(ctx, order)
		require.NoError(t, err)

		updated, err := svc.Find(ctx, order.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.OrderStatusConfirmed, updated.Status)
	})

	t.Run("non existing order", func(t *testing.T) {
		ghost := &domain.Order{ID: 9999, Status: domain.OrderStatusConfirmed}
		err := svc.Update(ctx, ghost)
		assert.Error(t, err)
	})
}

func TestOrderService_Delete(t *testing.T) {
	svc := setupOrderService()

	order := &domain.Order{
		Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}},
	}
	require.NoError(t, svc.Store(ctx, order))

	t.Run("existing order", func(t *testing.T) {
		err := svc.Delete(ctx, order.ID)
		require.NoError(t, err)

		_, err = svc.Find(ctx, order.ID)
		assert.Error(t, err)
	})

	t.Run("non existing order", func(t *testing.T) {
		err := svc.Delete(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestOrderService_List(t *testing.T) {
	svc := setupOrderService()

	require.NoError(t, svc.Store(ctx, &domain.Order{Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}}}))
	require.NoError(t, svc.Store(ctx, &domain.Order{Items: []domain.OrderItem{{ItemID: 2, Quantity: 1}}}))

	orders, err := svc.List(ctx, dto.OrderFilter{})
	require.NoError(t, err)
	assert.Len(t, orders, 2)
}
