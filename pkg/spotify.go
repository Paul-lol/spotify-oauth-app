package spotify

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify"
	rspotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Authenticator interface {
	GetToken(ctx context.Context) (*oauth2.Token, error)
}

type Client struct {
	auth Authenticator
}

func NewClient() *Client {
	return &Client{
		auth: NewSpotifyAuth(),
	}
}

func (c *Client) GetRecommendations(ctx context.Context) error {
	token, err := c.auth.GetToken(ctx)
	if err != nil {
		return err
	}
	client := spotifyauth.New().Client(ctx, token)
	spotClient := rspotify.New(client)

	ids, err := c.getArtistIds(ctx, spotClient)
	if err != nil {
		return err
	}

	recs, err := spotClient.GetRecommendations(ctx, rspotify.Seeds{Artists: ids}, rspotify.NewTrackAttributes())
	if err != nil {
		return err
	}

	fmt.Println(recs)
	return nil
}

func (c *Client) getArtistIds(ctx context.Context, client *rspotify.Client) ([]rspotify.ID, error) {
	result, err := client.Search(ctx, "Asake", spotify.SearchTypeArtist)
	if err != nil {
		return nil, err
	}

	return []rspotify.ID{result.Artists.Artists[0].ID}, nil
}
