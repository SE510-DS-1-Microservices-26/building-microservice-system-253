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

type mockItemCategoryService struct {
	categories map[uint]*domain.ItemCategory
	nextID     uint
}

func newMockItemCategoryService() *mockItemCategoryService {
	return &mockItemCategoryService{categories: make(map[uint]*domain.ItemCategory)}
}

func (m *mockItemCategoryService) Store(_ context.Context, c *domain.ItemCategory) error {
	m.nextID++
	c.ID = m.nextID
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.categories[c.ID] = c
	return nil
}

func (m *mockItemCategoryService) Find(_ context.Context, id uint) (*domain.ItemCategory, error) {
	c, ok := m.categories[id]
	if !ok {
		return nil, fmt.Errorf("category %d: %w", id, customErrors.ErrNotFound)
	}
	return c, nil
}

func (m *mockItemCategoryService) List(_ context.Context, _ dto.ListFilter) ([]domain.ItemCategory, error) {
	list := make([]domain.ItemCategory, 0, len(m.categories))
	for _, c := range m.categories {
		list = append(list, *c)
	}
	return list, nil
}

func (m *mockItemCategoryService) Count(_ context.Context) (uint, error) {
	return uint(len(m.categories)), nil
}

func (m *mockItemCategoryService) Update(_ context.Context, c *domain.ItemCategory) error {
	if _, ok := m.categories[c.ID]; !ok {
		return fmt.Errorf("category %d: %w", c.ID, customErrors.ErrNotFound)
	}
	m.categories[c.ID] = c
	return nil
}

func (m *mockItemCategoryService) Delete(_ context.Context, id uint) error {
	if _, ok := m.categories[id]; !ok {
		return fmt.Errorf("category %d: %w", id, customErrors.ErrNotFound)
	}
	delete(m.categories, id)
	return nil
}

func setupItemCategoryApp() (*fiber.App, *mockItemCategoryService) {
	svc := newMockItemCategoryService()
	app := newTestApp()
	routes.ItemCategoryRoutes(app, NewItemCategoryHandlers(svc, newTestPagination()))
	return app, svc
}

func TestItemCategoryHandlers_Create(t *testing.T) {
	app, _ := setupItemCategoryApp()

	t.Run("valid request returns 201", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/item-categories", strings.NewReader(`{"name":"Mains"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing name returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/item-categories", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestItemCategoryHandlers_Show(t *testing.T) {
	app, svc := setupItemCategoryApp()
	require.NoError(t, svc.Store(context.Background(), &domain.ItemCategory{Name: "Mains"}))

	t.Run("existing category returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/item-categories/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing category returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/item-categories/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/item-categories/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestItemCategoryHandlers_List(t *testing.T) {
	app, svc := setupItemCategoryApp()
	require.NoError(t, svc.Store(context.Background(), &domain.ItemCategory{Name: "Mains"}))
	require.NoError(t, svc.Store(context.Background(), &domain.ItemCategory{Name: "Drinks"}))

	req := httptest.NewRequest(http.MethodGet, "/item-categories", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestItemCategoryHandlers_Update(t *testing.T) {
	app, svc := setupItemCategoryApp()
	require.NoError(t, svc.Store(context.Background(), &domain.ItemCategory{Name: "Mains"}))

	t.Run("existing category returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/item-categories/1", strings.NewReader(`{"name":"Updated"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing category returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/item-categories/9999", strings.NewReader(`{"name":"Updated"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestItemCategoryHandlers_Delete(t *testing.T) {
	app, svc := setupItemCategoryApp()
	require.NoError(t, svc.Store(context.Background(), &domain.ItemCategory{Name: "Mains"}))

	t.Run("existing category returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/item-categories/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing category returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/item-categories/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
