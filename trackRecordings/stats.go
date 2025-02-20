package main

import (
	"fmt"
	"sort"
	"strings"
)

func PrintStats(tracks []Track) {
	stats := make(map[string]map[string]map[string]struct{}) // name -> orchestra -> year -> struct{}

	tmap := make(map[string]string) // name -> genre
	for _, track := range tracks {
		tName := capitalize(ConvertUTFToASCII(track.Name))
		if track.Year == "" {
			track.Year = "19xx"
		}
		if _, ok := stats[tName]; !ok {
			stats[tName] = make(map[string]map[string]struct{})
		}
		if _, ok := stats[tName][orchestra(track.Orchestra)]; !ok {
			stats[tName][orchestra(track.Orchestra)] = make(map[string]struct{})
		}
		stats[tName][orchestra(track.Orchestra)][track.Year] = struct{}{}
		tmap[tName] = capitalize(track.Genre)
	}

	fmt.Println("### Statistics")
	for tName, orchestras := range stats {
		if genre, ok := tmap[tName]; ok {
			fmt.Printf("Track: %s (%s)\n", tName, genre)
		} else {
			fmt.Printf("Track: %s\n", tName)
		}
		total := 0
		for orchestra, years := range orchestras {
			yearsList := make([]string, 0, len(years))
			for year := range years {
				yearsList = append(yearsList, year)
			}
			sort.Strings(yearsList)
			fmt.Printf("  Recorded %d times by %s in years: %s\n", len(years), orchestra, strings.Join(yearsList, ", "))
			total += len(years)
		}
		fmt.Printf("- total %d times\n", total)
	}
}
