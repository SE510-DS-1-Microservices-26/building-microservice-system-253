package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/dto"
	customErrors "cafeteria-delivery/internal/errors"
	"cafeteria-delivery/internal/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockItemService struct {
	items  map[uint]*domain.Item
	nextID uint
}

func newMockItemService() *mockItemService {
	return &mockItemService{items: make(map[uint]*domain.Item)}
}

func (m *mockItemService) Store(_ context.Context, item *domain.Item) error {
	m.nextID++
	item.ID = m.nextID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.items[item.ID] = item
	return nil
}

func (m *mockItemService) Find(_ context.Context, id uint) (*domain.Item, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, fmt.Errorf("item %d: %w", id, customErrors.ErrNotFound)
	}
	return item, nil
}

func (m *mockItemService) List(_ context.Context, _ dto.ItemFilter) ([]domain.Item, error) {
	list := make([]domain.Item, 0, len(m.items))
	for _, item := range m.items {
		list = append(list, *item)
	}
	return list, nil
}

func (m *mockItemService) Count(_ context.Context, _ dto.ItemFilter) (uint, error) {
	return uint(len(m.items)), nil
}

func (m *mockItemService) Update(_ context.Context, item *domain.Item) error {
	if _, ok := m.items[item.ID]; !ok {
		return fmt.Errorf("item %d: %w", item.ID, customErrors.ErrNotFound)
	}
	m.items[item.ID] = item
	return nil
}

func (m *mockItemService) Delete(_ context.Context, id uint) error {
	if _, ok := m.items[id]; !ok {
		return fmt.Errorf("item %d: %w", id, customErrors.ErrNotFound)
	}
	delete(m.items, id)
	return nil
}

func setupItemApp() (*fiber.App, *mockItemService) {
	svc := newMockItemService()
	app := newTestApp()
	routes.ItemRoutes(app, NewItemHandlers(svc, newTestPagination()))
	return app, svc
}

func TestItemHandlers_Create(t *testing.T) {
	app, _ := setupItemApp()

	t.Run("valid request returns 201", func(t *testing.T) {
		body := `{"name":"Burger","price":9.99,"quantity":50}`
		req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing required fields returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"name":"Burger"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestItemHandlers_Show(t *testing.T) {
	app, svc := setupItemApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Item{Name: "Burger", Price: 9.99, Quantity: 50}))

	t.Run("existing item returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing item returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestItemHandlers_List(t *testing.T) {
	app, svc := setupItemApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Item{Name: "Burger", Price: 9.99}))
	require.NoError(t, svc.Store(context.Background(), &domain.Item{Name: "Cola", Price: 2.99}))

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestItemHandlers_Update(t *testing.T) {
	app, svc := setupItemApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Item{Name: "Burger", Price: 9.99, Quantity: 50}))

	body := `{"name":"Burger","price":12.99,"quantity":50}`

	t.Run("existing item returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/items/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing item returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/items/9999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestItemHandlers_Delete(t *testing.T) {
	app, svc := setupItemApp()
	require.NoError(t, svc.Store(context.Background(), &domain.Item{Name: "Burger", Price: 9.99}))

	t.Run("existing item returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing item returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/items/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
