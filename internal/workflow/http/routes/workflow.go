package routes

import (
	"cafeteria-delivery/internal/workflow/core/ports"

	"github.com/gofiber/fiber/v2"
)

func WorkflowRoutes(router fiber.Router, h ports.WorkflowHandlers) {
	router.Post("/workflows/:action", h.Start)
	router.Get("/workflows/:workflowId", h.Show)
}
