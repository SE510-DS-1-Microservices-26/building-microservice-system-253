package handlers

import (
	"errors"

	"cafeteria-delivery/internal/base/config"
	customErrors "cafeteria-delivery/internal/base/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

func newTestPagination() *config.Pagination {
	v := viper.New()
	v.Set("api.default_limit", 15)
	v.Set("api.max_allowed_limit", 100)
	return config.NewPaginationConfig(v)
}

func newTestApp() *fiber.App {
	return fiber.New(fiber.Config{ErrorHandler: testErrorHandler})
}

func testErrorHandler(ctx *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{"message": fiberErr.Message})
	}
	switch {
	case errors.Is(err, customErrors.ErrNotFound):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrValidation):
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrConflict):
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, customErrors.ErrInvalidID):
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	case errors.Is(err, pgx.ErrNoRows):
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": customErrors.ErrNotFound.Error()})
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "internal server error"})
}
