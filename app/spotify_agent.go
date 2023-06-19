package main

import (
	"context"
	spotify "github.com/zmb3/spotify/v2"
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
// todo unused
func (spotifyAgent *SpotifyAgent) GetAudioAnalysis() (*spotify.AudioAnalysis, error) {
	analysis, err := spotifyAgent.client.GetAudioAnalysis(context.Background(), spotifyAgent.trackId)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

// GetAudioFeatures get the Features of the track
func (spotifyAgent *SpotifyAgent) GetAudioFeatures() (*spotify.AudioFeatures, error) {
	features, err := spotifyAgent.client.GetAudioFeatures(context.Background(), spotifyAgent.trackId)
	if err != nil {
		return nil, err
	}
	return features[0], nil
}
