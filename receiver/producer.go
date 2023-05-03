package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/saleebm/music-mood-analyzer/shared"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func queue(track *shared.Track) {
	// Initialize RabbitMQ connection and SpotifyAgent client
	rabbitConnStr := os.Getenv("RABBITMQ_CONN_STR")
	if len(rabbitConnStr) == 0 {
		log.Fatalf("Missing rabbit mq conn str")
	}

	conn, err := amqp.Dial(rabbitConnStr)
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("track: %+v\n", track)
	body, err := json.Marshal(track)
	if err != nil {
		fmt.Println("Unable to marshall track", err)
		return
	}

	err = ch.PublishWithContext(ctx,
		"",             // exchange
		"song_updates", // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
			ContentEncoding: "",
		})
	shared.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", track.TrackId)
}
