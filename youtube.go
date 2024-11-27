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

		artist := getArtist(title, discography)
		playlistID := createYoutubePlaylist(service, title)
		for _, track := range discography.Tracks {
			if track.Orchestra != "" {
				artist = track.Orchestra
			}
			year := strings.Split(track.Year, "-")[0]
			query := fmt.Sprintf("%s %s %v", track.Name, artist, year)
			if Config.Verbose > 0 {
				fmt.Println("searching for", query)
			}
			addToYoutubePlaylist(service, playlistID, query)
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
