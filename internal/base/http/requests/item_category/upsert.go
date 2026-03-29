package item_category

import (
	"fmt"

	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UpsertRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

func NewUpsertRequest() UpsertRequest {
	return UpsertRequest{}
}

func (r *UpsertRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}
	if err := validator.Validate(r); err != nil {
		return fmt.Errorf("%w: %s", customErrors.ErrValidation, err.Error())
	}
	return nil
}
