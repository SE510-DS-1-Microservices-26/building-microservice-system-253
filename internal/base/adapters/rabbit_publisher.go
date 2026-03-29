package adapters

import (
	"cafeteria-delivery/pkg/events"
	"cafeteria-delivery/pkg/rabbitmq"
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	ch *amqp.Channel
}

func NewRabbitPublisher(conn *rabbitmq.Connection) *RabbitPublisher {
	return &RabbitPublisher{ch: conn.Ch}
}

// PublishOrderCreated - publishes core-item.created event to rabbit mq
func (p *RabbitPublisher) PublishOrderCreated(ctx context.Context, event events.CoreItemCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.ch.PublishWithContext(ctx, rabbitmq.ExchangeName, events.CoreItemCreatedRoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
