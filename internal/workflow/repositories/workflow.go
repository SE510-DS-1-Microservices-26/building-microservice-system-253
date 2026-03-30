package repositories

import (
	"cafeteria-delivery/internal/workflow/core/domain"
	"cafeteria-delivery/internal/workflow/core/ports"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowRepository struct {
	pool *pgxpool.Pool
}

var _ ports.WorkflowRepository = (*WorkflowRepository)(nil)

func NewWorkflowRepository(pool *pgxpool.Pool) *WorkflowRepository {
	return &WorkflowRepository{pool: pool}
}

func (r *WorkflowRepository) Store(ctx context.Context, wf *domain.WorkflowInstance) error {
	query, args, err := squirrel.Insert(domain.WorkflowTable{}.TableName()).
		Columns("workflow_id", "type", "state", "payload", "last_error").
		Values(wf.WorkflowID, wf.Type, wf.State, wf.Payload, wf.LastError).
		Suffix("RETURNING created_at, updated_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build store query: %w", err)
	}

	return r.pool.QueryRow(ctx, query, args...).Scan(&wf.CreatedAt, &wf.UpdatedAt)
}

func (r *WorkflowRepository) Update(ctx context.Context, wf *domain.WorkflowInstance) error {
	query, args, err := squirrel.Update(domain.WorkflowTable{}.TableName()).
		Set("state", wf.State).
		Set("last_error", wf.LastError).
		Set("updated_at", time.Now().UTC()).
		Where(squirrel.Eq{"workflow_id": wf.WorkflowID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

func (r *WorkflowRepository) Find(ctx context.Context, id uuid.UUID) (*domain.WorkflowInstance, error) {
	query, args, err := squirrel.Select(domain.WorkflowTable{}.Columns()...).
		From(domain.WorkflowTable{}.TableName()).
		Where(squirrel.Eq{"workflow_id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build find query: %w", err)
	}

	var wf domain.WorkflowInstance
	if err = r.pool.QueryRow(ctx, query, args...).Scan(wf.Fields()...); err != nil {
		return nil, fmt.Errorf("failed to scan workflow: %w", err)
	}

	return &wf, nil
}
