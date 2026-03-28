package responses

import "cafeteria-delivery/internal/core/domain"

type OrdersList struct {
	Data  []domain.Order `json:"data"`
	Count uint           `json:"count"`
}
