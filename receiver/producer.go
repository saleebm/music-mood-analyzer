package main

import (
	"context"
	"github.com/saleebm/music-mood-analyzer/shared"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func queue(trackId string) {
	// todo env for rabbit
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",             // exchange
		"song_updates", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(trackId),
		})
	shared.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", trackId)
}
