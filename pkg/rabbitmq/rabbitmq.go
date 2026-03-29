package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn *amqp.Connection
	Ch   *amqp.Channel
}

const ExchangeName = "cafeteria"

// Dial - initialize connection and channel, declares exchange
func Dial(url string) (*Connection, error) {
	// connect to rabbit mq
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	// open a new channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// declare a new, durable topic exchange, named cafeteria
	if err = ch.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	c := &Connection{conn: conn, Ch: ch}
	return c, nil
}

// Close - releases channel and rabbit mq connection
func (c *Connection) Close() {
	c.Ch.Close()
	c.conn.Close()
}
