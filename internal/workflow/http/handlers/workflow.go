package handlers

import (
	"cafeteria-delivery/internal/workflow/core/ports"
	"cafeteria-delivery/internal/workflow/http/requests"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WorkflowHandlers struct {
	serv ports.WorkflowService
}

var _ ports.WorkflowHandlers = (*WorkflowHandlers)(nil)

func NewWorkflowHandlers(service ports.WorkflowService) *WorkflowHandlers {
	return &WorkflowHandlers{serv: service}
}

func (h *WorkflowHandlers) Start(ctx *fiber.Ctx) error {
	action := ctx.Params("action")

	switch action {
	case "place-order":
		req := &requests.PlaceOrderRequest{}
		if err := req.Parse(ctx); err != nil {
			return err
		}

		wf, err := h.serv.StartPlaceOrder(ctx.UserContext(), req.ToPortsRequest())
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusCreated).JSON(wf)

	default:
		return fiber.NewError(fiber.StatusBadRequest, "unknown workflow action: "+action)
	}
}

func (h *WorkflowHandlers) Show(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("workflowId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid workflow id")
	}

	wf, err := h.serv.Find(ctx.UserContext(), id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(wf)
}
