package model

import (
	"cafeteria-delivery/internal/base/core/domain"
	orderRequest "cafeteria-delivery/internal/base/http/requests/order"
)

type OrderModelMapper struct{}

func NewOrderModelMapper() OrderModelMapper {
	return OrderModelMapper{}
}

func (m OrderModelMapper) NewFromCreateRequest(req orderRequest.CreateRequest) domain.Order {
	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		})
	}

	return domain.Order{
		UserID: req.UserID,
		Items:  items,
	}
}

func (m OrderModelMapper) NewFromUpdateRequest(req orderRequest.UpdateRequest, id uint) domain.Order {
	return domain.Order{
		ID:     id,
		Status: req.Status,
	}
}
