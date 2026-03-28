package domain

import "time"

type ItemCategory struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *ItemCategory) Fields() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	}
}

type ItemCategoryTable struct{}

func (ItemCategoryTable) TableName() string {
	return "item_categories"
}

func (ItemCategoryTable) Columns() []string {
	return []string{
		"item_categories.id",
		"item_categories.name",
		"item_categories.created_at",
		"item_categories.updated_at",
	}
}
