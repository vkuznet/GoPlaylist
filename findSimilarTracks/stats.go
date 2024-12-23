package main

import (
	"fmt"
	"sort"
	"strings"
)

func PrintStats(tracks []Track) {
	stats := make(map[string]map[string]map[string]struct{}) // name -> orchestra -> year -> struct{}

	for _, track := range tracks {
		tName := capitalize(track.Name)
		if track.Year == "" {
			track.Year = "19xx"
		}
		if _, ok := stats[tName]; !ok {
			stats[tName] = make(map[string]map[string]struct{})
		}
		if _, ok := stats[tName][track.Orchestra]; !ok {
			stats[tName][track.Orchestra] = make(map[string]struct{})
		}
		stats[tName][track.Orchestra][track.Year] = struct{}{}
	}

	fmt.Println("### Statistics")
	for name, orchestras := range stats {
		fmt.Printf("Track: %s\n", name)
		total := 0
		for orchestra, years := range orchestras {
			yearsList := make([]string, 0, len(years))
			for year := range years {
				yearsList = append(yearsList, year)
			}
			sort.Strings(yearsList)
			fmt.Printf("  Performed %d times by %s in years: %s\n", len(years), orchestra, strings.Join(yearsList, ", "))
			total += len(years)
		}
		fmt.Printf("- total %d times\n", total)
	}
}
