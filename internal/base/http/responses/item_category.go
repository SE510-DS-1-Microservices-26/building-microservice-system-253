package responses

import "cafeteria-delivery/internal/base/core/domain"

type ItemCategoriesList struct {
	Data  []domain.ItemCategory `json:"data"`
	Count uint                  `json:"count"`
}
