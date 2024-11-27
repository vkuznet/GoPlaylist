package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

// Add YouTube and Spotify API setup here (omitted for brevity)
var spotifyId, spotifySecret, youtubeId, youtubSecret string

// Use the readXMLFile function
func main() {
	var file string
	flag.StringVar(&file, "file", "", "xml file")
	var config string
	flag.StringVar(&config, "config", "", "configuration file")
	flag.Parse()

	err := parseConfig(config)
	if err != nil {
		log.Fatalf("Fail to parse config file %s, error %v", config, err)
	}

	title := Config.PlaylistTitle
	if title == "" {
		arr := strings.Split(file, "/")
		title = strings.Replace(arr[len(arr)-1], ".xml", "", -1)
	}
	fmt.Printf("creating %s playlist: %s\n", Config.Service, title)

	// read provided file
	discography, err := readXMLFile(file)
	if err != nil {
		log.Fatalf("Error reading XML file: %v", err)
	}

	// choose a client to use
	if strings.ToLower(Config.Service) == "spotify" {
		setupSpotifyClient(title, discography)
	} else {
		setupYouTubeService(title, discography)
	}
}
