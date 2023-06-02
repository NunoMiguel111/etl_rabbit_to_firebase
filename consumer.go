package main

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func consume_from_rabbit(channel *amqp091.Channel, queue_name string) (<-chan amqp091.Delivery, error) {

	// Consume all messages from
	messages, err := channel.Consume(
		queue_name, // Queue name
		"",         // Consumer tag
		true,       // Auto-acknowledge
		false,      // Exclusive
		false,      // No-local
		false,      // No wait
		nil,        // Arguments
	)

	if err != nil {
		return nil, fmt.Errorf("Error consuming messages: %v", err)
	}

	return messages, err
}
