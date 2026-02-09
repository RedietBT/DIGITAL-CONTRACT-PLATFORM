package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserEventHandler interface {
	HandleUserCreatedEvent(ctx context.Context, event models.UserCreatedEvent) error
	HandleUserDeletedEvent(ctx context.Context, userID string) error
}

type ProfileConsumer struct {
	channel *amqp.Channel
	handler UserEventHandler
}

func NewProfileConsumer(conn *amqp.Connection, handler UserEventHandler) (*ProfileConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &ProfileConsumer{
		channel: ch,
		handler: handler,
	}, nil
}

// Start now contains the merged logic for both Create and Delete
func (c *ProfileConsumer) Start() error {
	// 1. Declare Exchange (Fanout)
	err := c.channel.ExchangeDeclare(
		"user_events",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 2. Declare Queue
	q, err := c.channel.QueueDeclare(
		"profile_creation_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 3. Bind Queue to Exchange
	err = c.channel.QueueBind(q.Name, "", "user_events", false, nil)
	if err != nil {
		return err
	}

	// 4. Start Consuming
	msgs, err := c.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var raw map[string]interface{}
			if err := json.Unmarshal(d.Body, &raw); err != nil {
				log.Printf("❌ Failed to unmarshal message: %v", err)
				continue
			}

			ctx := context.Background()

			// LOGIC: Check 'email' for Create OR 'user_id' for Delete
			if _, isCreate := raw["email"]; isCreate {
				var event models.UserCreatedEvent
				json.Unmarshal(d.Body, &event)
				log.Printf("📥 Consumer: Handling User Created for %s", event.UserID)
				_ = c.handler.HandleUserCreatedEvent(ctx, event)

			} else if userID, ok := raw["user_id"].(string); ok {
				log.Printf("📥 Consumer: Handling User Deleted for %s", userID)
				// This satisfies the requirement to sync the profile_schema.profile table
				_ = c.handler.HandleUserDeletedEvent(ctx, userID)
			}
		}
	}()

	return nil
}