package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/saleebm/music-mood-analyzer/shared"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Missing .env, %s", err.Error())
	}

	http.HandleFunc("/track", handleTrack)
	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Printf("Error: %v\nStack trace:\n%s", err, debug.Stack())
		return
	}
}

// handleTrack get a track to queue into sentiment analysis processor
func handleTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("wrong method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var track *shared.Track
	//todo fix this and test using a test response writer
	err := json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "Error decoding JSON: %v", err)
		if err != nil {
			return
		}
		return
	}

	_, err = fmt.Fprintf(w, "Received track ID: %s\nUUIDs for track: %+v", track.TrackId, track.Uuids)
	if err != nil {
		return
	}
	if len(track.TrackId) > 0 {
		queue(track)
	}
}
