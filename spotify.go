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
		log.Println("Spotify client successfully authenticated")

		user, err := client.CurrentUser(ctx)
		if err != nil {
			http.Error(w, "Couldn't get current user", http.StatusForbidden)
			log.Fatalf("Couldn't get current user: %v", err)
			return
		}

		// obtain artist either from title of discography
		artist := getArtist(title, discography)

		// check if playlist already exist, if not we will create it
		spotifyID, err := getSpotifyPlaylistIDByName(client, title)
		if err != nil {
			log.Println("Unable to lookup playlist ID for", title)
			spotifyID = createSpotifyPlaylist(client, user.ID, title)
		}

		// fetch existing tracks in our playlist
		tracks, err := getSpotifyTracksForPlaylistID(client, spotifyID)

		for idx, track := range discography.Tracks {
			if track.Orchestra != "" {
				artist = track.Orchestra
			}
			//             year := strings.Split(track.Year, "-")[0]
			//             query := fmt.Sprintf("track:%s year:%v artist:%s", track.Name, year, artist)
			query := fmt.Sprintf("track:%s artist:%s", track.Name, artist)
			if Config.Verbose > 0 {
				fmt.Printf("query idx: %d track: %s\n", idx, query)
			}
			if inList(track.Name, tracks) {
				if Config.Verbose > 0 {
					fmt.Printf("idx: %d track: %s, already exist in playlist, skipping...\n", idx, query)
				}
				log.Println("")
			} else {
				if Config.Verbose > 0 {
					fmt.Printf("idx: %d track: %s\n", idx, query)
				}
				addToSpotifyPlaylist(client, spotifyID, query)
			}
		}
		purl := constructSpotifyPlaylistURL(spotifyID)
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

func constructSpotifyPlaylistURL(spotifyID spotify.ID) string {
	return fmt.Sprintf("https://open.spotify.com/playlist/%v", spotifyID)
}

func getSpotifyPlaylistIDByName(client *spotify.Client, playlistName string) (spotify.ID, error) {
	searchResults, err := client.Search(context.Background(), playlistName, spotify.SearchTypePlaylist)
	if err != nil {
		return "", fmt.Errorf("error searching playlists: %v", err)
	}

	if searchResults.Playlists == nil || len(searchResults.Playlists.Playlists) == 0 {
		return "", fmt.Errorf("no playlists found for name: %s", playlistName)
	}

	// Return the ID of the first playlist in the search results
	return searchResults.Playlists.Playlists[0].ID, nil
}

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
