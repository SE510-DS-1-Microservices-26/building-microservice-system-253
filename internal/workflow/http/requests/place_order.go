package requests

import (
	"fmt"

	"cafeteria-delivery/internal/workflow/core/ports"
	customErrors "cafeteria-delivery/internal/workflow/errors"
	"cafeteria-delivery/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type PlaceOrderRequest struct {
	UserID uint             `json:"user_id" validate:"required,min=1"`
	Items  []PlaceOrderItem `json:"items"   validate:"required,min=1,dive"`
}

type PlaceOrderItem struct {
	ItemID   uint `json:"item_id"  validate:"required,min=1"`
	Quantity int  `json:"quantity" validate:"required,min=1"`
}

func (r *PlaceOrderRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	if err := validator.Validate(r); err != nil {
		return fmt.Errorf("%w: %s", customErrors.ErrValidation, err.Error())
	}
	return nil
}

func (r *PlaceOrderRequest) ToPortsRequest() ports.PlaceOrderRequest {
	items := make([]ports.PlaceOrderItem, len(r.Items))
	for i, it := range r.Items {
		items[i] = ports.PlaceOrderItem{ItemID: it.ItemID, Quantity: it.Quantity}
	}
	return ports.PlaceOrderRequest{UserID: r.UserID, Items: items}
}
