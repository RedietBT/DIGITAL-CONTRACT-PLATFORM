package broker

import (
	"encoding/json"
	"log"

	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/models"
	"github.com/RedietBT/DIGITAL-CONTRACT-PLATFORM/backend/internal/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SignatureConsumer struct {
	conn *amqp.Connection
	repo repository.SignatureRepository
}

func NewSignatureConsumer(conn *amqp.Connection, repo repository.SignatureRepository) (*SignatureConsumer, error) {
	return &SignatureConsumer{
		conn: conn, 
		repo: repo,
	}, nil
}

func (c *SignatureConsumer) Listen() {
	ch, err := c.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Exchange names must match what the Contract service uses to publish
	err = ch.ExchangeDeclare(
		"contract_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	q, err := ch.QueueDeclare(
		"signature_service_contract_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	// We bind to hear when a contract is published
	err = ch.QueueBind(
		q.Name,
		"contract.published",
		"contract_events",
		false,
		nil,
	)

	msgs, err := ch.Consume(
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
			var event models.ContractPublishedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Panicf("Error decoding message: %v", err)
				continue
			}
		}
	}()

	select {}
}