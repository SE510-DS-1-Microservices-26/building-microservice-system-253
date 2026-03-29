package routes

import (
	"cafeteria-delivery/internal/users/core/ports"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(engine fiber.Router, handler ports.UserHandlers) {
	users := engine.Group("/users")

	users.Get("", handler.List).Name("listUsers")
	users.Get("/:id", handler.Show).Name("showUser")
	users.Post("", handler.Create).Name("createUser")
	users.Put("/:id", handler.Update).Name("updateUser")
	users.Delete("/:id", handler.Delete).Name("deleteUser")
}
