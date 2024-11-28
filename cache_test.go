package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCacheInit(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tmpDir := t.TempDir()
	cache := Cache{}

	service := "spotify"
	cache.Init(service, tmpDir)

	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Fatalf("Expected directory %s was not created", tmpDir)
	}
}

func TestCacheAddAndCheckTrack(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tmpDir := t.TempDir()
	cache := Cache{}

	service := "spotify"
	cache.Init(service, tmpDir)

	title := "MyPlaylist"
	playlistID := "12345"
	trackName := "MyTrack.mp3"

	err := cache.AddTrack(title, playlistID, trackName)
	if err != nil {
		t.Fatal(err)
	}

	if !cache.CheckTrack(title, playlistID, trackName) {
		t.Fatalf("Expected track %s to exist in cache but it does not", trackName)
	}
}

func TestCacheLoad(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cache := Cache{}
	service := "spotify"
	tmpDir := t.TempDir()

	cache.Init(service, tmpDir)

	title := "MyPlaylist"
	playlistID := "12345"
	trackNames := []string{"Track1.mp3", "Track2.mp3"}
	for _, trackName := range trackNames {
		err := cache.AddTrack(title, playlistID, trackName)
		if err != nil {
			fmt.Printf("fail to add track %s to cache %+v", trackName, cache)
		}
	}

	tracks, err := cache.Load(service, title, playlistID)
	if err != nil {
		t.Fatal(err)
	}

	if len(tracks) != len(trackNames) {
		fmt.Println("### generated tracks")
		for _, t := range trackNames {
			fmt.Println(t)
		}
		fmt.Println("### loaded tracks")
		for _, t := range tracks {
			fmt.Println(t)
		}
		t.Fatalf("Expected %d tracks in cache, found %d", len(trackNames), len(tracks))
	}

	for _, trackName := range trackNames {
		found := false
		for _, t := range tracks {
			if t == trackName {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Track %s not found in cache map", trackName)
		}
	}
}
