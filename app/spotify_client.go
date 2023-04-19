package main

import (
	"context"
	"github.com/saleebm/music-mood-analyzer/shared"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"os"
)

func NewSpotifyClient(ctx context.Context) *spotify.Client {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	shared.FailOnError(err, "couldn't get token")

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	return client
}
