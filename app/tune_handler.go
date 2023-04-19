package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zmb3/spotify/v2"
	"log"
	"time"
)

type TuneHandler struct {
	limiter <-chan time.Time
	client  *spotify.Client
}

func NewTuneHandler(limiter <-chan time.Time, client *spotify.Client) *TuneHandler {
	return &TuneHandler{limiter: limiter, client: client}
}

func (tuneHandler *TuneHandler) HandleMsg(msg amqp.Delivery) {
	reqTrackId := string(msg.Body)
	tuneHandler.processTrack(reqTrackId)
}

func (tuneHandler *TuneHandler) processTrack(reqTrackId string) {
	<-tuneHandler.limiter
	fmt.Println("request", time.Now())
	fmt.Printf("Message: %s\n", reqTrackId)
	spotifyAgent := NewSpotifyAgent(reqTrackId, tuneHandler.client)

	features, err := spotifyAgent.GetAudioFeatures()
	if err != nil {
		log.Printf("Error getting audio features for %s: %v", reqTrackId, err)
		return
	}

	moodStore := spotifyAgent.ProcessAudioAnalysis(features)
	if err != nil {
		log.Printf("Error processing audio analysis for %s: %v", reqTrackId, err)
		return
	}

	fmt.Printf("\tMood for %s: %s\n\tColor: %s\n", reqTrackId, moodStore.mood, moodStore.color)
}
