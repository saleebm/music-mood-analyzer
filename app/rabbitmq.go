package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(connectionString string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

func ConsumeSongs(conn *amqp.Connection, queueName string, handleMessage func(amqp.Delivery)) error {
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %w", err)
	}

	go func() {
		for msg := range messages {
			handleMessage(msg) // todo consider requeue back if failed
			msg.Ack(false)
		}
	}()

	return nil
}
