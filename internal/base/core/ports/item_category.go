package ports

import (
	"context"

	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"

	"github.com/gofiber/fiber/v2"
)

type ItemCategoryHandlers interface {
	Create(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}

type ItemCategoryService interface {
	Store(ctx context.Context, category *domain.ItemCategory) error
	Find(ctx context.Context, id uint) (*domain.ItemCategory, error)
	List(ctx context.Context, filter dto.ListFilter) ([]domain.ItemCategory, error)
	Count(ctx context.Context) (uint, error)
	Update(ctx context.Context, category *domain.ItemCategory) error
	Delete(ctx context.Context, id uint) error
}

type ItemCategoryRepository interface {
	Store(ctx context.Context, category *domain.ItemCategory) error
	Find(ctx context.Context, id uint) (*domain.ItemCategory, error)
	List(ctx context.Context, filter dto.ListFilter) ([]domain.ItemCategory, error)
	Count(ctx context.Context) (uint, error)
	Update(ctx context.Context, category *domain.ItemCategory) error
	Delete(ctx context.Context, id uint) error
}
