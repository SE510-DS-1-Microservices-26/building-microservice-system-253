package handlers

import (
	"cafeteria-delivery/internal/base/core/ports"
	"fmt"
	"strconv"

	"cafeteria-delivery/internal/base/config"
	_ "cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/internal/base/http/mappers/model"
	itemRequest "cafeteria-delivery/internal/base/http/requests/item"
	"cafeteria-delivery/internal/base/http/responses"

	"github.com/gofiber/fiber/v2"
)

type ItemHandlers struct {
	serv ports.ItemService
	pgn  *config.Pagination
}

func NewItemHandlers(service ports.ItemService, pagination *config.Pagination) *ItemHandlers {
	return &ItemHandlers{serv: service, pgn: pagination}
}

var _ ports.ItemHandlers = (*ItemHandlers)(nil)

// Create godoc
// @Summary Create item
// @Tags items
// @Accept json
// @Produce json
// @Param data body itemRequest.UpsertRequest true "Item"
// @Success 201 {object} domain.Item
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /items [post]
func (h *ItemHandlers) Create(ctx *fiber.Ctx) error {
	req := itemRequest.NewUpsertRequest()
	if err := req.Parse(ctx); err != nil {
		return err
	}

	item := model.NewItemModelMapper().NewFromCreateRequest(req)

	if err := h.serv.Store(ctx.UserContext(), &item); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(item)
}

// Show godoc
// @Summary Get item
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} domain.Item
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /items/{id} [get]
func (h *ItemHandlers) Show(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: item id", customErrors.ErrInvalidID)
	}

	item, err := h.serv.Find(ctx.UserContext(), uint(id))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(item)
}

// List godoc
// @Summary List items
// @Tags items
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param order query string false "asc or desc"
// @Param sortBy query string false "Field to sort by"
// @Param category_id query int false "Filter by category"
// @Success 200 {object} responses.ItemsList
// @Failure 500 {object} responses.ErrorResponse
// @Router /items [get]
func (h *ItemHandlers) List(ctx *fiber.Ctx) error {
	limit := uint64(ctx.QueryInt("limit"))
	offset := uint64(ctx.QueryInt("offset"))
	order := dto.Order(ctx.Query("order"))
	sortBy := ctx.Query("sortBy")
	listFilter := dto.NewListFilter(
		limit,
		offset,
		h.pgn.DefaultLimit(),
		h.pgn.MaxLimit(),
		order,
		sortBy,
	)

	var categoryID *uint
	if val := ctx.QueryInt("category_id"); val > 0 {
		c := uint(val)
		categoryID = &c
	}

	filter := dto.ItemFilter{
		ListFilter: listFilter,
		CategoryID: categoryID,
	}

	items, err := h.serv.List(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	count, err := h.serv.Count(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.ItemsList{Data: items, Count: count})
}

// Update godoc
// @Summary Update item
// @Tags items
// @Accept json
// @Param id path int true "Item ID"
// @Param data body itemRequest.UpsertRequest true "Item"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /items/{id} [put]
func (h *ItemHandlers) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: item id", customErrors.ErrInvalidID)
	}

	req := itemRequest.NewUpsertRequest()
	if err = req.Parse(ctx); err != nil {
		return err
	}

	item := model.NewItemModelMapper().NewFromUpdateRequest(req, uint(id))

	if err = h.serv.Update(ctx.UserContext(), &item); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary Delete item
// @Tags items
// @Param id path int true "Item ID"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /items/{id} [delete]
func (h *ItemHandlers) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: item id", customErrors.ErrInvalidID)
	}

	if err = h.serv.Delete(ctx.UserContext(), uint(id)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
