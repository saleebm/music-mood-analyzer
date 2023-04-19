package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

type Track struct {
	TrackId string `json:"trackId"`
}

func main() {
	http.HandleFunc("/track", handleTrack)
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Printf("Error: %v\nStack trace:\n%s", err, debug.Stack())
		return
	}
}

func handleTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var track Track
	err := json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	fmt.Fprintf(w, "Received track ID: %s", track.TrackId)
	if len(track.TrackId) > 0 {
		queue(track.TrackId)
	}
}
