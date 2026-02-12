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

	go func ()  {
		for d := range msgs {
			var event models.UserDeletedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				continue // Skip if message format is wrong
			}

			// Only act if the action is DELETE
			if event.Action == "DELETE" {
				log.Printf("🗑️ Contract Consumer: Removing data for User %s", event.UserID)
				
				// Using Background context because this is an async process
				ctx := context.Background()

				//We call the repo to delete all contracts ownwd by this user
				err := c.repo.DeleteAllByUserID(ctx, event.UserID)
				if err != nil {
					log.Printf("❌ Failed to delete contracts for user %s: %v", event.UserID, err)
				}
			}
		}
	}()

	return nil
}