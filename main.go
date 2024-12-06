package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Add YouTube and Spotify API setup here (omitted for brevity)
var spotifyId, spotifySecret, youtubeId, youtubSecret string

// local cache
var cache *Cache

// Use the readXMLFile function
func main() {
	var file string
	flag.StringVar(&file, "file", "", "xml or csv file to read")
	var config string
	flag.StringVar(&config, "config", "", "configuration file")
	var showTracks bool
	flag.BoolVar(&showTracks, "tracks", false, "show tracks and exit")
	var sortBy string
	flag.StringVar(&sortBy, "sortBy", "", "sort tracks by attribute: orchestra, artist, year, genre, vocal")
	var sortOrder string
	flag.StringVar(&sortOrder, "sortOrder", "ascending", "sort order: ascending or descending")
	flag.Parse()

	err := parseConfig(config)
	if err != nil {
		log.Fatalf("Fail to parse config file %s, error %v", config, err)
	}
	if Config.Verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	if showTracks {
		Config.Verbose = 0
		discography, _ := readFile(file, sortBy, sortOrder)
		for _, track := range discography.Tracks {
			fmt.Printf("%+v\n", track)
		}
		return
	}

	title := Config.PlaylistTitle
	if title == "" {
		arr := strings.Split(file, "/")
		title = strings.Replace(arr[len(arr)-1], ".xml", "", -1)
	}
	fmt.Printf("creating %s playlist: %s\n", Config.Service, title)

	// read provided file
	discography, err := readFile(file, sortBy, sortOrder)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	// choose a client to use
	cache = &Cache{}
	if strings.ToLower(Config.Service) == "spotify" {
		cdir := fmt.Sprintf("%s/.goplaylist/%s", os.Getenv("HOME"), "spotify")
		cache.Init("spotify", cdir)
		setupSpotifyClient(title, discography)
	} else {
		cdir := fmt.Sprintf("%s/.goplaylist/%s", os.Getenv("HOME"), "youtube")
		cache.Init("youtube", cdir)
		setupYouTubeService(title, discography)
	}
}
