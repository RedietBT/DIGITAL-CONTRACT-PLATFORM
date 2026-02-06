package broker

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProvider struct {
	conn *amqp.Connection // The physical cable connecting your app to RabbitMQ
	channel *amqp.channel // A "virtual" connection inside the cable (efficient for many tasks)
}

// NewRabbitMQProvider connects to the broker and declares the Exchange
func NewRabbitMQProvider(url string) (*RabbitMQProvider, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("faild to open a channel: %w", err)
	}

	// Declare the Exchange (The Post Office)
	// We use 'direct' routing to send messages to specific queues based on a key
	err = ch.ExchangeDeclare(
		"user_events",   //Exchange name
		"direct",       //Type
		true,          // Durable: stays alive if RabbitMQ restarts
		false,        // Auto-deleted
		false,       // Internal
		false,      // No-wait
		nil,       // Arguments
	)

	return &RabbitMQProvider{conn: conn, channel: ch}, err
}

// PublishUserCreated sends a JSON payload to the exchange
func (r *RabbitMQProvider) PublishUserCreated(ctx context.Context, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return r.channel.PublishWithContext(ctx,
	"user_events",
	routingKey,
	false,
	false,
	amqp.Publishing{
		ContentType: "application/json",
		Body: body,
	})
}

	func (r *RabbitMQProvider) Close() {
		r.channel.Close()
		r.conn.Close()
	}