package main

// config module
//
// Copyright (c) 2020 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Configuration stores server configuration parameters
type Configuration struct {
	YoutubeId     string `json:"youtube_id"`
	YoutubeSecret string `json:"youtube_secret"`
	SpotifyId     string `json:"spotify_id"`
	SpotifySecret string `json:"spotify_secret"`
	CallbackPort  int    `json:"callback_port"`
	Service       string `json:"service"`
	PlaylistTitle string `json:"playlist_title"`
	Verbose       int    `json:"verbose"`
}

// Config variable represents configuration object
var Config Configuration

// helper function to parse server configuration file
func parseConfig(configFile string) error {
	data, err := os.ReadFile(filepath.Clean(configFile))
	if err != nil {
		log.Println("Unable to read", err)
		return err
	}
	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Println("Unable to parse", err)
		return err
	}
	return nil
}
