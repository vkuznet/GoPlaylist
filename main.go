package main

import (
	"encoding/json"
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
	var title string
	flag.StringVar(&title, "title", "", "title of new playlist")
	var showTracks bool
	flag.BoolVar(&showTracks, "tracks", false, "show tracks and exit")
	var sortBy string
	flag.StringVar(&sortBy, "sortBy", "", "sort tracks by attribute: orchestra, artist, year, genre, vocal")
	var filterBy string
	flag.StringVar(&filterBy, "filterBy", "", "filter tracks conditions: {key:value, key:value}")
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

	var filters map[string]string
	if filterBy != "" {
		err := json.Unmarshal([]byte(filterBy), &filters)
		if err != nil {
			log.Fatal(err)
		}
	}

	// by default use playlist title
	ptitle := Config.PlaylistTitle
	if ptitle == "" {
		// if it is not parsed from input file we'll use name of the file itself
		arr := strings.Split(file, "/")
		if strings.HasSuffix(file, ".xml") {
			ptitle = strings.Replace(arr[len(arr)-1], ".xml", "", -1)
		} else if strings.HasSuffix(file, ".csv") {
			ptitle = strings.Replace(arr[len(arr)-1], ".csv", "", -1)
		}
	}
	if title != "" {
		// use title provided via option
		ptitle = title
	}
	if ptitle == "" {
		log.Fatal("empty playlist title")
	}
	fmt.Printf("creating %s playlist: %s\n", Config.Service, ptitle)

	// if asked for tracks only, display them and exit
	if showTracks {
		Config.Verbose = 0
		discography, _ := readFile(file, sortBy, sortOrder, filters)
		for _, track := range discography.Tracks {
			fmt.Printf("%+v\n", track)
		}
		return
	}

	// read provided file
	discography, err := readFile(file, sortBy, sortOrder, filters)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	// choose a client to use
	cache = &Cache{}
	if strings.ToLower(Config.Service) == "spotify" {
		cdir := fmt.Sprintf("%s/.goplaylist/%s", os.Getenv("HOME"), "spotify")
		cache.Init("spotify", cdir)
		setupSpotifyClient(ptitle, discography)
	} else {
		cdir := fmt.Sprintf("%s/.goplaylist/%s", os.Getenv("HOME"), "youtube")
		cache.Init("youtube", cdir)
		setupYouTubeService(ptitle, discography)
	}
}
