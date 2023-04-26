package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/saleebm/music-mood-analyzer/shared"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Missing .env, %s", err.Error())
	}
	// Initialize RabbitMQ connection and SpotifyAgent client
	rabbitConnStr := os.Getenv("RABBITMQ_CONN_STR")
	if len(rabbitConnStr) == 0 {
		log.Fatalf("Missing rabbit mq conn str")
	}

	rabbitConn, err := ConnectToRabbitMQ(rabbitConnStr)
	shared.FailOnError(err, "Error connecting to RabbitMQ")
	defer rabbitConn.Close()

	limiter := time.Tick(1000 * time.Millisecond) // process one every second

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
