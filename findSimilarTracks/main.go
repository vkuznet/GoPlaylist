package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func parseXMLFile(filename string) ([]Track, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var discography Discography
	err = xml.Unmarshal(byteValue, &discography)
	if err != nil {
		return nil, err
	}

	for i := range discography.Tracks {
		// Patch orchestra if missing
		if discography.Tracks[i].Orchestra == "" {
			discography.Tracks[i].Orchestra = discography.Orchestra
		}
	}

	return discography.Tracks, nil
}

func findSimilarTracks(tracks []Track) []Track {
	seen := make(map[string][]Track)
	var result []Track

	// Group tracks by name
	for _, track := range tracks {
		seen[track.Name] = append(seen[track.Name], track)
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
		key := fmt.Sprintf("%s|%s|%s", track.Name, track.Year, track.Orchestra)
		unique[key] = track
	}

	var finalTracks []Track
	for _, track := range unique {
		finalTracks = append(finalTracks, track)
	}

	return finalTracks
}

func writeXMLFile(outputFile string, tracks []Track) error {
	// Sort tracks by name
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Name < tracks[j].Name
	})

	tracksWrapper := Tracks{Tracks: tracks}
	data, err := xml.MarshalIndent(tracksWrapper, "", "  ")
	if err != nil {
		return err
	}

	data = append([]byte(xml.Header), data...)
	err = ioutil.WriteFile(outputFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var inputPattern, outputFile string
	var verbose, stats bool
	flag.StringVar(&inputPattern, "input", "*.xml", "Input XML files pattern")
	flag.StringVar(&outputFile, "output", "output.xml", "Output XML file")
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.BoolVar(&stats, "stats", false, "stats output")
	flag.Parse()

	// Use filepath.Walk to collect all matching files recursively
	arr := strings.Split(inputPattern, ".")
	ext := arr[len(arr)-1]
	var files []string
	err := filepath.Walk(filepath.Dir(inputPattern), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ext) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil || len(files) == 0 {
		fmt.Println("No input files found")
		os.Exit(1)
	}

	// Print each individual file name
	if verbose {
		for _, fname := range files {
			fmt.Println("Found", fname)
		}
	}

	var allTracks []Track
	for _, file := range files {
		tracks, err := parseXMLFile(file)
		if err != nil {
			fmt.Printf("Error parsing file %s: %v\n", file, err)
			os.Exit(1)
		}
		allTracks = append(allTracks, tracks...)
	}

	similarTracks := findSimilarTracks(allTracks)

	fmt.Printf("Collected %d tracks from %d files", len(allTracks), len(files))
	err = writeXMLFile(outputFile, similarTracks)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote %d similar tracks to %s\n", len(similarTracks), outputFile)
	if stats {
		PrintStats(similarTracks)
	}
}
