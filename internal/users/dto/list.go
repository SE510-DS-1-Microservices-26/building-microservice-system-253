package dto

import "fmt"

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type ListFilter struct {
	Limit  uint64
	Offset uint64
	Sort   string
}

func NewListFilter(limit, offset, defaultLimit, maxLimit uint64, orderBy Order, sortBy string) ListFilter {
	if limit == 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	if orderBy != ASC && orderBy != DESC {
		orderBy = DESC
	}
	if sortBy == "" {
		sortBy = "id"
	}

	return ListFilter{
		Limit:  limit,
		Offset: offset,
		Sort:   fmt.Sprintf("%s %s", sortBy, orderBy),
	}
}
