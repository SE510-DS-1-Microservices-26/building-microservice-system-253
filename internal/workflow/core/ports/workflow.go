package ports

import (
	"cafeteria-delivery/internal/workflow/core/domain"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorkflowHandlers interface {
	Start(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
}

type WorkflowService interface {
	StartPlaceOrder(ctx context.Context, req PlaceOrderRequest) (*domain.WorkflowInstance, error)
	Find(ctx context.Context, id uuid.UUID) (*domain.WorkflowInstance, error)
}

type WorkflowRepository interface {
	Store(ctx context.Context, wf *domain.WorkflowInstance) error
	Update(ctx context.Context, wf *domain.WorkflowInstance) error
	Find(ctx context.Context, id uuid.UUID) (*domain.WorkflowInstance, error)
}

type UsersClientPort interface {
	ValidateUser(ctx context.Context, userID uint) error
}

type CoreClientPort interface {
	CreateOrder(ctx context.Context, userID uint, items []CreateOrderItem) (uint, error)
	CancelOrder(ctx context.Context, orderID uint) error
}

type CreateOrderItem struct {
	ItemID   uint
	Quantity int
}

type PlaceOrderRequest struct {
	UserID uint             `json:"user_id"`
	Items  []PlaceOrderItem `json:"items"`
}

type PlaceOrderItem struct {
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
}
