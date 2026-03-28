package routes

import (
	_ "cafeteria-delivery/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func DocsRoutes(engine *fiber.App) {
	engine.Get("/docs/*", fiberSwagger.WrapHandler)
}
