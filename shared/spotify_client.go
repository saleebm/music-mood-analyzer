package shared

import (
	"context"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
)

func NewSpotifyClient(ctx context.Context) *spotify.Client {
	spotifyId := os.Getenv("SPOTIFY_ID")
	spotifySecret := os.Getenv("SPOTIFY_SECRET")
	if len(spotifyId) == 0 || len(spotifySecret) == 0 {
		log.Panicln("missing SPOTIFY_ID or SPOTIFY_SECRET")
	}
	config := &clientcredentials.Config{
		ClientID:     spotifyId,
		ClientSecret: spotifySecret,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	FailOnError(err, "couldn't get token")

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient, spotify.WithRetry(true))
	return client
}
