package model

import (
	"cafeteria-delivery/internal/users/core/domain"
	userRequest "cafeteria-delivery/internal/users/http/requests/user"
)

type UserModelMapper struct{}

func NewUserModelMapper() *UserModelMapper {
	return &UserModelMapper{}
}

func (m *UserModelMapper) NewFromUpsertRequest(req *userRequest.UpsertRequest) domain.User {
	return domain.User{
		Name:  req.Name,
		Email: req.Email,
	}
}

func (m *UserModelMapper) NewFromUpsertRequestWithID(req *userRequest.UpsertRequest, id uint) domain.User {
	return domain.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}
}
