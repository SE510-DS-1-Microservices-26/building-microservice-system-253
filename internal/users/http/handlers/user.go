package handlers

import (
	"cafeteria-delivery/internal/users/core/ports"
	"fmt"
	"strconv"

	"cafeteria-delivery/internal/users/config"
	_ "cafeteria-delivery/internal/users/core/domain"
	"cafeteria-delivery/internal/users/dto"
	customErrors "cafeteria-delivery/internal/users/errors"
	"cafeteria-delivery/internal/users/http/mappers/model"
	userRequest "cafeteria-delivery/internal/users/http/requests/user"
	"cafeteria-delivery/internal/users/http/responses"

	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	serv ports.UserService
	pgn  *config.Pagination
}

func NewUserHandlers(service ports.UserService, pagination *config.Pagination) *UserHandlers {
	return &UserHandlers{serv: service, pgn: pagination}
}

var _ ports.UserHandlers = (*UserHandlers)(nil)

// Create godoc
// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param data body userRequest.UpsertRequest true "User"
// @Success 201 {object} domain.User
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users [post]
func (h *UserHandlers) Create(ctx *fiber.Ctx) error {
	req := userRequest.NewUpsertRequest()
	if err := req.Parse(ctx); err != nil {
		return err
	}

	user := model.NewUserModelMapper().NewFromUpsertRequest(req)

	if err := h.serv.Store(ctx.UserContext(), &user); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// Show godoc
// @Summary Get user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandlers) Show(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: user id", customErrors.ErrInvalidID)
	}

	user, err := h.serv.Find(ctx.UserContext(), uint(id))
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

// List godoc
// @Summary List users
// @Tags users
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param order query string false "asc or desc"
// @Param sortBy query string false "Field to sort by"
// @Success 200 {object} responses.UsersList
// @Failure 500 {object} responses.ErrorResponse
// @Router /users [get]
func (h *UserHandlers) List(ctx *fiber.Ctx) error {
	limit := uint64(ctx.QueryInt("limit"))
	offset := uint64(ctx.QueryInt("offset"))
	order := dto.Order(ctx.Query("order"))
	sortBy := ctx.Query("sortBy")
	filter := dto.NewListFilter(limit, offset, h.pgn.DefaultLimit(), h.pgn.MaxLimit(), order, sortBy)

	users, err := h.serv.List(ctx.UserContext(), filter)
	if err != nil {
		return err
	}

	count, err := h.serv.Count(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(responses.UsersList{Data: users, Count: count})
}

// Update godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Param id path int true "User ID"
// @Param data body userRequest.UpsertRequest true "User"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 422 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandlers) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: user id", customErrors.ErrInvalidID)
	}

	req := userRequest.NewUpsertRequest()
	if err = req.Parse(ctx); err != nil {
		return err
	}

	user := model.NewUserModelMapper().NewFromUpsertRequestWithID(req, uint(id))

	if err = h.serv.Update(ctx.UserContext(), &user); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Delete godoc
// @Summary Delete user
// @Tags users
// @Param id path int true "User ID"
// @Success 200
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandlers) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return fmt.Errorf("%w: user id", customErrors.ErrInvalidID)
	}

	if err = h.serv.Delete(ctx.UserContext(), uint(id)); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
