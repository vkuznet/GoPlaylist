package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var xmlPattern, musicDir, orchestra, matchMode string
	var dryRun, fixTitle bool
	var verbose int
	flag.StringVar(&xmlPattern, "xml", "*.xml", "Glob pattern for XML discography files")
	mdir := fmt.Sprintf("%s/Music/iTunes/iTunes Music", os.Getenv("HOME"))
	flag.StringVar(&musicDir, "musicDir", mdir, "Directory containing iTunes music")
	flag.StringVar(&orchestra, "orchestra", "", "Orchestra name to match subdirectories")
	flag.BoolVar(&dryRun, "dryRun", false, "If set, tags won't be written to files")
	flag.IntVar(&verbose, "verbose", 0, "verbose mode")
	flag.StringVar(&matchMode, "matchMode", "strict", "Match mode: strict or fuzzy")
	flag.BoolVar(&fixTitle, "fixTitle", false, "If set, titles of songs will be fixde accodring to discography")

	flag.Parse()

	// Load discography
	tracks, err := ParseXML(xmlPattern)
	if err != nil {
		log.Fatalf("Failed to parse XML: %v", err)
	}
	fmt.Printf("Found %d tracks in %s\n", len(tracks), xmlPattern)

	if orchestra == "" {
		log.Fatalln("no orchestra is provided")
	}

	// Find matching music files
	files := []string{}
	err = filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp3") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to scan music directory: %v", err)
	}
	fmt.Printf("Found %d files in %s\n", len(files), musicDir)

	// Perform matching
	matches := MatchTracks(files, tracks, MatchMode(matchMode), verbose)
	fmt.Printf("Found %d matches\n", len(matches))

	// Process results
	for _, match := range matches {
		fmt.Printf("Matched: %s <-> %s\n", match.Track.Name, match.FilePath)
		err := UpdateTags(orchestra, match.FilePath, match.Track, fixTitle, dryRun, verbose)
		if err != nil {
			fmt.Printf("Error updating tags for %s: %v\n", match.FilePath, err)
		}
	}
}
