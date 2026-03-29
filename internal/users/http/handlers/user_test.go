package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cafeteria-delivery/internal/users/core/domain"
	"cafeteria-delivery/internal/users/dto"
	customErrors "cafeteria-delivery/internal/users/errors"
	"cafeteria-delivery/internal/users/http/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserService struct {
	users  map[uint]*domain.User
	nextID uint
}

func newMockUserService() *mockUserService {
	return &mockUserService{users: make(map[uint]*domain.User)}
}

func (m *mockUserService) Store(_ context.Context, user *domain.User) error {
	for _, u := range m.users {
		if u.Email == user.Email {
			return fmt.Errorf("email already exists: %w", customErrors.ErrConflict)
		}
	}
	m.nextID++
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	snapshot := *user
	m.users[user.ID] = &snapshot
	return nil
}

func (m *mockUserService) Find(_ context.Context, id uint) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d: %w", id, customErrors.ErrNotFound)
	}
	return u, nil
}

func (m *mockUserService) List(_ context.Context, _ dto.ListFilter) ([]domain.User, error) {
	list := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserService) Count(_ context.Context) (uint, error) {
	return uint(len(m.users)), nil
}

func (m *mockUserService) Update(_ context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return fmt.Errorf("user %d: %w", user.ID, customErrors.ErrNotFound)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserService) Delete(_ context.Context, id uint) error {
	if _, ok := m.users[id]; !ok {
		return fmt.Errorf("user %d: %w", id, customErrors.ErrNotFound)
	}
	delete(m.users, id)
	return nil
}

func setupUserApp() (*fiber.App, *mockUserService) {
	svc := newMockUserService()
	app := newTestApp()
	routes.UserRoutes(app, NewUserHandlers(svc, newTestPagination()))
	return app, svc
}

func TestUserHandlers_Create(t *testing.T) {
	app, _ := setupUserApp()

	t.Run("valid request returns 201", func(t *testing.T) {
		body := `{"name":"Alice","email":"alice@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing name returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"alice@example.com"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("invalid email returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"Alice","email":"not-an-email"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		body := `{"name":"Alice2","email":"alice@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

func TestUserHandlers_Show(t *testing.T) {
	app, svc := setupUserApp()
	require.NoError(t, svc.Store(context.Background(), &domain.User{Name: "Alice", Email: "alice@example.com"}))

	t.Run("existing user returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing user returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestUserHandlers_List(t *testing.T) {
	app, svc := setupUserApp()
	require.NoError(t, svc.Store(context.Background(), &domain.User{Name: "Alice", Email: "alice@example.com"}))
	require.NoError(t, svc.Store(context.Background(), &domain.User{Name: "Bob", Email: "bob@example.com"}))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandlers_Update(t *testing.T) {
	app, svc := setupUserApp()
	require.NoError(t, svc.Store(context.Background(), &domain.User{Name: "Alice", Email: "alice@example.com"}))

	body := `{"name":"Alice Updated","email":"alice.updated@example.com"}`

	t.Run("existing user returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing user returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/9999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("missing required fields returns 422", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(`{"name":"Alice"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestUserHandlers_Delete(t *testing.T) {
	app, svc := setupUserApp()
	require.NoError(t, svc.Store(context.Background(), &domain.User{Name: "Alice", Email: "alice@example.com"}))

	t.Run("existing user returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing user returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
