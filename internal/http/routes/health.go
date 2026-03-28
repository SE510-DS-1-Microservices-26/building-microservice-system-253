package routes

import (
	"github.com/gofiber/fiber/v2"
)

func HealthRoutes(engine *fiber.App) {
	engine.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
}
