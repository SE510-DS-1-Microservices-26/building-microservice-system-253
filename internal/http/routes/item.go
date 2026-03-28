package routes

import (
	"cafeteria-delivery/internal/core/ports"

	"github.com/gofiber/fiber/v2"
)

func ItemRoutes(engine *fiber.App, handler ports.ItemHandlers) {
	items := engine.Group("/items")

	items.Get("", handler.List).Name("listItems")
	items.Get("/:id", handler.Show).Name("showItem")
	items.Post("", handler.Create).Name("createItem")
	items.Put("/:id", handler.Update).Name("updateItem")
	items.Delete("/:id", handler.Delete).Name("deleteItem")
}
