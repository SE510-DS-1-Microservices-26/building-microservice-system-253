package responses

import "cafeteria-delivery/internal/core/domain"

type ItemCategoriesList struct {
	Data  []domain.ItemCategory `json:"data"`
	Count uint                  `json:"count"`
}
