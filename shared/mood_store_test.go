package shared

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	"os"
	"testing"
)

func TestMoodStoreExportSentiment(t *testing.T) {
	err := godotenv.Load("../.env")
	FailOnError(err, "Failed to load env")

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
	uuidsSlice := []string{"5a6d6ee9-5d8b-4569-b9d0-563a2c850591"}[:]
	moodStore := MoodStore{Features: payload, TrackId: "20bJBbPapGQ4bqs0YcA9xY", Mood: "happy", Color: "#FFFF00", Uuids: uuidsSlice}
	moodStore.ExportSentiment()
}
