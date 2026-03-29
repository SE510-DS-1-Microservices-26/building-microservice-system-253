package routes

import (
	"cafeteria-delivery/internal/base/core/ports"

	"github.com/gofiber/fiber/v2"
)

func ItemCategoryRoutes(engine fiber.Router, handler ports.ItemCategoryHandlers) {
	categories := engine.Group("/item-categories")

	categories.Get("", handler.List).Name("listItemCategories")
	categories.Get("/:id", handler.Show).Name("showItemCategory")
	categories.Post("", handler.Create).Name("createItemCategory")
	categories.Put("/:id", handler.Update).Name("updateItemCategory")
	categories.Delete("/:id", handler.Delete).Name("deleteItemCategory")
}
