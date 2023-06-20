package main

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/saleebm/music-mood-analyzer/shared"
	spotify "github.com/zmb3/spotify/v2"
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
	var track *shared.Track
	err := json.Unmarshal(msg.Body, &track)
	if err != nil {
		log.Printf("Unable to unmarshall track\n Error %s\n", err.Error())
		return
	}
	moodStore, err := tuneHandler.processTrack(track)
	fmt.Printf("%+v\n", moodStore)
	shared.FailOnError(err, "Failed to process track")
}

func (tuneHandler *TuneHandler) processTrack(track *shared.Track) (*MoodStore, error) {
	<-tuneHandler.limiter
	fmt.Printf("track %+v", track)
	trackId := track.TrackId
	fmt.Println("--------\nRequest", time.Now().Format("2006-01-02 3:4:5 pm"))
	fmt.Printf("Message: %s\n", trackId)
	spotifyAgent := NewSpotifyAgent(trackId, tuneHandler.client)

	features, err := spotifyAgent.GetAudioFeatures()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting audio Features for %s: %v", trackId, err))
	}

	moodStore := NewMoodStore(features, track)
	moodStore.GetSentiment(features)
	moodStore.ExportSentiment() // export back to node
	moodStore.WriteResults()

	fmt.Printf("Mood for %s: %s\nColor: %s\n", trackId, moodStore.Mood, moodStore.Color)
	return moodStore, nil
}
