package handlers

import (
	"cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/core/ports"
	"fmt"
	"strconv"

	"cafeteria-delivery/internal/base/config"
	_ "cafeteria-delivery/internal/base/core/domain"
	"cafeteria-delivery/internal/base/dto"
	customErrors "cafeteria-delivery/internal/base/errors"
	"cafeteria-delivery/internal/base/http/mappers/model"
	orderRequest "cafeteria-delivery/internal/base/http/requests/order"
	"cafeteria-delivery/internal/base/http/responses"

	"github.com/gofiber/fiber/v2"
)

type OrderHandlers struct {
	serv ports.OrderService
	pgn  *config.Pagination
}

func NewOrderHandlers(service ports.OrderService, pagination *config.Pagination) *OrderHandlers {
	return &OrderHandlers{serv: service, pgn: pagination}
}

var _ ports.OrderHandlers = (*OrderHandlers)(nil)

// Create godoc
// @Summary Create order
// @Tags orders
// @Accept json
// @Produce json
// @Param data body orderRequest.CreateRequest true "Order"
// @Success 201 {object} domain.Order
// @Failure 404 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /orders [post]
func (h *OrderHandlers) Create(ctx *fiber.Ctx) error {
	req := orderRequest.NewCreateRequest()
	if err := req.Parse(ctx); err != nil {
		return err
	}

	order := model.NewOrderModelMapper().NewFromCreateRequest(req)

	if err := h.serv.Store(ctx.UserContext(), &order); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(order)
}

// Show godoc
// @Summary Get order
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandlers) Show(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: order id", customErrors.ErrInvalidID)
	}

	order, err := h.serv.Find(ctx.UserContext(), uint(id))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(order)
}

// List godoc
// @Summary List orders
// @Tags orders
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param order query string false "asc or desc"
// @Param sortBy query string false "Field to sort by"
// @Param user_id query int false "Filter by user"
// @Param status query int false "Filter by status (1=pending 2=confirmed 3=ready 4=delivered 5=cancelled)"
// @Success 200 {object} responses.OrdersList
// @Failure 500 {object} responses.ErrorResponse
// @Router /orders [get]
func (h *OrderHandlers) List(ctx *fiber.Ctx) error {
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

	var userID *uint
	if val := ctx.QueryInt("user_id"); val > 0 {
		u := uint(val)
		userID = &u
	}

	var status *domain.OrderStatus
	if val := ctx.QueryInt("status"); val > 0 {
		s := domain.OrderStatus(val)
		status = &s
	}

	filter := dto.OrderFilter{
		ListFilter: listFilter,
		UserID:     userID,
		Status:     status,
	}

	orders, err := h.serv.List(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	count, err := h.serv.Count(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.OrdersList{Data: orders, Count: count})
}

// Update godoc
// @Summary Update order status
// @Tags orders
// @Accept json
// @Param id path int true "Order ID"
// @Param data body orderRequest.UpdateRequest true "Status"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /orders/{id} [put]
func (h *OrderHandlers) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: order id", customErrors.ErrInvalidID)
	}

	req := orderRequest.NewUpdateRequest()
	if err = req.Parse(ctx); err != nil {
		return err
	}

	order := model.NewOrderModelMapper().NewFromUpdateRequest(req, uint(id))

	if err = h.serv.Update(ctx.UserContext(), &order); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary Delete order
// @Tags orders
// @Param id path int true "Order ID"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /orders/{id} [delete]
func (h *OrderHandlers) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: order id", customErrors.ErrInvalidID)
	}

	if err = h.serv.Delete(ctx.UserContext(), uint(id)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
