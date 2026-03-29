package repositories

import (
	"cafeteria-delivery/internal/users/core/domain"
	"cafeteria-delivery/internal/users/core/ports"
	"context"
	"errors"
	"fmt"
	"time"

	"cafeteria-delivery/internal/users/dto"
	customErrors "cafeteria-delivery/internal/users/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

var _ ports.UserRepository = (*UserRepository)(nil)

func NewUserRepository(p *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: p}
}

// Store - create new user
func (r *UserRepository) Store(ctx context.Context, user *domain.User) error {
	query, args, err := squirrel.Insert(domain.UserTable{}.TableName()).
		Columns("name", "email").
		Values(user.Name, user.Email).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if err = r.pool.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("user with such email already exists: %w", customErrors.ErrConflict)
		}
		return fmt.Errorf("failed to scan values: %w", err)
	}

	return nil
}

// Find - find user
func (r *UserRepository) Find(ctx context.Context, id uint) (*domain.User, error) {
	query, args, err := squirrel.Select(domain.UserTable{}.Columns()...).
		From(domain.UserTable{}.TableName()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user domain.User
	if err = r.pool.QueryRow(ctx, query, args...).Scan(user.Fields()...); err != nil {
		return nil, fmt.Errorf("failed to scan values: %w", err)
	}

	return &user, nil
}

// List - list users
func (r *UserRepository) List(ctx context.Context, filter dto.ListFilter) ([]domain.User, error) {
	query, args, err := squirrel.Select(domain.UserTable{}.Columns()...).
		From(domain.UserTable{}.TableName()).
		OrderBy(filter.Sort).
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err = rows.Scan(u.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan values: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// Count - count users
func (r *UserRepository) Count(ctx context.Context) (uint, error) {
	query, args, err := squirrel.Select("COUNT(id)").
		From(domain.UserTable{}.TableName()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var count uint
	if err = r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count: %w", err)
	}

	return count, nil
}

// Update - update user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query, args, err := squirrel.Update(domain.UserTable{}.TableName()).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": user.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("user with such email already exists: %w", customErrors.ErrConflict)
		}
		return fmt.Errorf("failed to execute query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d: %w", user.ID, customErrors.ErrNotFound)
	}

	return nil
}

// Delete - delete user
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	query, args, err := squirrel.Delete(domain.UserTable{}.TableName()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user %d: %w", id, customErrors.ErrNotFound)
	}

	return nil
}
