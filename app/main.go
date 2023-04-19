package main

import (
	"context"
	"github.com/saleebm/music-mood-analyzer/shared"
	"log"
	"time"
)

func main() {
	// Initialize RabbitMQ connection and SpotifyAgent client
	//todo conn string
	rabbitConnStr := "amqp://guest:guest@localhost:5672/"

	rabbitConn, err := ConnectToRabbitMQ(rabbitConnStr)
	shared.FailOnError(err, "Error connecting to RabbitMQ")
	defer rabbitConn.Close()

	limiter := time.Tick(200 * time.Millisecond)

	ctx := context.Background()
	client := NewSpotifyClient(ctx)
	tuneHandler := NewTuneHandler(limiter, client)

	err = ConsumeSongs(rabbitConn, "song_updates", tuneHandler.HandleMsg)
	if err != nil {
		log.Fatalf("Error consuming song updates: %v", err)
	}

	// Run forever
	select {}
}
