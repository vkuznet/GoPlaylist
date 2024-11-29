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
	track := Track{Name: "Name", Year: "Year", Orchestra: "Orchestra"}

	err := cache.AddTrack(title, playlistID, track)
	if err != nil {
		t.Fatal(err)
	}

	if !cache.CheckTrack(title, playlistID, track) {
		t.Fatalf("Expected track %s to exist in cache but it does not", track.String())
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
	tracks := []Track{
		Track{Name: "Name", Year: "Year", Orchestra: "Orchestra1"},
		Track{Name: "Name", Year: "Year", Orchestra: "Orchestra2"},
	}
	for _, track := range tracks {
		err := cache.AddTrack(title, playlistID, track)
		if err != nil {
			fmt.Printf("fail to add track %s to cache %+v", track.String(), cache)
		}
	}

	cacheTracks, err := cache.Load(service, title, playlistID)
	if err != nil {
		t.Fatal(err)
	}

	if len(tracks) != len(cacheTracks) {
		fmt.Println("### generated tracks")
		for _, t := range tracks {
			fmt.Println(t)
		}
		fmt.Println("### cached tracks")
		for _, t := range cacheTracks {
			fmt.Println(t)
		}
		t.Fatalf("Expected %d tracks in cache, found %d", len(tracks), len(cacheTracks))
	}

	for _, track := range tracks {
		found := false
		for _, t := range cacheTracks {
			if t.String() == track.String() {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Track %s not found in cache map", track.String())
		}
	}
}
