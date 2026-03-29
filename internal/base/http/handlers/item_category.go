package handlers

import (
	_ "cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/core/ports"
	"fmt"
	"strconv"

	"cafeteria-delivery/internal/base/config"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/internal/base/http/mappers/model"
	itemCategoryRequest "cafeteria-delivery/internal/base/http/requests/item_category"
	"cafeteria-delivery/internal/base/http/responses"

	"github.com/gofiber/fiber/v2"
)

type ItemCategoryHandlers struct {
	serv ports.ItemCategoryService
	pgn  *config.Pagination
}

func NewItemCategoryHandlers(service ports.ItemCategoryService, pagination *config.Pagination) *ItemCategoryHandlers {
	return &ItemCategoryHandlers{serv: service, pgn: pagination}
}

var _ ports.ItemCategoryHandlers = (*ItemCategoryHandlers)(nil)

// Create godoc
// @Summary Create item category
// @Tags item-categories
// @Accept json
// @Produce json
// @Param data body itemCategoryRequest.UpsertRequest true "Category"
// @Success 201 {object} domain.ItemCategory
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /item-categories [post]
func (h *ItemCategoryHandlers) Create(ctx *fiber.Ctx) error {
	req := itemCategoryRequest.NewUpsertRequest()
	if err := req.Parse(ctx); err != nil {
		return err
	}

	category := model.NewItemCategoryModelMapper().NewFromCreateRequest(req)

	if err := h.serv.Store(ctx.UserContext(), &category); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(category)
}

// Show godoc
// @Summary Get item category
// @Tags item-categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.ItemCategory
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /item-categories/{id} [get]
func (h *ItemCategoryHandlers) Show(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: category id", customErrors.ErrInvalidID)
	}

	category, err := h.serv.Find(ctx.UserContext(), uint(id))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(category)
}

// List godoc
// @Summary List item categories
// @Tags item-categories
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param order query string false "asc or desc"
// @Param sortBy query string false "Field to sort by"
// @Success 200 {object} responses.ItemCategoriesList
// @Failure 500 {object} responses.ErrorResponse
// @Router /item-categories [get]
func (h *ItemCategoryHandlers) List(ctx *fiber.Ctx) error {
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

	categories, err := h.serv.List(ctx.UserContext(), listFilter)
	if err != nil {
		return err
	}

	count, err := h.serv.Count(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.ItemCategoriesList{Data: categories, Count: count})
}

// Update godoc
// @Summary Update item category
// @Tags item-categories
// @Accept json
// @Param id path int true "Category ID"
// @Param data body itemCategoryRequest.UpsertRequest true "Category"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /item-categories/{id} [put]
func (h *ItemCategoryHandlers) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: category id", customErrors.ErrInvalidID)
	}

	req := itemCategoryRequest.NewUpsertRequest()
	if err = req.Parse(ctx); err != nil {
		return err
	}

	category := model.NewItemCategoryModelMapper().NewFromUpdateRequest(req, uint(id))

	if err = h.serv.Update(ctx.UserContext(), &category); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary Delete item category
// @Tags item-categories
// @Param id path int true "Category ID"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /item-categories/{id} [delete]
func (h *ItemCategoryHandlers) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: category id", customErrors.ErrInvalidID)
	}

	if err = h.serv.Delete(ctx.UserContext(), uint(id)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
