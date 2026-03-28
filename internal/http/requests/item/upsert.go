package item

import (
	"fmt"

	customErrors "cafeteria-delivery/internal/errors"
	"cafeteria-delivery/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type UpsertRequest struct {
	CategoryID  *uint   `json:"category_id"  validate:"omitempty,min=1"`
	Name        string  `json:"name"         validate:"required,max=255"`
	Description string  `json:"description"  validate:"omitempty,max=1000"`
	ImageURL    string  `json:"image_url"    validate:"omitempty,url"`
	Price       float64 `json:"price"        validate:"required,gt=0"`
	Quantity    int     `json:"quantity"     validate:"min=0"`
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
