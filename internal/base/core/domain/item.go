package domain

import "time"

type Item struct {
	ID          uint      `json:"id"`
	CategoryID  *uint     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (i *Item) Fields() []interface{} {
	return []interface{}{
		&i.ID,
		&i.CategoryID,
		&i.Name,
		&i.Description,
		&i.ImageURL,
		&i.Price,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	}
}

type ItemTable struct{}

func (ItemTable) TableName() string {
	return "items"
}

func (ItemTable) Columns() []string {
	return []string{
		"items.id",
		"items.category_id",
		"items.name",
		"items.description",
		"items.image_url",
		"items.price",
		"items.quantity",
		"items.created_at",
		"items.updated_at",
	}
}
