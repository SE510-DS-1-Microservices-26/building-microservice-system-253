package ports

import (
	"context"

	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type OrderHandlers interface {
	Create(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}

type OrderService interface {
	Store(ctx context.Context, order *domain.Order) error
	Find(ctx context.Context, id uint) (*domain.Order, error)
	List(ctx context.Context, filter dto.OrderFilter) ([]domain.Order, error)
	Count(ctx context.Context, filter dto.OrderFilter) (uint, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id uint) error
}

type OrderRepository interface {
	Store(ctx context.Context, order *domain.Order) error
	Find(ctx context.Context, id uint) (*domain.Order, error)
	List(ctx context.Context, filter dto.OrderFilter) ([]domain.Order, error)
	Count(ctx context.Context, filter dto.OrderFilter) (uint, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id uint) error
}
