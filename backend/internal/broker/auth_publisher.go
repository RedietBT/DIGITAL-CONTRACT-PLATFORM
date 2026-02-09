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