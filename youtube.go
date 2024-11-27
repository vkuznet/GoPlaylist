package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	parts := strings.Split(part, ",")
	call := service.Channels.List(parts)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func youtubeChannels() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
}

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
		log.Printf("Youtube client successfully authenticated: %+v", service)

		artist := getArtist(title, discography)
		playlistID := createPlaylist(service, title)
		for _, track := range discography.Tracks {
			//             year := strings.Split(track.Year, "-")[0]
			//             query := fmt.Sprintf("track:%s year:%v artist:%s", track.Name, year, artist)
			query := fmt.Sprintf("track:%s artist:%s", track.Name, artist)
			if Config.Verbose > 0 {
				fmt.Println("searching for", query)
			}
			addToPlaylist(service, playlistID, track.Name)
		}
		msg := fmt.Sprintf("New playlist \"%s\" is created with ID: %s", title, playlistID)
		log.Println(msg)
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

func createPlaylist(service *youtube.Service, title string) string {
	playlist := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       title,
			Description: "Playlist",
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "public",
		},
	}

	snippet := []string{"snippet.status"}
	createdPlaylist, err := service.Playlists.Insert(snippet, playlist).Do()
	if err != nil {
		log.Fatalf("Error creating YouTube playlist: %v", err)
	}
	return createdPlaylist.Id
}

func addToPlaylist(service *youtube.Service, playlistID, trackName string) {
	searchResp, err := service.Search.List([]string{"id"}).
		Q(trackName).
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
