package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheInit(t *testing.T) {
	tmpDir := t.TempDir()
	cache := Cache{Map: make(map[string]CacheEntry)}

	service := "spotify"
	cache.Init(service, tmpDir)

	expectedDir := filepath.Join(tmpDir, ".goplaylist", service)
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Fatalf("Expected directory %s was not created", expectedDir)
	}

	expectedFile := filepath.Join(expectedDir, "cache.csv")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Expected cache file %s was not created", expectedFile)
	}
}

func TestCacheAddAndCheckTrack(t *testing.T) {
	tmpDir := t.TempDir()
	cache := Cache{Map: make(map[string]CacheEntry)}

	service := "spotify"
	cache.Init(service, tmpDir)

	title := "MyPlaylist"
	playlistID := "12345"
	trackName := "MyTrack.mp3"

	cache.AddTrack(title, playlistID, trackName)

	if !cache.CheckTrack(title, playlistID, trackName) {
		t.Fatalf("Expected track %s to exist in cache but it does not", trackName)
	}
}

func TestCacheLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cache := Cache{Map: make(map[string]CacheEntry)}

	service := "spotify"
	cache.Init(service, tmpDir)

	title := "MyPlaylist"
	playlistID := "12345"
	trackNames := []string{"Track1.mp3", "Track2.mp3"}
	for _, trackName := range trackNames {
		cache.AddTrack(title, playlistID, trackName)
	}

	cache.Load(service)

	entry, exists := cache.Map[title]
	if !exists {
		t.Fatalf("Expected playlist %s to exist in cache map", title)
	}

	if len(entry.Tracks) != len(trackNames) {
		t.Fatalf("Expected %d tracks in cache, found %d", len(trackNames), len(entry.Tracks))
	}

	for _, trackName := range trackNames {
		found := false
		for _, t := range entry.Tracks {
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
