package model

import (
	"cafeteria-delivery/internal/core/domain"
	itemRequest "cafeteria-delivery/internal/http/requests/item"
)

type ItemModelMapper struct{}

func NewItemModelMapper() ItemModelMapper {
	return ItemModelMapper{}
}

func (m ItemModelMapper) NewFromCreateRequest(req itemRequest.UpsertRequest) domain.Item {
	return m.newFromUpsertRequest(req)
}

func (m ItemModelMapper) NewFromUpdateRequest(req itemRequest.UpsertRequest, id uint) domain.Item {
	item := m.newFromUpsertRequest(req)
	item.ID = id
	return item
}

func (m ItemModelMapper) newFromUpsertRequest(req itemRequest.UpsertRequest) domain.Item {
	return domain.Item{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
}
