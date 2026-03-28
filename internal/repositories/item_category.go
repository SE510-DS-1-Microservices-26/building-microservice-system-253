package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/core/ports"
	"cafeteria-delivery/internal/dto"
	customErrors "cafeteria-delivery/internal/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemCategoryRepository struct {
	pool *pgxpool.Pool
}

var _ ports.ItemCategoryRepository = (*ItemCategoryRepository)(nil)

func NewItemCategoryRepository(p *pgxpool.Pool) *ItemCategoryRepository {
	return &ItemCategoryRepository{pool: p}
}

// Store - create new item category
func (r *ItemCategoryRepository) Store(ctx context.Context, category *domain.ItemCategory) error {
	query, args, err := squirrel.Insert(domain.ItemCategoryTable{}.TableName()).
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if err = r.pool.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("category with such name already exists: %w", customErrors.ErrConflict)
		}
		return fmt.Errorf("failed to scan category values: %w", err)
	}

	return nil
}

// Find - find the category
func (r *ItemCategoryRepository) Find(ctx context.Context, id uint) (*domain.ItemCategory, error) {
	query, args, err := squirrel.Select(domain.ItemCategoryTable{}.Columns()...).
		From(domain.ItemCategoryTable{}.TableName()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var category domain.ItemCategory
	if err = r.pool.QueryRow(ctx, query, args...).Scan(category.Fields()...); err != nil {
		return nil, fmt.Errorf("failed to scan category values: %w", err)
	}

	return &category, nil
}

// List - list categories
func (r *ItemCategoryRepository) List(ctx context.Context, filter dto.ListFilter) ([]domain.ItemCategory, error) {
	query, args, err := squirrel.Select(domain.ItemCategoryTable{}.Columns()...).
		From(domain.ItemCategoryTable{}.TableName()).
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

	categories := make([]domain.ItemCategory, 0)
	for rows.Next() {
		var category domain.ItemCategory
		if err = rows.Scan(category.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan category values: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Count - count categories
func (r *ItemCategoryRepository) Count(ctx context.Context) (uint, error) {
	query, args, err := squirrel.Select("COUNT(id)").
		From(domain.ItemCategoryTable{}.TableName()).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var count uint
	if err = r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count value: %w", err)
	}

	return count, nil
}

// Update - update the category
func (r *ItemCategoryRepository) Update(ctx context.Context, category *domain.ItemCategory) error {
	query, args, err := squirrel.Update(domain.ItemCategoryTable{}.TableName()).
		Set("name", category.Name).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": category.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("category with such name already exists: %w", customErrors.ErrConflict)
		}
		return fmt.Errorf("failed to execute query: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category %d: %w", category.ID, customErrors.ErrNotFound)
	}

	return nil
}

// Delete - delete the category
func (r *ItemCategoryRepository) Delete(ctx context.Context, id uint) error {
	query, args, err := squirrel.Delete(domain.ItemCategoryTable{}.TableName()).
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
		return fmt.Errorf("category %d: %w", id, customErrors.ErrNotFound)
	}

	return nil
}
