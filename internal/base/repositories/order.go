package repositories

import (
	"cafeteria-delivery/internal/base/core/ports"
	"context"
	"fmt"
	"time"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

var _ ports.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(p *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: p}
}

// Store - create new order with order items
func (r *OrderRepository) Store(ctx context.Context, order *domain.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// build and execute order query
	orderQuery, orderArgs, err := squirrel.Insert(domain.OrderTable{}.TableName()).
		Columns("user_id", "status").
		Values(order.UserID, order.Status).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build order query: %w", err)
	}

	if err = tx.QueryRow(ctx, orderQuery, orderArgs...).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return fmt.Errorf("failed to scan order values: %w", err)
	}

	// build and execute order items query
	itemsInsert := squirrel.Insert(domain.OrderItemTable{}.TableName()).
		Columns("order_id", "item_id", "quantity", "unit_price").
		PlaceholderFormat(squirrel.Dollar)

	for i, item := range order.Items {
		itemsInsert = itemsInsert.Values(order.ID, item.ItemID, item.Quantity, item.UnitPrice)
		order.Items[i].OrderID = order.ID
	}

	itemsQuery, itemsArgs, err := itemsInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build order items query: %w", err)
	}

	if _, err = tx.Exec(ctx, itemsQuery, itemsArgs...); err != nil {
		return fmt.Errorf("failed to insert order items: %w", err)
	}

	// collect items ids and quantities to decrement in one batch update
	itemIDs := make([]uint, len(order.Items))
	quantities := make([]int, len(order.Items))
	for i, item := range order.Items {
		itemIDs[i] = item.ItemID
		quantities[i] = item.Quantity
	}

	decrementQuery := `
		UPDATE items
		SET quantity = items.quantity - delta.qty
		FROM (SELECT unnest($1::int[]) AS id, unnest($2::int[]) AS qty) AS delta
		WHERE items.id = delta.id`

	if _, err = tx.Exec(ctx, decrementQuery, itemIDs, quantities); err != nil {
		return fmt.Errorf("failed to decrement items quantities: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Find - find an order with order items
func (r *OrderRepository) Find(ctx context.Context, id uint) (*domain.Order, error) {
	orderQuery, orderArgs, err := squirrel.Select(domain.OrderTable{}.Columns()...).
		From(domain.OrderTable{}.TableName()).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var order domain.Order
	if err = r.pool.QueryRow(ctx, orderQuery, orderArgs...).Scan(order.Fields()...); err != nil {
		return nil, fmt.Errorf("failed to scan order values: %w", err)
	}

	itemsQuery, itemsArgs, err := squirrel.Select(domain.OrderItemTable{}.Columns()...).
		From(domain.OrderItemTable{}.TableName()).
		Where(squirrel.Eq{"order_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build order items query: %w", err)
	}

	rows, err := r.pool.Query(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute order items query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err = rows.Scan(item.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan order item values: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}

// List - list the orders
func (r *OrderRepository) List(ctx context.Context, filter dto.OrderFilter) ([]domain.Order, error) {
	// fetch orders data first
	builder := squirrel.Select(domain.OrderTable{}.Columns()...).
		From(domain.OrderTable{}.TableName()).
		OrderBy(filter.Sort).
		Limit(filter.Limit).
		Offset(filter.Offset).
		PlaceholderFormat(squirrel.Dollar)

	// check the user filter
	if filter.UserID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": *filter.UserID})
	}
	// check the status filter
	if filter.Status != nil {
		builder = builder.Where(squirrel.Eq{"status": *filter.Status})
	}

	// build a query
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// execute query to fetch orders
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// scan orders values
	orders := make([]domain.Order, 0)
	for rows.Next() {
		var order domain.Order
		if err = rows.Scan(order.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan order values: %w", err)
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return orders, nil
	}

	// extract orders ids
	orderIDs := make([]uint, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.ID
	}

	// build order items query
	itemsQuery, itemsArgs, err := squirrel.Select(domain.OrderItemTable{}.Columns()...).
		From(domain.OrderItemTable{}.TableName()).
		Where(squirrel.Eq{"order_id": orderIDs}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build order items query: %w", err)
	}

	// execute order items query
	itemRows, err := r.pool.Query(ctx, itemsQuery, itemsArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute order items query: %w", err)
	}
	defer itemRows.Close()

	// scan order items values
	items := make(map[uint][]domain.OrderItem)
	for itemRows.Next() {
		var item domain.OrderItem
		if err = itemRows.Scan(item.Fields()...); err != nil {
			return nil, fmt.Errorf("failed to scan order item values: %w", err)
		}
		items[item.OrderID] = append(items[item.OrderID], item)
	}

	for i, o := range orders {
		orders[i].Items = items[o.ID]
	}

	return orders, nil
}

// Count - count the orders
func (r *OrderRepository) Count(ctx context.Context, filter dto.OrderFilter) (uint, error) {
	builder := squirrel.Select("COUNT(id)").
		From(domain.OrderTable{}.TableName()).
		PlaceholderFormat(squirrel.Dollar)

	// check the user filter
	if filter.UserID != nil {
		builder = builder.Where(squirrel.Eq{"user_id": *filter.UserID})
	}
	// check the status filter
	if filter.Status != nil {
		builder = builder.Where(squirrel.Eq{"status": *filter.Status})
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

// Update - update the order
func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query, args, err := squirrel.Update(domain.OrderTable{}.TableName()).
		Set("status", order.Status).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"id": order.ID}).
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
		return fmt.Errorf("order %d: %w", order.ID, customErrors.ErrNotFound)
	}

	return nil
}

// Delete - delete the order
func (r *OrderRepository) Delete(ctx context.Context, id uint) error {
	query, args, err := squirrel.Delete(domain.OrderTable{}.TableName()).
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
		return fmt.Errorf("order %d: %w", id, customErrors.ErrNotFound)
	}

	return nil
}
