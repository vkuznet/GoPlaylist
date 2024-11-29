package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Cache represents local cache object
type Cache struct {
	Dir    string
	Tracks []Track
}

// Init method initialize cache
func (c *Cache) Init(service, dir string) {
	c.Dir = dir

	// Check if the directory exists
	info, err := os.Stat(c.Dir)
	if os.IsNotExist(err) {
		// Directory does not exist, create it
		err = os.MkdirAll(c.Dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Cache directory created:", c.Dir)
	} else if err != nil {
		log.Fatal(err)
	} else if !info.IsDir() {
		log.Fatalf("Path %s exists but is not a directory", c.Dir)
	}
}

// helper function to construct cache file
func (c *Cache) cacheFile(title, playlistID string) (string, error) {
	cdir := fmt.Sprintf("%s/%s/%s", c.Dir, title, playlistID)
	err := os.MkdirAll(cdir, os.ModePerm)
	if err != nil {
		return "", err
	}
	cacheFile := fmt.Sprintf("%s/cache.txt", cdir)
	if _, e := os.Stat(cacheFile); os.IsNotExist(e) {
		file, err := os.Create(cacheFile)
		if err != nil {
			log.Fatalf("Error creating cache file: %v", err)
		}
		file.Close()
		fmt.Println("Cache file created:", cacheFile)
	}
	return cacheFile, nil
}

// AddTrack adds new track object to cache
func (c *Cache) AddTrack(title, playlistID string, track Track) error {
	var file *os.File
	var err error

	// Check if the track already exists
	if c.CheckTrack(title, playlistID, track) {
		return nil // No need to add, already exists
	}

	// Open the cache file in append mode
	cacheFile, err := c.cacheFile(title, playlistID)
	if err != nil {
		return err
	}

	// If the file doesn't exist, create it
	if _, e := os.Stat(cacheFile); os.IsNotExist(e) {
		file, err = os.Create(cacheFile)
		if err != nil {
			log.Fatalf("Error creating cache file: %v", err)
		}
		file.Close()
		fmt.Println("Cache file created:", cacheFile)
	} else {
		// open file in append mode
		file, err = os.OpenFile(cacheFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error opening cache file: %w", err)
		}
	}
	defer file.Close()

	// Append the track to our cache file
	if _, err := file.WriteString(track.String() + "\n"); err != nil {
		return fmt.Errorf("error writing to cache file: %w", err)
	}

	return nil
}

// CheckTrack checks track within existing local cache
func (c *Cache) CheckTrack(title, playlistID string, track Track) bool {
	cacheFile, err := c.cacheFile(title, playlistID)
	if err != nil {
		return false
	}

	// Open the cache file
	file, err := os.Open(cacheFile)
	if err != nil {
		log.Fatalf("Error opening cache file: %v", err)
	}
	defer file.Close()

	// Perform a line-by-line scan
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == track.String() {
			return true
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error scanning cache file: %v", err)
	}

	return false
}

// Load method loads track from local cache and return list of track objects
func (c *Cache) Load(service, title, playlistID string) ([]Track, error) {
	var tracks []Track

	// Open the cache file
	cacheFile, err := c.cacheFile(title, playlistID)
	if err != nil {
		return tracks, err
	}
	file, err := os.Open(cacheFile)
	if err != nil {
		return tracks, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return tracks, err
	}
	for _, t := range strings.Split(string(data), "\n") {
		if t != "" {
			track := constructTrack(t)
			tracks = append(tracks, track)
		}
	}
	return tracks, nil
}
