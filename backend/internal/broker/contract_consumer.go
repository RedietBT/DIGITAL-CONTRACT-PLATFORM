package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ContractConsumer struct {
	channel *amqp.Channel
	repo    repository.ContractRepository
}

func NewContractConsumer(conn *amqp.Connection, repo repository.ContractRepository) (*ContractConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return  &ContractConsumer{
		channel:  ch,
		repo: repo,
	}, nil
}

func (c *ContractConsumer) Start() error {
	// 1. Declare the Exchange (Must match Auth Service)
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

	// 2. Declare a unique queue for the Contract Service
	q, err := c.channel.QueueDeclare(
		"contract_service_user_queue", // Unique name so Contract gets its own copy of messages
		true, false, false, false, nil,
	)

	// 3. Bind the queue to the exchage
	err = c.channel.QueueBind(q.Name, "", "user_events", false, nil)

	// 4. Start Consuming
	msgs, err := c.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	
	go func() {
        for d := range msgs {
            // 1. Try to unmarshal as UserCreatedEvent first
            var createEvent models.UserCreatedEvent
            if err := json.Unmarshal(d.Body, &createEvent); err == nil && createEvent.UserID != "" {
                log.Printf("👤 Contract Service: Noticed new user created: %s", createEvent.UserID)
                // Optional: You could save this to a local 'users' table if you want 
                // to validate IDs before assigning them to contracts.
                continue 
            }

            // 2. Fallback to your existing Delete logic
            var deleteEvent models.UserDeletedEvent
            json.Unmarshal(d.Body, &deleteEvent)
            if deleteEvent.Action == "DELETE" {
                log.Printf("🗑️ Contract Consumer: Removing data for User %s", deleteEvent.UserID)
                c.repo.DeleteAllByUserID(context.Background(), deleteEvent.UserID)
            }
        }
    }()
    return nil
}