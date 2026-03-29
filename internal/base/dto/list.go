package dto

import (
	"fmt"
)

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

func NewListFilter(limit, offset, configDefault, configMax uint64, orderBy Order, sortBy string) ListFilter {
	if limit == 0 {
		limit = configDefault
	}
	if limit > configMax {
		limit = configMax
	}

	if orderBy != ASC && orderBy != DESC {
		orderBy = DESC
	}

	if sortBy == "" {
		sortBy = "id"
	}

	sort := fmt.Sprintf("%s %s", sortBy, orderBy)

	return ListFilter{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
	}
}
