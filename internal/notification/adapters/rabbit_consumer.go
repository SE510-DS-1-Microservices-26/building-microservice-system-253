package adapters

import (
	"cafeteria-delivery/internal/notification/core/domain"
	"cafeteria-delivery/internal/notification/core/services"
	"cafeteria-delivery/pkg/events"
	"cafeteria-delivery/pkg/rabbitmq"
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "notification.core-item.created"

type RabbitConsumer struct {
	ch      *amqp.Channel
	service *services.NotificationService
}

func NewRabbitConsumer(conn *rabbitmq.Connection, service *services.NotificationService) (*RabbitConsumer, error) {
	ch := conn.Ch

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err = ch.QueueBind(q.Name, events.CoreItemCreatedRoutingKey, rabbitmq.ExchangeName, false, nil); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	c := &RabbitConsumer{ch: ch, service: service}

	return c, nil
}

// Start - starts consuming messages
func (c *RabbitConsumer) Start(ctx context.Context) error {
	messages, err := c.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	log.Println("notification consumer started...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-messages:
			if !ok {
				return fmt.Errorf("rabbitmq channel closed")
			}
			c.handle(ctx, msg)
		}
	}
}

// handle - helper method to process rabbit mq delivery
func (c *RabbitConsumer) handle(ctx context.Context, msg amqp.Delivery) {
	// unmarshal item-created event payload
	var event events.CoreItemCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("failed to unmarshal event: %v — nacking", err)
		if err = msg.Nack(false, false); err != nil {
			log.Printf("failed to nack message: %v", err)
		}
		return
	}

	// create the new notification record
	n := &domain.Notification{
		EventID:     event.EventID,
		CoreItemID:  event.CoreItemID,
		OwnerUserID: event.OwnerUserID,
		Summary:     event.Summary,
		OccurredAt:  event.OccurredAt,
	}

	// store new record
	if err := c.service.Store(ctx, n); err != nil {
		log.Printf("failed to store notification for event %s: %v", event.EventID, err)
		if err = msg.Nack(false, false); err != nil {
			log.Printf("failed to nack message: %v", err)
		}
		return
	}

	// acknowledge the message
	if err := msg.Ack(false); err != nil {
		log.Printf("failed to ack message: %v", err)
	}
}
