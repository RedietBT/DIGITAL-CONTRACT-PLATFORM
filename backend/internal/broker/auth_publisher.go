package broker

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthPublisher struct {
	channel *amqp.Channel
}

func NewAuthPublisher(conn *amqp.Connection) (*AuthPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// We declare the exchange here to ensure it exists before we shout
	err = ch.ExchangeDeclare(
		"user_events", 
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	return &AuthPublisher{channel: ch}, err
}

func (p *AuthPublisher) PublishUserCreated(ctx context.Context, event interface{}) error {
	body, _ := json.Marshal(event)
	return p.channel.PublishWithContext(ctx,
	          "user_events",
              "",
              false,
              false,
              amqp.Publishing{
	            ContentType: "application/json",
	            Body: body,})
}

// PublishUserDeleted sends a message to RabbitMQ indicating a user has been removed.
func (p *AuthPublisher) PublishUserDeleted(ctx context.Context, userID string) error {
	// 1. Prepare the message body (just the ID is enough for deletion)
	body, err := json.Marshal(map[string]string{
		"user_id": userID,
		"action":  "DELETE",
	})
	if err != nil {
		return err
	}

	// 2. Publish to the "user.events" exchange
	// We use a routing key like "user.deleted"
	return p.channel.PublishWithContext(ctx,
		"user_events",  // exchange name
		"user.deleted", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}