package model

import (
	"cafeteria-delivery/internal/core/domain"
	"cafeteria-delivery/internal/http/requests/item_category"
)

type ItemCategoryModelMapper struct{}

func NewItemCategoryModelMapper() ItemCategoryModelMapper {
	return ItemCategoryModelMapper{}
}

func (m ItemCategoryModelMapper) NewFromCreateRequest(req item_category.UpsertRequest) domain.ItemCategory {
	return domain.ItemCategory{
		Name: req.Name,
	}
}

func (m ItemCategoryModelMapper) NewFromUpdateRequest(req item_category.UpsertRequest, id uint) domain.ItemCategory {
	return domain.ItemCategory{
		ID:   id,
		Name: req.Name,
	}
}
