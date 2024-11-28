package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type CacheEntry struct {
	PlaylistID string
	Tracks     []string
}

type Cache struct {
	Map       map[string]CacheEntry
	mu        sync.Mutex
	cacheFile string
}

func (c *Cache) Init(service, dir string) {
	if c.Map == nil {
		c.Map = make(map[string]CacheEntry)
	}
	serviceDir := filepath.Join(dir, ".goplaylist", service)

	// Check if the directory exists
	info, err := os.Stat(serviceDir)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(serviceDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Cache directory created:", serviceDir)
	} else if err != nil {
		log.Fatal(err)
	} else if !info.IsDir() {
		log.Fatalf("Path %s exists but is not a directory", serviceDir)
	}

	// Initialize the cache file path
	c.cacheFile = filepath.Join(serviceDir, "cache.csv")

	// If the file doesn't exist, create it
	if _, err := os.Stat(c.cacheFile); os.IsNotExist(err) {
		file, err := os.Create(c.cacheFile)
		if err != nil {
			log.Fatalf("Error creating cache file: %v", err)
		}
		file.Close()
		fmt.Println("Cache file created:", c.cacheFile)
	}
}

func (c *Cache) AddTrack(title, playlistID, trackName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Open the cache file in append mode
	file, err := os.OpenFile(c.cacheFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("Error opening cache file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the track information to the file
	err = writer.Write([]string{title, playlistID, trackName})
	if err != nil {
		log.Fatalf("Error writing to cache file: %v", err)
	}
}

func (c *Cache) CheckTrack(title, playlistID, trackName string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Open the cache file
	file, err := os.Open(c.cacheFile)
	if err != nil {
		log.Fatalf("Error opening cache file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading cache file: %v", err)
	}

	// Search for the track in the records
	for _, record := range records {
		if len(record) >= 3 && record[0] == title && record[1] == playlistID && record[2] == trackName {
			return true
		}
	}

	return false
}

func (c *Cache) Load(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Open the cache file
	file, err := os.Open(c.cacheFile)
	if err != nil {
		log.Fatalf("Error opening cache file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading cache file: %v", err)
	}

	// Populate the cache map
	c.Map = make(map[string]CacheEntry)
	for _, record := range records {
		if len(record) >= 3 {
			title := record[0]
			playlistID := record[1]
			trackName := record[2]

			entry, exists := c.Map[title]
			if !exists {
				entry = CacheEntry{PlaylistID: playlistID, Tracks: []string{}}
			}
			entry.Tracks = append(entry.Tracks, trackName)
			c.Map[title] = entry
		}
	}
}
