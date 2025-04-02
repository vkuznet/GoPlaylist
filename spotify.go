package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	spotify "github.com/zmb3/spotify/v2"
	auth "github.com/zmb3/spotify/v2/auth"
)

// helper function to setup spotify client
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
		log.Println("Spotify client successfully authenticated")

		user, err := client.CurrentUser(ctx)
		if err != nil {
			http.Error(w, "Couldn't get current user", http.StatusForbidden)
			log.Fatalf("Couldn't get current user: %v", err)
			return
		}

		// obtain orchestra either from title of discography
		orchestra := getOrchestra(title, discography)

		// check if playlist already exist, if not we will create it
		playlistID, err := getSpotifyPlaylistIDByName(client, title)
		if err != nil {
			log.Println("Unable to lookup playlist ID for", title)
			playlistID = createSpotifyPlaylist(client, user.ID, title)
		}

		// load cache entries for our playlist
		tracks, err := cache.Load("spotify", title, string(playlistID))

		// fetch existing tracks in our playlist
		//         tracks, err := getSpotifyTracksForPlaylistID(client, playlistID)
		if err != nil {
			log.Printf("unable to find tracks for playlist '%s' (%v), error %v", title, playlistID, err)
		}

		for idx, track := range discography.Tracks {
			if track.Orchestra != "" {
				orchestra = track.Orchestra
			}
			year := strings.Split(track.Year, "-")[0]
			//             query := fmt.Sprintf("track:%s year:%v artist:%s", track.Name, year, artist)
			query := fmt.Sprintf("track:%s artist:%s", track.Name, orchestra)
			if Config.Verbose > 0 {
				fmt.Printf("query idx: %4d track: %s\n", idx, query)
			}
			trk := Track{Name: track.Name, Year: year, Orchestra: orchestra, Artist: track.Artist}
			if inList(trk, tracks) {
				fmt.Printf("idx: %4d query: %s, already exist in playlist, skipping...\n", idx, query)
			} else {
				fmt.Printf("idx: %4d track: %s\n", idx, query)
				if err := addToSpotifyPlaylist(client, playlistID, query); err == nil {
					// add track to local cache if was successfully added to playlist
					cache.AddTrack(title, string(playlistID), trk)
				}
			}
		}
		purl := constructSpotifyPlaylistURL(playlistID)
		msg := fmt.Sprintf("New playlist <a href=\"%s\">%s</a> is created", purl, title)
		log.Println(msg)
		w.Header().Set("Content-Type", "text/html")
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

// helper function to create spotify playlist
func createSpotifyPlaylist(client *spotify.Client, userID, title string) spotify.ID {
	ctx := context.Background()
	playlist, err := client.CreatePlaylistForUser(ctx, userID, title, "Playlist created for Orquesta Típica", true, false)
	if err != nil {
		log.Fatalf("Error creating Spotify playlist: %v", err)
	}
	return playlist.ID
}

// helper function to add track to spotify playlist
func addToSpotifyPlaylist(client *spotify.Client, playlistID spotify.ID, trackName string) error {
	ctx := context.Background()
	searchResults, err := client.Search(ctx, trackName, spotify.SearchTypeTrack)
	if err != nil || len(searchResults.Tracks.Tracks) == 0 {
		log.Printf("Error finding track: %v", err)
		return err
	}
	trackID := searchResults.Tracks.Tracks[0].ID

	_, err = client.AddTracksToPlaylist(ctx, playlistID, trackID)
	if err != nil {
		log.Printf("Error adding track to playlist: %v", err)
		return err
	}
	return nil
}

// helper function to construct spotify playlist URL
func constructSpotifyPlaylistURL(playlistID spotify.ID) string {
	return fmt.Sprintf("https://open.spotify.com/playlist/%v", playlistID)
}

// helper function to get spotify playlist IDa by provided name
func getSpotifyPlaylistIDByName(client *spotify.Client, playlistName string) (spotify.ID, error) {
	// Fetch all playlists for the authenticated user
	playlists, err := client.CurrentUsersPlaylists(context.Background())
	if err != nil {
		return "", fmt.Errorf("error fetching user's playlists: %v", err)
	}

	// Iterate through the user's playlists to find the one matching the name
	for _, playlist := range playlists.Playlists {
		if playlist.Name == playlistName {
			return playlist.ID, nil
		}
	}

	return "", fmt.Errorf("no playlist found with name: %s", playlistName)
}

// helper function to get spotify tracks for given playlist ID
func getSpotifyTracksForPlaylistID(client *spotify.Client, playlistID spotify.ID) ([]string, error) {
	var tracks []string
	playlist, err := client.GetPlaylist(context.Background(), playlistID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving playlist: %v", err)
	}

	for _, item := range playlist.Tracks.Tracks {
		tracks = append(tracks, fmt.Sprintf("%s by %s", item.Track.Name, item.Track.Artists[0].Name))
	}

	return tracks, nil
}
