package main

import (
	_ "context"

	spotify "github.com/osinor/spotify-app/pkg"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

func main() {
	spotify.Do()
	// client := spotify.NewClient()
	// if err := client.GetRecommendations(context.Background()); err != nil {
	// 	log.Fatal(err)
	// }
}
