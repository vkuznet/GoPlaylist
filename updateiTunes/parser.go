// parser.go
package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// Track represents a single track entry in the XML file.
type Track struct {
	Name     string `xml:"name,attr"`
	Vocal    string `xml:"vocal,attr"`
	Year     string `xml:"year,attr"`
	Genre    string `xml:"genre,attr"`
	Composer string `xml:"composer,attr"`
	Author   string `xml:"author,attr"`
	Label    string `xml:"label,attr"`
}

// Discography is the root element of the XML file containing multiple tracks.
type Discography struct {
	Tracks []Track `xml:"track"`
}

// ParseXML accepts a glob pattern (e.g., "Francisco Canaro*.xml"), reads all matching XML files,
// and returns a flattened slice of all tracks.
func ParseXML(pattern string) ([]Track, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("error parsing glob pattern: %v", err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no XML files matched the pattern: %s", pattern)
	}

	var allTracks []Track

	for _, file := range matches {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read file %s: %v\n", file, err)
			continue
		}

		var discography Discography
		if err := xml.Unmarshal(data, &discography); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not parse XML in file %s: %v\n", file, err)
			continue
		}

		allTracks = append(allTracks, discography.Tracks...)
	}

	return allTracks, nil
}
