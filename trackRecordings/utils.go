package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// helper function to write XML file
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

// helper function to parse intput XML file
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
			discography.Tracks[i].Orchestra = orchestra(discography.Orchestra)
		}
	}

	return discography.Tracks, nil
}

// helper function to capitalize word
func capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	word = strings.ToLower(word)
	return string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
}

// helper function convertUTFToASCII removes diacritics
func convertUTFToASCII(input string) string {
	// Normalize the string to decompose characters (NFD form)
	t := norm.NFD.String(input)

	// Filter out non-spacing marks (diacritics)
	var result strings.Builder
	for _, r := range t {
		if unicode.IsMark(r) {
			continue // Skip diacritical marks
		}
		result.WriteRune(r)
	}
	return result.String()
}

// ConvertUTFToASCII removes diacritics and converts the string to lowercase ASCII
func ConvertUTFToASCII(input string) string {
	// Convert to lowercase
	return strings.ToLower(convertUTFToASCII(input))
}

// helper function to normalize track orchestra
func orchestra(o string) string {
	return convertUTFToASCII(strings.Trim(strings.Split(o, "(")[0], " "))
}
