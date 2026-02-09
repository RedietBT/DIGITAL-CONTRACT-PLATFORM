package broker

import (
	"context"
	"encoding/json"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserEventHandler interface{
	HandleUserCreatedEvent(ctx context.Context, event models.UserCreatedEvent) error
}

type ProfileConsumer struct {
	channel *amqp.Channel
	handler UserEventHandler
}

func NewProfileConsumer(conn *amqp.Connection, handler UserEventHandler) (*ProfileConsumer, error){
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &ProfileConsumer{
		channel: ch,
		handler: handler,
	}, nil
}

func (c *ProfileConsumer) Start(){
	// Declare exchage and queue
	_= c.channel.ExchangeDeclare(
		"user_events",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	q, _ := c.channel.QueueDeclare(
		"profile_creation_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	_ = c.channel.QueueBind( q.Name, "", "user_events", false, nil)
	
	msgs, _ := c.channel.Consume(q.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			var event models.UserCreatedEvent
			json.Unmarshal(d.Body, &event)
			c.handler.HandleUserCreatedEvent(context.Background(), event)
		}
	}()
}

