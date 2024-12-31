package spotify

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/spotify"
)

type SpotifyAuth struct {
}

func NewSpotifyAuth() *SpotifyAuth {
	return &SpotifyAuth{}
}

func (s *SpotifyAuth) GetToken(ctx context.Context) (*oauth2.Token, error) {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.Endpoint.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}
