package main

import (
	"fmt"
	"github.com/zmb3/spotify/v2"
)

var sentimentColors = map[string]string{
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
	features spotify.AudioFeatures
	mood     string
	color    string
}

func (moodStore *MoodStore) getSentiment(features *spotify.AudioFeatures) string {
	valence := features.Valence
	energy := features.Energy
	danceability := features.Danceability
	tempo := features.Tempo

	fmt.Printf("%v, %v, %v, %v\n", valence, energy, danceability, tempo)

	// Determine the sentiment based on the audio features
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
		fmt.Printf("Missing mood for\n\tValence: %v,\nEnergy: %v,\nTempo: %v,\nDanceability: %v\n", valence, energy, tempo, danceability)
		return ""
	}
}

func NewMoodStore(features *spotify.AudioFeatures) (moodStore *MoodStore) {
	moodStore = &MoodStore{features: *features}
	mood := moodStore.getSentiment(features)
	color := sentimentColors[mood]
	if len(mood) > 0 && len(color) == 0 {
		fmt.Printf("Missing color for mood, %s\n", mood)
	}
	moodStore.mood = mood
	moodStore.color = color
	return
}
