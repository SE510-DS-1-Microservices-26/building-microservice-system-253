package responses

import (
	"cafeteria-delivery/internal/users/core/domain"
)

type UsersList struct {
	Data  []domain.User `json:"data"`
	Count uint          `json:"count"`
}
