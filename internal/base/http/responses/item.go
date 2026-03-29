package responses

import "cafeteria-delivery/internal/base/core/domain"

type ItemsList struct {
	Data  []domain.Item `json:"data"`
	Count uint          `json:"count"`
}
