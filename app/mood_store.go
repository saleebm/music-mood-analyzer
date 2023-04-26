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
	"happy":      "#FFFF00", // yellow
	"energetic":  "#FFA500", // orange
	"excited":    "#FF0000", // red
	"sad":        "#0000FF", // blue
	"calm":       "#008000", // green
	"mysterious": "#800080", // purple
	"romantic":   "#FFC0CB", // pink
	"nostalgic":  "#964B00", // brown
	"hopeful":    "#ADD8E6", // light blue
	"angry":      "#8B0000", // dark red
	"anxious":    "#808080", // gray
	"confused":   "#40E0D0", // turquoise
	"powerful":   "#00008B", // dark blue
	"playful":    "#90EE90", // light green
	"dreamy":     "#E6E6FA", // lavender
	"peaceful":   "#D8BFD8", // light purple
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
	if valence >= 0.5 && energy >= 0.5 && danceability >= 0.5 && tempo >= 120 {
		return "happy"
	} else if valence <= 0.5 && energy >= 0.5 && danceability >= 0.5 && tempo >= 120 {
		return "excited"
	} else if valence <= 0.5 && energy >= 0.5 && danceability >= 0.5 {
		return "energetic"
	} else if valence <= 0.5 && energy <= 0.5 && danceability >= 0.5 && tempo <= 100 {
		return "sad"
	} else if valence >= 0.5 && energy <= 0.5 && danceability <= 0.5 && tempo <= 100 {
		return "calm"
	} else if valence <= 0.5 && energy <= 0.5 && danceability <= 0.5 && tempo <= 100 {
		return "mysterious"
	} else if valence >= 0.5 && energy <= 0.5 && danceability <= 0.5 && tempo >= 100 && tempo <= 120 {
		return "romantic"
	} else if valence <= 0.5 && energy <= 0.5 && danceability <= 0.5 && tempo > 100 {
		return "nostalgic"
	} else if valence >= 0.5 && energy >= 0.5 && danceability >= 0.5 && tempo >= 100 && tempo <= 120 {
		return "hopeful"
	} else if valence <= 0.5 && energy <= 0.5 && danceability >= 0.5 && tempo >= 120 {
		return "angry"
	} else if valence <= 0.5 && energy <= 0.5 && danceability >= 0.5 && tempo <= 120 {
		return "anxious"
	} else if valence >= 0.4 && valence <= 0.6 && energy >= 0.4 && energy <= 0.6 && danceability >= 0.5 && danceability <= 0.7 && tempo >= 100 &&
		tempo <= 120 {
		return "confused"
	} else if energy >= 0.7 && tempo >= 120 && valence <= 0.5 && danceability <= 0.5 {
		return "powerful"
	} else if valence >= 0.7 && danceability >= 0.7 && energy >= 0.5 && tempo >= 100 && tempo <= 120 {
		return "playful"
	} else if valence >= 0.7 && energy <= 0.3 && danceability >= 0.5 && tempo <= 100 {
		return "dreamy"
	} else if valence >= 0.7 && energy <= 0.3 && danceability <= 0.5 && tempo <= 100 {
		return "peaceful"
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
	fmt.Printf("%+v\n", w)
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
		tmpDir = goPath + "/saleebm/music-Mood-analyzer/tmp/"
	}
	filename := fmt.Sprintf("%s%s-%s-%s.json", tmpDir, moodStore.TrackId, mood, color)
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
