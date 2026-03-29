package routes

import (
	"cafeteria-delivery/internal/base/core/ports"

	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(engine fiber.Router, handler ports.OrderHandlers) {
	orders := engine.Group("/orders")

	orders.Get("", handler.List).Name("listOrders")
	orders.Get("/:id", handler.Show).Name("showOrder")
	orders.Post("", handler.Create).Name("createOrder")
	orders.Put("/:id", handler.Update).Name("updateOrder")
	orders.Delete("/:id", handler.Delete).Name("deleteOrder")
}
