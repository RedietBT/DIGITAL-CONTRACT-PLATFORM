package broker

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connect returns a raw RabbitMQ connection.
// Services will pass this connectiom into their specific Publishers/Consumers.
func Connect(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	return conn, nil
}