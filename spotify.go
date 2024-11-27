package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spotify "github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"
)

func setupSpotifyClient(title string, discography *Discography) {
	ctx := context.Background()

	// Create a new authenticator for Spotify
	auth := auth.New(
		auth.WithClientID(Config.SpotifyId),
		auth.WithClientSecret(Config.SpotifySecret),
		auth.WithRedirectURL(callbackUrl()),
		auth.WithScopes(auth.ScopePlaylistModifyPublic),
	)

	// Handle token via a callback URL and redirect
	state := "random-state-string"
	var client *spotify.Client
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(ctx, state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			log.Fatalf("Couldn't get token: %v", err)
			return
		}

		client = spotify.New(auth.Client(ctx, token))
		log.Printf("Spotify client successfully authenticated: %+v", client)

		user, err := client.CurrentUser(ctx)
		if err != nil {
			http.Error(w, "Couldn't get current user", http.StatusForbidden)
			log.Fatalf("Couldn't get current user: %v", err)
			return
		}
		spotifyID := createSpotifyPlaylist(client, user.ID, title)
		artist := getArtist(title, discography)
		for _, track := range discography.Tracks {
			//             year := strings.Split(track.Year, "-")[0]
			//             query := fmt.Sprintf("track:%s year:%v artist:%s", track.Name, year, artist)
			query := fmt.Sprintf("track:%s artist:%s", track.Name, artist)
			if Config.Verbose > 0 {
				fmt.Println("searching for", query)
			}
			addToSpotifyPlaylist(client, spotifyID, query)
		}
		msg := fmt.Sprintf("New playlist \"%s\" is created. Spotify ID: %s", title, spotifyID)
		log.Println(msg)
		w.Write([]byte(msg))
	})

	// Start a web server to complete the auth flow
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.CallbackPort), nil))
	}()

	// Redirect user to Spotify's auth page
	authURL := auth.AuthURL(state)
	log.Printf("Please log in to Spotify by visiting the following page in your browser:\n%s", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
}

func createSpotifyPlaylist(client *spotify.Client, userID, title string) spotify.ID {
	ctx := context.Background()
	playlist, err := client.CreatePlaylistForUser(ctx, userID, title, "Playlist created for Orquesta Típica", true, false)
	if err != nil {
		log.Fatalf("Error creating Spotify playlist: %v", err)
	}
	return playlist.ID
}

func addToSpotifyPlaylist(client *spotify.Client, playlistID spotify.ID, trackName string) {
	ctx := context.Background()
	searchResults, err := client.Search(ctx, trackName, spotify.SearchTypeTrack)
	if err != nil || len(searchResults.Tracks.Tracks) == 0 {
		log.Printf("Error finding track: %v", err)
		return
	}
	trackID := searchResults.Tracks.Tracks[0].ID

	_, err = client.AddTracksToPlaylist(ctx, playlistID, trackID)
	if err != nil {
		log.Printf("Error adding track to playlist: %v", err)
	}
}
