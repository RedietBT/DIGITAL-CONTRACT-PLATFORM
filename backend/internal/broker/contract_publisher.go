package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	amqp "github.com/rabbitmq/amqp091-go" // Standard modern library
)

type ContractPublisher struct {
	ch *amqp.Channel
}

func NewContractPublisher(conn *amqp.Connection) (*ContractPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 1. Declare the Exchange (Topic is great for microservices)
	err = ch.ExchangeDeclare(
		"contract_events", 
		"topic",           
		true,              
		false,             
		false,             
		false,             
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ContractPublisher{ch: ch}, nil
}

// PublishSignatureCreated notifies other services that a new signature was added
func (p *ContractPublisher) PublishSignatureCreated(ctx context.Context, event models.SignatureCreatedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Use PublishWithContext for modern RabbitMQ support
	err = p.ch.PublishWithContext(
		ctx,
		"contract_events",   // Exchange
		"signature.created", // Routing Key (Contract Service will listen for this)
		false,               
		false,               
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Printf("Failed to publish signature event: %v", err)
		return err
	}

	log.Printf(" [AMQP] Published Signature: %s for Contract: %s", event.SignatureID, event.ContractID)
	return nil
}

// (Optional) You can keep this if the Signature Service ever needs to republish contract data
func (p *ContractPublisher) PublishContractPublished(ctx context.Context, event models.ContractPublishedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.ch.PublishWithContext(
		ctx,
		"contract_events",
		"contract.published",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}