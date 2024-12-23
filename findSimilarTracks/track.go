package main

import (
	"encoding/xml"
	"fmt"
)

type Track struct {
	Name      string `xml:"name,attr"`
	Year      string `xml:"year,attr"`
	Orchestra string `xml:"orchestra,attr"`
	Genre     string `xml:"genre,attr"`
}

type Discography struct {
	Orchestra string  `xml:"orchestra,attr"`
	Tracks    []Track `xml:"track"`
}

type Tracks struct {
	XMLName xml.Name `xml:"tracks"`
	Tracks  []Track  `xml:"track"`
}

func findSimilarTracks(tracks []Track) []Track {
	seen := make(map[string][]Track)
	var result []Track

	// Group tracks by name
	for _, track := range tracks {
		tName := capitalize(track.Name)
		seen[tName] = append(seen[tName], track)
	}

	// Find tracks with the same name but different year or orchestra
	for _, group := range seen {
		if len(group) > 1 {
			for i := 0; i < len(group)-1; i++ {
				for j := i + 1; j < len(group); j++ {
					if group[i].Year != group[j].Year || group[i].Orchestra != group[j].Orchestra {
						result = append(result, group[i], group[j])
					}
				}
			}
		}
	}

	// Remove duplicates from result
	unique := make(map[string]Track)
	for _, track := range result {
		tName := capitalize(track.Name)
		key := fmt.Sprintf("%s|%s|%s", tName, track.Year, track.Orchestra)
		unique[key] = track
	}

	var finalTracks []Track
	for _, track := range unique {
		finalTracks = append(finalTracks, track)
	}

	return finalTracks
}
