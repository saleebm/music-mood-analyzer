package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/saleebm/music-mood-analyzer/shared"
	"github.com/zmb3/spotify/v2"
	"net/http"
	"os"
)

var SentimentColors = map[string]string{
	"Happy":      "#FDE74C",
	"Excited":    "#FFA07A",
	"Lively":     "#FF5733",
	"Melancholy": "#6495ED",
	"Serene":     "#ADD8E6",
	"Hopeful":    "#90EE90",
	"Fiery":      "#FF6347",
	"Anxious":    "#8B0000",
	"Brooding":   "#2F4F4F",
	"Easygoing":  "#F5DEB3",
	"Dreamy":     "#D8BFD8",
	"Peaceful":   "#F0E68C",
}

type MoodStore struct {
	Features spotify.AudioFeatures `json:"features"`
	Mood     string                `json:"mood"`
	Color    string                `json:"color"`
	TrackId  string                `json:"trackId"`
	Uuids    []string              `json:"uuids"`
}

func NewMoodStore(features *spotify.AudioFeatures, track *shared.Track) (moodStore *MoodStore) {
	moodStore = &MoodStore{Features: *features, TrackId: track.TrackId, Uuids: track.Uuids}
	return // love it
}

func (moodStore *MoodStore) GetSentiment(features *spotify.AudioFeatures) string {
	valence := features.Valence
	energy := features.Energy
	danceability := features.Danceability
	tempo := features.Tempo

	fmt.Printf("%v, %v, %v, %v\n", valence, energy, danceability, tempo)

	// Determine the sentiment based on the audio Features
	if danceability >= 0.5 && energy >= 0.5 && tempo < 100 && valence >= 0.5 {
		return "Happy"
	} else if danceability >= 0.7 && energy > 0.7 && tempo >= 100 && valence >= 0.3 {
		return "Excited"
	} else if danceability >= 0.5 && energy >= 0.7 && tempo >= 120 && valence >= 0.5 {
		return "Lively"
	} else if danceability <= 0.5 && energy <= 0.5 && tempo <= 80 && valence <= 0.5 {
		return "Melancholy"
	} else if danceability <= 0.5 && energy <= 0.5 && tempo <= 80 && valence > 0.7 {
		return "Serene"
	} else if danceability >= 0.5 && energy >= 0.7 && tempo < 120 && valence >= 0.7 {
		return "Hopeful"
	} else if danceability < 0.75 && energy >= 0.5 && tempo >= 80 && tempo < 120 && valence < 0.7 {
		return "Fiery"
	} else if danceability < 0.5 && energy >= 0.5 && tempo > 80 && tempo < 120 && valence < 0.7 {
		return "Anxious"
	} else if danceability < 0.5 && energy >= 0.5 && tempo >= 120 && valence < 0.7 {
		return "Brooding"
	} else if danceability < 0.5 && energy >= 0.5 && energy < 0.75 && tempo < 80 && valence < 0.7 {
		return "Easygoing"
	} else if danceability < 0.5 && energy <= 0.5 && tempo > 60 && tempo < 80 && valence < 0.7 {
		return "Dreamy"
	} else if danceability < 0.5 && energy <= 0.5 && tempo <= 80 && valence <= 0.7 {
		return "Peaceful"
	} else {
		fmt.Printf("Missing Mood for\n\tValence: %v,\t\nEnergy: %v,\t\nTempo: %v,\t\nDanceability: %v\n", valence, energy, tempo, danceability)
		return ""
	}
}

func (moodStore *MoodStore) ExportSentiment() {
	endpointUrl := os.Getenv("ENDPOINT_URL")
	endpointSecret := os.Getenv("ENDPOINT_SECRET")
	if len(endpointSecret) == 0 || len(endpointUrl) == 0 {
		fmt.Println("Missing env ENDPOINT_SECRET or ENDPOINT_URL")
		return
	}
	w := new(bytes.Buffer)
	err := json.NewEncoder(w).Encode(&moodStore)
	if err != nil {
		fmt.Println("Unable to encode Mood store", err)
		return
	}
	//fmt.Printf("%+v\n", w)
	req, err := http.NewRequest("POST", endpointUrl, w)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret", endpointSecret)
	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request")
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Sent request to save sentiment. Status code %d \n", resp.StatusCode)
}

func (moodStore *MoodStore) WriteResults() {
	mood := moodStore.Mood
	color := moodStore.Color
	features := moodStore.Features
	env := os.Getenv("GO_ENV")
	tmpDir := os.TempDir()
	if env == "development" {
		// Set tmp directory to the path of the Go folder
		goPath := os.Getenv("GOPATH")
		tmpDir = goPath + "/saleebm/music-Mood-analyzer/tmp"
	}
	filename := fmt.Sprintf("%s/%s-%s-%s.json", tmpDir, moodStore.TrackId, mood, color)
	exists, err := shared.Exists(filename)
	if exists {
		fmt.Println("File " + filename + " already exists")
		return
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

	fmt.Printf("%+v saved to %+v\n", moodStore.TrackId, filename)
}
