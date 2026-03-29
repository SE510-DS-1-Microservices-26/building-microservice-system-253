package repositories

import (
	"cafeteria-delivery/internal/notification/core/domain"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(p *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: p}
}

// Store - insert new notification
func (r *NotificationRepository) Store(ctx context.Context, n *domain.Notification) error {
	columns := []string{"event_id", "core_item_id", "owner_user_id", "summary", "occurred_at"}
	values := []interface{}{n.EventID, n.CoreItemID, n.OwnerUserID, n.Summary, n.OccurredAt}
	query, args, err := squirrel.Insert(domain.NotificationTable{}.TableName()).
		Columns(columns...).
		Values(values...).
		Suffix("ON CONFLICT (event_id) DO NOTHING RETURNING id, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if err = r.pool.QueryRow(ctx, query, args...).Scan(&n.ID, &n.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("duplicate event, skipping: %s", n.EventID)
			return nil
		}
		return fmt.Errorf("failed to scan notification values: %w", err)
	}

	return nil
}
