package ports

import (
	"cafeteria-delivery/internal/users/core/domain"
	"context"

	"cafeteria-delivery/internal/users/dto"

	"github.com/gofiber/fiber/v2"
)

type UserHandlers interface {
	Create(ctx *fiber.Ctx) error
	Show(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}

type UserService interface {
	Store(ctx context.Context, user *domain.User) error
	Find(ctx context.Context, id uint) (*domain.User, error)
	List(ctx context.Context, filter dto.ListFilter) ([]domain.User, error)
	Count(ctx context.Context) (uint, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}

type UserRepository interface {
	Store(ctx context.Context, user *domain.User) error
	Find(ctx context.Context, id uint) (*domain.User, error)
	List(ctx context.Context, filter dto.ListFilter) ([]domain.User, error)
	Count(ctx context.Context) (uint, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
}
