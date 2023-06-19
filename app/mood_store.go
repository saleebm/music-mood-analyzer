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
	"Angry":        "#bf0000",
	"Fiery":        "#f74959",
	"Lively":       "#ff5733",
	"Nervous":      "#f9a044",
	"Anxious":      "#f05404",
	"Worried":      "#e1a90b",
	"Concerned":    "#ea9708",
	"Confused":     "#f7bf73",
	"Afraid":       "#f7db16",
	"Surprised":    "#f8d23b",
	"Peaceful":     "#f0e68c",
	"Excited":      "#c8e719",
	"Amazed":       "#00a900",
	"Happy":        "#fde74c",
	"Hopeful":      "#90ee90",
	"Easygoing":    "#f5deb3",
	"Brooding":     "#2f4f4f",
	"Grief":        "#0070ab",
	"Melancholy":   "#6495ed",
	"Serene":       "#add8e6",
	"Disappointed": "#0033c8",
	"Upset":        "#4600c2",
	"Disdain":      "#6800a1",
	"Horrify":      "#cd21c3",
	"Mysterious":   "#ce0078",
	"Cranky":       "#f54587",
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
	if danceability < 0.5 && energy >= 0.75 && tempo < 80 && valence < 0.3 {
		return "Angry"
	} else if danceability >= 0.7 && energy >= 0.8 && tempo > 60 && tempo < 120 && valence >= 0.7 {
		return "Fiery"
	} else if danceability >= 0.7 && energy >= 0.5 && energy < 0.8 && tempo >= 120 && valence >= 0.4 && valence < 0.7 {
		return "Lively"
	} else if danceability <= 0.4 && energy < 0.6 && tempo >= 100 && valence <= 0.5 {
		return "Nervous"
	} else if danceability <= 0.6 && energy >= 0.6 && tempo >= 80 && valence > 0.2 && valence <= 0.5 {
		return "Anxious"
	} else if danceability <= 0.4 && energy <= 0.5 && tempo > 80 && valence <= 0.4 {
		return "Worried"
	} else if danceability <= 0.5 && energy < 0.6 && tempo < 120 && valence > 0.4 && valence < 0.7 {
		return "Concerned"
	} else if danceability > 0.5 && energy <= 0.5 && tempo <= 100 && valence > 0.4 {
		return "Confused"
	} else if danceability < 0.6 && energy < 0.6 && tempo < 120 && valence < 0.3 {
		return "Afraid"
	} else if danceability < 0.5 && energy <= 0.4 && tempo <= 90 && valence >= 0.7 {
		return "Peaceful"
	} else if danceability >= 0.7 && energy >= 0.8 && tempo >= 120 && valence >= 0.7 {
		return "Excited"
	} else if danceability >= 0.3 && danceability < 0.7 && energy >= 0.6 && tempo > 120 && valence >= 0.7 {
		return "Amazed"
	} else if danceability >= 0.7 && energy >= 0.7 && tempo < 120 && valence >= 0.7 {
		return "Happy"
	} else if danceability >= 0.5 && energy >= 0.3 && energy < 0.7 && tempo < 120 && valence >= 0.7 {
		return "Hopeful"
	} else if danceability >= 0.5 && energy < 0.5 && tempo < 120 && valence >= 0.5 {
		return "Easygoing"
	} else if danceability > 0.2 && danceability <= 0.5 && energy <= 0.5 && tempo < 80 && valence <= 0.3 {
		return "Brooding"
	} else if danceability <= 0.4 && energy <= 0.4 && tempo > 90 && valence <= 0.3 {
		return "Grief"
	} else if danceability > 0.4 && energy <= 0.4 && tempo <= 90 && valence <= 0.3 {
		return "Melancholy"
	} else if danceability < 0.5 && energy <= 0.4 && tempo <= 90 && valence > 0.3 && valence < 0.7 {
		return "Serene"
	} else if danceability <= 0.5 && energy <= 0.5 && tempo < 80 && valence <= 0.4 {
		return "Disappointed"
	} else if danceability <= 0.2 && energy < 0.6 && tempo <= 80 && valence <= 0.3 {
		return "Upset"
	} else if danceability <= 0.3 && energy >= 0.6 && tempo < 100 && valence <= 0.4 {
		return "Disdain"
	} else if danceability <= 0.3 && energy >= 0.7 && tempo >= 100 && valence <= 0.2 {
		return "Horrify"
	} else if danceability > 0.3 && danceability <= 0.6 && energy < 0.6 && tempo < 100 && valence <= 0.5 {
		return "Mysterious"
	} else if danceability <= 0.3 && energy <= 0.5 && tempo > 80 && valence <= 0.3 {
		return "Cranky"
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
