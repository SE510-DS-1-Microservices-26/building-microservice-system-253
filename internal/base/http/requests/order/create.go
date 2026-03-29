package order

import (
	"fmt"

	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ItemRequest struct {
	ItemID   uint `json:"item_id"  validate:"required,min=1"`
	Quantity int  `json:"quantity" validate:"required,min=1"`
}

type CreateRequest struct {
	UserID *uint         `json:"user_id"`
	Items  []ItemRequest `json:"items" validate:"required,min=1,dive"`
}

func NewCreateRequest() CreateRequest {
	return CreateRequest{}
}

func (r *CreateRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	if err := validator.Validate(r); err != nil {
		return fmt.Errorf("%w: %s", customErrors.ErrValidation, err.Error())
	}
	return nil
}
