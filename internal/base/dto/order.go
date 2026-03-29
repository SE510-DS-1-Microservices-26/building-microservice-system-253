package dto

import "cafeteria-delivery/internal/base/core/domain"

type OrderFilter struct {
	ListFilter
	UserID *uint
	Status *domain.OrderStatus
}
