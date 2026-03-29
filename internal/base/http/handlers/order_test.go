package handlers

import (
	"cafeteria-delivery/internal/base/core/domain"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/internal/base/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockOrderService struct {
	orders map[uint]*domain.Order
	nextID uint
}

func newMockOrderService() *mockOrderService {
	return &mockOrderService{orders: make(map[uint]*domain.Order)}
}

func (m *mockOrderService) Store(_ context.Context, order *domain.Order) error {
	m.nextID++
	order.ID = m.nextID
	order.Status = domain.OrderStatusPending
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderService) Find(_ context.Context, id uint) (*domain.Order, error) {
	o, ok := m.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %d: %w", id, customErrors.ErrNotFound)
	}
	return o, nil
}

func (m *mockOrderService) List(_ context.Context, _ dto.OrderFilter) ([]domain.Order, error) {
	list := make([]domain.Order, 0, len(m.orders))
	for _, o := range m.orders {
		list = append(list, *o)
	}
	return list, nil
}

func (m *mockOrderService) Count(_ context.Context, _ dto.OrderFilter) (uint, error) {
	return uint(len(m.orders)), nil
}

func (m *mockOrderService) Update(_ context.Context, order *domain.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return fmt.Errorf("order %d: %w", order.ID, customErrors.ErrNotFound)
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderService) Delete(_ context.Context, id uint) error {
	if _, ok := m.orders[id]; !ok {
		return fmt.Errorf("order %d: %w", id, customErrors.ErrNotFound)
	}
	delete(m.orders, id)
	return nil
}

func setupOrderApp() (*fiber.App, *mockOrderService) {
	svc := newMockOrderService()
	app := newTestApp()
	routes.OrderRoutes(app, NewOrderHandlers(svc, newTestPagination()))
	return app, svc
}

func TestOrderHandlers_Create(t *testing.T) {
	app, _ := setupOrderApp()

	t.Run("valid request returns 201", func(t *testing.T) {
		body := `{"user_id":1,"items":[{"item_id":1,"quantity":2}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing items returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"user_id":1}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("item with zero quantity returns 422", func(t *testing.T) {
		body := `{"user_id":1,"items":[{"item_id":1,"quantity":0}]}`
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestOrderHandlers_Show(t *testing.T) {
	app, svc := setupOrderApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Order{
		Items: []domain.OrderItem{{ItemID: 1, Quantity: 2, UnitPrice: 9.99}},
	}))

	t.Run("existing order returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing order returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestOrderHandlers_List(t *testing.T) {
	app, svc := setupOrderApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Order{Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}}}))
	require.NoError(t, svc.Store(context.Background(), &domain.Order{Items: []domain.OrderItem{{ItemID: 2, Quantity: 3}}}))

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOrderHandlers_Update(t *testing.T) {
	app, svc := setupOrderApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Order{
		Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}},
	}))

	t.Run("existing order returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/orders/1", strings.NewReader(`{"status":2}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing order returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/orders/9999", strings.NewReader(`{"status":2}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid status returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/orders/1", strings.NewReader(`{"status":99}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestOrderHandlers_Delete(t *testing.T) {
	app, svc := setupOrderApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Order{
		Items: []domain.OrderItem{{ItemID: 1, Quantity: 1}},
	}))

	t.Run("existing order returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing order returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/orders/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
