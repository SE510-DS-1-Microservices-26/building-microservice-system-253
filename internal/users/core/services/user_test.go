package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cafeteria-delivery/internal/users/core/domain"
	"cafeteria-delivery/internal/users/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

type mockUserRepository struct {
	users  map[uint]*domain.User
	nextID uint
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[uint]*domain.User)}
}

func (m *mockUserRepository) Store(_ context.Context, user *domain.User) error {
	m.nextID++
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	snapshot := *user
	m.users[user.ID] = &snapshot
	return nil
}

func (m *mockUserRepository) Find(_ context.Context, id uint) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

func (m *mockUserRepository) List(_ context.Context, _ dto.ListFilter) ([]domain.User, error) {
	list := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		list = append(list, *u)
	}
	return list, nil
}

func (m *mockUserRepository) Count(_ context.Context) (uint, error) {
	return uint(len(m.users)), nil
}

func (m *mockUserRepository) Update(_ context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return fmt.Errorf("user %d not found", user.ID)
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) Delete(_ context.Context, id uint) error {
	if _, ok := m.users[id]; !ok {
		return fmt.Errorf("user %d not found", id)
	}
	delete(m.users, id)
	return nil
}

func setupUserService() *UserService {
	return NewUserService(newMockUserRepository())
}

func TestUserService_Store(t *testing.T) {
	svc := setupUserService()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	err := svc.Store(ctx, user)

	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
}

func TestUserService_Find(t *testing.T) {
	svc := setupUserService()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	require.NoError(t, svc.Store(ctx, user))

	t.Run("existing user", func(t *testing.T) {
		found, err := svc.Find(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, "Alice", found.Name)
		assert.Equal(t, "alice@example.com", found.Email)
	})

	t.Run("non existing user", func(t *testing.T) {
		_, err := svc.Find(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestUserService_Update(t *testing.T) {
	svc := setupUserService()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	require.NoError(t, svc.Store(ctx, user))

	t.Run("existing user", func(t *testing.T) {
		user.Name = "Alice Updated"
		user.Email = "alice.updated@example.com"
		err := svc.Update(ctx, user)
		require.NoError(t, err)

		updated, err := svc.Find(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Alice Updated", updated.Name)
		assert.Equal(t, "alice.updated@example.com", updated.Email)
	})

	t.Run("non existing user", func(t *testing.T) {
		ghost := &domain.User{ID: 9999, Name: "Ghost", Email: "ghost@example.com"}
		err := svc.Update(ctx, ghost)
		assert.Error(t, err)
	})
}

func TestUserService_Delete(t *testing.T) {
	svc := setupUserService()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	require.NoError(t, svc.Store(ctx, user))

	t.Run("existing user", func(t *testing.T) {
		err := svc.Delete(ctx, user.ID)
		require.NoError(t, err)

		_, err = svc.Find(ctx, user.ID)
		assert.Error(t, err)
	})

	t.Run("non existing user", func(t *testing.T) {
		err := svc.Delete(ctx, 9999)
		assert.Error(t, err)
	})
}

func TestUserService_List(t *testing.T) {
	svc := setupUserService()

	require.NoError(t, svc.Store(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"}))
	require.NoError(t, svc.Store(ctx, &domain.User{Name: "Bob", Email: "bob@example.com"}))

	users, err := svc.List(ctx, dto.ListFilter{})
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserService_Count(t *testing.T) {
	svc := setupUserService()

	require.NoError(t, svc.Store(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"}))
	require.NoError(t, svc.Store(ctx, &domain.User{Name: "Bob", Email: "bob@example.com"}))

	count, err := svc.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint(2), count)
}