package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func setupYouTubeService(title string, discography *Discography) {
	if Config.YoutubeSecret == "" {
		return
	}
	ctx := context.Background()

	// OAuth2 configuration
	config := &oauth2.Config{
		ClientID:     Config.YoutubeId,
		ClientSecret: Config.YoutubeSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: callbackUrl(),
		Scopes:      []string{youtube.YoutubeForceSslScope},
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// obtain code from HTTP request
		code := r.URL.Query().Get("code")

		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Fatalf("Unable to retrieve token from web %v", err)
		}

		client := config.Client(ctx, token)

		// Create the YouTube service
		service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Failed to create YouTube client: %v", err)
		}
		log.Println("Youtube client successfully authenticated")

		// obtain artist either from title of discography
		artist := getArtist(title, discography)

		// check if playlist already exist, if not we will create it
		playlistID, err := getYoutubePlaylistIDByName(service, title)
		if err != nil {
			log.Printf("Unable to lookup playlist ID for '%s', error %v", title, err)
			playlistID = createYoutubePlaylist(service, title)
		}

		// load cache entries for our playlist
		tracks, err := cache.Load("youtube", title, string(playlistID))

		// fetch existing tracks in our playlist
		//         tracks, err := getYoutubeTracksForPlaylistID(service, playlistID)
		if err != nil {
			log.Printf("unable to find tracks for playlist '%s' (%v), error %v", title, playlistID, err)
		}

		for idx, track := range discography.Tracks {
			if track.Orchestra != "" {
				artist = track.Orchestra
			}
			year := strings.Split(track.Year, "-")[0]
			query := fmt.Sprintf("%s %s %v", track.Name, artist, year)
			if inList(track.Name, tracks) {
				if Config.Verbose > 0 {
					fmt.Printf("idx: %d query: %s, already exist in playlist, skipping...\n", idx, query)
				}
				log.Println("")
			} else {
				if Config.Verbose > 0 {
					fmt.Printf("idx: %d track: %s\n", idx, query)
				}
				addToYoutubePlaylist(service, playlistID, query)
				// add track to local cache
				cache.AddTrack(title, playlistID, track.Name)
			}
		}
		purl := constructYouTubePlaylistURL(playlistID)
		msg := fmt.Sprintf("New playlist <a href=\"%s\">%s</a> is created", purl, title)
		log.Println(msg)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(msg))
	})

	// Start a web server to complete the auth flow
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Config.CallbackPort), nil))
	}()

	// Obtain a token (this part typically includes a user authentication process)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

}

func createYoutubePlaylist(service *youtube.Service, title string) string {
	playlist := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       title,
			Description: "Playlist",
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "public",
		},
	}

	snippet := []string{"snippet", "status"}
	createdPlaylist, err := service.Playlists.Insert(snippet, playlist).Do()
	if err != nil {
		log.Fatalf("Error creating YouTube playlist: %v", err)
	}
	return createdPlaylist.Id
}

func addToYoutubePlaylist(service *youtube.Service, playlistID, query string) {
	searchResp, err := service.Search.List([]string{"id"}).
		Q(query).
		MaxResults(1).
		Type("video").
		Do()
	if err != nil || len(searchResp.Items) == 0 {
		log.Printf("Error finding video: %v", err)
		return
	}
	videoID := searchResp.Items[0].Id.VideoId

	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}
	_, err = service.PlaylistItems.Insert([]string{"snippet"}, playlistItem).Do()
	if err != nil {
		log.Printf("Error adding video to playlist: %v", err)
	}
}

func constructYouTubePlaylistURL(playlistID string) string {
	return fmt.Sprintf("https://www.youtube.com/playlist?list=%s", playlistID)
}

func getYoutubePlaylistIDByName(service *youtube.Service, playlistName string) (string, error) {
	if Config.Verbose > 0 {
		log.Printf("lookup playlist with title: '%s' in user's account", playlistName)
	}

	// Fetch playlists owned by the authenticated user
	playlistsResp, err := service.Playlists.List([]string{"snippet"}).
		Mine(true).
		Do()
	if err != nil {
		return "", fmt.Errorf("error fetching user's playlists: %v", err)
	}

	// Iterate through the user's playlists to find the one matching the name
	for _, playlist := range playlistsResp.Items {
		if playlist.Snippet.Title == playlistName {
			if Config.Verbose > 0 {
				log.Printf("found existing playlist %s", playlist.Id)
			}
			return playlist.Id, nil
		}
	}

	return "", fmt.Errorf("no playlist found with name: %s", playlistName)
}

func getYoutubeTracksForPlaylistID(service *youtube.Service, playlistID string) ([]string, error) {
	var tracks []string
	nextPageToken := ""

	for {
		playlistItemsResp, err := service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(playlistID).
			MaxResults(50). // Maximum allowed by YouTube API
			PageToken(nextPageToken).
			Do()

		if err != nil {
			return nil, fmt.Errorf("error retrieving playlist items: %v", err)
		}

		for _, item := range playlistItemsResp.Items {
			if Config.Verbose > 0 {
				log.Printf("adding track %s to from existing playlist", item.Snippet.Title)
			}
			tracks = append(tracks, fmt.Sprintf("%s (%s)", item.Snippet.Title, item.Snippet.ResourceId.VideoId))
		}

		// Check if there's another page of results
		nextPageToken = playlistItemsResp.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return tracks, nil
}
