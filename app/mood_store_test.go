package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/saleebm/music-mood-analyzer/shared"
	"github.com/zmb3/spotify/v2"
	"os"
	"testing"
)

func TestMoodStoreExportSentiment(t *testing.T) {
	err := godotenv.Load("../.env")
	shared.FailOnError(err, "Failed to load env")

	fileName := "../tmp/0cOyhnhy13lc5G5nr4EF0q-happy-#FFFF00.json"
	file, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Failed to open file: %s", err.Error())
	}
	var payload spotify.AudioFeatures
	err = json.Unmarshal(file, &payload)
	if err != nil {
		t.Fatalf("Failed to read Features: %s", err.Error())
	}
	fmt.Printf("%+v\n", payload)
	moodStore := MoodStore{Features: payload, TrackId: "0cOyhnhy13lc5G5nr4EF0q", Mood: "happy", Color: "#FFFF00"}
	moodStore.ExportSentiment()
}
