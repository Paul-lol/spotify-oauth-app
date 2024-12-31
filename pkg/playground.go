package spotify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	privUser = &spotify.PrivateUser{}
	auth     = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	// ch       = make(chan *spotify.Client)
	state = "abc123"
)

func Do() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/output", printUserInfo)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	//	go func() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	//	}()

	// wait for auth to complete
	// client := <-ch
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	httpClient := spotifyauth.New().Client(r.Context(), token)
	client := spotify.New(httpClient)

	// search for albums with the name Sempiternal
	results, err := client.Search(r.Context(), "Playboy", spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// select the top album
	item := results.Albums.Albums[0]
	fmt.Println(item)

	// get the current user info
	privUser, err = client.CurrentUser(r.Context())
	if err != nil {
		log.Fatalf("error getting user: %s", err)
	}

	return
	// http.Redirect()
	// json.NewEncoder(w).Encode(me)
	// get tracks from album
	res, err := client.GetAlbumTracks(r.Context(), item.ID, spotify.Market("US"))

	if err != nil {
		log.Fatal("error getting tracks ....", err.Error())
		return
	}

	// *display in tabular form using TabWriter
	wr := tabwriter.NewWriter(os.Stdout, 10, 2, 3, ' ', 0)
	fmt.Fprintf(wr, "%s\t%s\t%s\t%s\t\n\n", "Songs", "Energy", "Danceability", "Valence")

	// loop through tracks
	for _, track := range res.Tracks {

		// retrieve features
		features, err := client.GetAudioFeatures(r.Context(), track.ID)
		if err != nil {
			log.Fatal("error getting audio features...", err.Error())
			return
		}
		fmt.Fprintf(w, "%s\t%v\t%v\t%v\t\n", track.Name, features[0].Energy, features[0].Danceability, features[0].Valence)
		wr.Flush()
	}
}

func printUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusContinue)
	json.NewEncoder(w).Encode(privUser)
}
