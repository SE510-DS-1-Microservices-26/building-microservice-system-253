package domain

import "time"

type OrderStatus int16

const (
	OrderStatusPending OrderStatus = iota + 1
	OrderStatusConfirmed
	OrderStatusReady
	OrderStatusDelivered
	OrderStatusCancelled
)

type Order struct {
	ID        uint        `json:"id"`
	UserID    *uint       `json:"user_id"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
}

func (o *Order) Fields() []interface{} {
	return []interface{}{
		&o.ID,
		&o.UserID,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	}
}

type OrderTable struct{}

func (OrderTable) TableName() string {
	return "orders"
}

func (OrderTable) Columns() []string {
	return []string{
		"orders.id",
		"orders.user_id",
		"orders.status",
		"orders.created_at",
		"orders.updated_at",
	}
}

// == OrderItem ==

type OrderItem struct {
	OrderID   uint    `json:"order_id"`
	ItemID    uint    `json:"item_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func (oi *OrderItem) Fields() []interface{} {
	return []interface{}{
		&oi.OrderID,
		&oi.ItemID,
		&oi.Quantity,
		&oi.UnitPrice,
	}
}

type OrderItemTable struct{}

func (OrderItemTable) TableName() string {
	return "order_items"
}

func (OrderItemTable) Columns() []string {
	return []string{
		"order_id",
		"item_id",
		"quantity",
		"unit_price",
	}
}
