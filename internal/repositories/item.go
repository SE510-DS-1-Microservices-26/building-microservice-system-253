package repositories

import (
	"context"
	"fmt"
	"time"

	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/core/ports"
	"cafeteria-delivery/internal/dto"
	customErrors "cafeteria-delivery/internal/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	pool *pgxpool.Pool
}

var _ ports.ItemRepository = (*ItemRepository)(nil)

func NewItemRepository(p *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{pool: p}
}

// Store - create new menu item
func (r *ItemRepository) Store(ctx context.Context, item *domain.Item) error {
	query, args, err := squirrel.Insert(domain.ItemTable{}.TableName()).
		Columns("category_id", "name", "description", "image_url", "price", "quantity").
		Values(item.CategoryID, item.Name, item.Description, item.ImageURL, item.Price, item.Quantity).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if err = r.pool.QueryRow(ctx, query, args...).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return fmt.Errorf("failed to scan item values: %w", err)
	}

	return nil
}

// Find - find a menu item
func (r *ItemRepository) Find(ctx context.Context, id uint) (*domain.Item, error) {
	query, args, err := squirrel.Select(domain.ItemTable{}.Columns()...).
		From(domain.ItemTable{}.TableName()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var item domain.Item
	if err = r.pool.QueryRow(ctx, query, args...).Scan(item.Fields()...); err != nil {
		return nil, fmt.Errorf("failed to scan item values: %w", err)
	}

	return &item, nil
}

// FindByIDs - find menu items by ids
func (r *ItemRepository) FindByIDs(ctx context.Context, ids []uint) ([]domain.Item, error) {
	query, args, err := squirrel.Select(domain.ItemTable{}.Columns()...).
		From(domain.ItemTable{}.TableName()).
		Where(squirrel.Eq{"id": ids}).
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

	items := make([]domain.Item, 0, len(ids))
	for rows.Next() {
		var item domain.Item
		if err = rows.Scan(item.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan item values: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// List - list menu items
func (r *ItemRepository) List(ctx context.Context, filter dto.ItemFilter) ([]domain.Item, error) {
	builder := squirrel.Select(domain.ItemTable{}.Columns()...).
		From(domain.ItemTable{}.TableName()).
		OrderBy(filter.Sort).
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(squirrel.Dollar)

	// check the category filter
	if filter.CategoryID != nil {
		builder = builder.Where(squirrel.Eq{"category_id": *filter.CategoryID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Item, 0)
	for rows.Next() {
		var item domain.Item
		if err = rows.Scan(item.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan item values: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

// Count - count menu items
func (r *ItemRepository) Count(ctx context.Context, filter dto.ItemFilter) (uint, error) {
	builder := squirrel.Select("COUNT(id)").
		From(domain.ItemTable{}.TableName()).
		PlaceholderFormat(squirrel.Dollar)

	// check the category filter
	if filter.CategoryID != nil {
		builder = builder.Where(squirrel.Eq{"category_id": *filter.CategoryID})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build query: %w", err)
	}

	var count uint
	if err = r.pool.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to scan count value: %w", err)
	}

	return count, nil
}

// Update - update menu item
func (r *ItemRepository) Update(ctx context.Context, item *domain.Item) error {
	query, args, err := squirrel.Update(domain.ItemTable{}.TableName()).
		Set("category_id", item.CategoryID).
		Set("name", item.Name).
		Set("description", item.Description).
		Set("image_url", item.ImageURL).
		Set("price", item.Price).
		Set("quantity", item.Quantity).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": item.ID}).
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
		return fmt.Errorf("item %d: %w", item.ID, customErrors.ErrNotFound)
	}

	return nil
}

// Delete - delete menu item
func (r *ItemRepository) Delete(ctx context.Context, id uint) error {
	query, args, err := squirrel.Delete(domain.ItemTable{}.TableName()).
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
		return fmt.Errorf("item %d: %w", id, customErrors.ErrNotFound)
	}

	return nil
}
