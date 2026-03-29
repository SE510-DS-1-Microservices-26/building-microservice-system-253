package domain

import (
	"time"
)

type Notification struct {
	ID          uint      `json:"id"`
	EventID     string    `json:"event_id"`
	CoreItemID  uint      `json:"core_item_id"`
	OwnerUserID *uint     `json:"owner_user_id"`
	Summary     string    `json:"summary"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (n *Notification) Fields() []interface{} {
	return []interface{}{
		&n.ID,
		&n.EventID,
		&n.CoreItemID,
		&n.OwnerUserID,
		&n.Summary,
		&n.OccurredAt,
		&n.CreatedAt,
	}
}

type NotificationTable struct{}

func (NotificationTable) TableName() string {
	return "notifications"
}

func (NotificationTable) Columns() []string {
	return []string{
		"notifications.id",
		"notifications.event_id",
		"notifications.core_item_id",
		"notifications.owner_user_id",
		"notifications.summary",
		"notifications.occurred_at",
		"notifications.created_at",
	}
}
