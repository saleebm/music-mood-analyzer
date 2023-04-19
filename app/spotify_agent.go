package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zmb3/spotify/v2"
	//"math/rand"
	"os"
	//"time"
)

type SpotifyAgent struct {
	trackId spotify.ID
	client  *spotify.Client
}

func NewSpotifyAgent(trackId string, client *spotify.Client) *SpotifyAgent {
	return &SpotifyAgent{trackId: spotify.ID(trackId), client: client}
}

// GetAudioAnalysis gets the audio analysis from spotify api
func (spotifyAgent SpotifyAgent) GetAudioAnalysis() (*spotify.AudioAnalysis, error) {
	analysis, err := spotifyAgent.client.GetAudioAnalysis(context.Background(), spotifyAgent.trackId)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

// GetAudioFeatures get the features of the track
func (spotifyAgent SpotifyAgent) GetAudioFeatures() (*spotify.AudioFeatures, error) {
	features, err := spotifyAgent.client.GetAudioFeatures(context.Background(), spotifyAgent.trackId)
	if err != nil {
		return nil, err
	}
	return features[0], nil
}

// ProcessAudioAnalysis extract mood information from the audio analysis.
func (spotifyAgent SpotifyAgent) ProcessAudioAnalysis(features *spotify.AudioFeatures) (moodStore *MoodStore) {
	// Initialize the mood store
	moodStore = NewMoodStore(features)
	spotifyAgent.WriteResults(moodStore) // log results
	return
}

// WriteResults Keeps a record of the analysis locally
func (spotifyAgent SpotifyAgent) WriteResults(moodStore *MoodStore) {
	mood := moodStore.mood
	color := moodStore.color
	features := moodStore.features
	env := os.Getenv("GO_ENV")
	tmpDir := os.TempDir()
	if env == "development" {
		// Set tmp directory to the path of the Go folder
		goPath := os.Getenv("GOPATH")
		tmpDir = goPath + "/tmp"
	}
	filename := fmt.Sprintf("%s%s-%s-%s.json", tmpDir, spotifyAgent.trackId, mood, color)
	fmt.Printf("%+v saved to %+v\n", spotifyAgent.trackId, filename)

	_, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0755)
	if errors.Is(err, os.ErrNotExist) {
		//todo
	}

	jsonBytes, err := json.Marshal(features)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Write the JSON to a file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("JSON saved to file:", filename)
}
