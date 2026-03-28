package ports

import (
	"context"

	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type ItemHandlers interface {
	Create(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}

type ItemService interface {
	Store(ctx context.Context, item *domain.Item) error
	Find(ctx context.Context, id uint) (*domain.Item, error)
	List(ctx context.Context, filter dto.ItemFilter) ([]domain.Item, error)
	Count(ctx context.Context, filter dto.ItemFilter) (uint, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id uint) error
}

type ItemRepository interface {
	Store(ctx context.Context, item *domain.Item) error
	Find(ctx context.Context, id uint) (*domain.Item, error)
	FindByIDs(ctx context.Context, ids []uint) ([]domain.Item, error)
	List(ctx context.Context, filter dto.ItemFilter) ([]domain.Item, error)
	Count(ctx context.Context, filter dto.ItemFilter) (uint, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id uint) error
}
