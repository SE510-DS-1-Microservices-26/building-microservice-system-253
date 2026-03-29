package events

import "time"

const CoreItemCreatedRoutingKey = "core-item.created"

type CoreItemCreatedEvent struct {
	EventID       string    `json:"event_id"`
	OccurredAt    time.Time `json:"occurred_at"`
	CorrelationID string    `json:"correlation_id"`
	CoreItemID    uint      `json:"core_item_id"`
	OwnerUserID   *uint     `json:"owner_user_id"`
	Summary       string    `json:"summary"`
}
