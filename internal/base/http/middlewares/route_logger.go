package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	DateTimeLayout = "15:04:05 02-01-2006"
)

func RouteLoggerMiddleware(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${method} ${path} ${status} ${latency} pid=${pid}\n",
		TimeFormat: DateTimeLayout,
	}))
}
