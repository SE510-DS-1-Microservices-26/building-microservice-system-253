package order

import (
	"fmt"

	"cafeteria-delivery/internal/base/core/domain"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UpdateRequest struct {
	Status domain.OrderStatus `json:"status" validate:"required,min=1,max=5"`
}

func NewUpdateRequest() UpdateRequest {
	return UpdateRequest{}
}

func (r *UpdateRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	if err := validator.Validate(r); err != nil {
		return fmt.Errorf("%w: %s", customErrors.ErrValidation, err.Error())
	}
	return nil
}
