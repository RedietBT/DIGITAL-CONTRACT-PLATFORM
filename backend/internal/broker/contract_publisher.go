package broker

import (
	"encoding/json"
	"log"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ContractPublisher struct {
	ch *amqp.Channel
}

func NewContractPublisher(conn *amqp.Connection) (*ContractPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 1. Declare the Exchange
	err = ch.ExchangeDeclare(
		"contract_events", // Name
		"topic",           // Type
		true,              // Durable
		false,             // Auto-deleted
		false,             // Internal
		false,             // No-wait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ContractPublisher{ch: ch}, nil
}

// PublishContractPublished sends the message to the Signature Service
func (p *ContractPublisher) PublishContractPublished(event models.ContractPublishedEvent) error {
	// 2. Convert the Event struct to JSON bytes
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// 3. Publish to the Exchange with the specific Routing Key
	err = p.ch.Publish(
		"contract_events",    // Exchange
		"contract.published", // Routing Key (The Signature Consumer listens for this!)
		false,                // Mandatory
		false,                // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish contract event: %v", err)
		return err
	}

	log.Printf(" [AMQP] Published Contract: %s", event.ContractID)
	return nil
}