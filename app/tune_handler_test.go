package main

import (
	"context"
	"testing"
	"time"
)

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTuneHandler(t *testing.T) {
	t.Run("It gets the real thing", func(t *testing.T) {
		ctx := context.Background()
		client := NewSpotifyClient(ctx)
		trackId := "4eSGSqP2TZvvX0kadZZttM"
		limiter := time.Tick(200 * time.Millisecond)
		tuneHandler := NewTuneHandler(limiter, client)

		tuneHandler.processTrack(trackId)
		// todo add assertions here to check the result
	})
	t.Run("multiple requests rate limit", func(t *testing.T) {
		ctx := context.Background()
		client := NewSpotifyClient(ctx)
		trackId := "4eSGSqP2TZvvX0kadZZttM"
		limiter := time.Tick(200 * time.Millisecond)
		tuneHandler := NewTuneHandler(limiter, client)

		for i := 0; i < 10; i++ {
			tuneHandler.processTrack(trackId)
		}
		// todo add assertions here to check the result
	})
}
