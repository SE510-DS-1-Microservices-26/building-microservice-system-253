package dto

import "cafeteria-delivery/internal/core/domain"

type OrderFilter struct {
	ListFilter
	UserID *uint
	Status *domain.OrderStatus
}
