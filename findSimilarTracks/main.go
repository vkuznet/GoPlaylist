package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

	fmt.Printf("Collected %d tracks from %d files\n", len(allTracks), len(files))
	err = writeXMLFile(outputFile, similarTracks)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote %d multi-performed tracks to %s\n", len(similarTracks), outputFile)
	if stats {
		PrintStats(similarTracks)
	}
}
