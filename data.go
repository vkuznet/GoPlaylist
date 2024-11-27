package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Track struct {
	Name      string `xml:"name,attr"`
	Genre     string `xml:"genre,attr,omitempty"`
	Year      string `xml:"year,attr"`
	Orchestra string `xml:"orchestra,attr"`
}

type Discography struct {
	Orchestra string  `xml:"orchestra,attr"`
	Tracks    []Track `xml:"track"`
}

func readXMLFile(filename string) (*Discography, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var discography Discography
	xml.Unmarshal(byteValue, &discography)

	var tracks []Track
	for _, track := range discography.Tracks {
		if track.Orchestra == "" && discography.Orchestra != "" {
			track.Orchestra = discography.Orchestra
			tracks = append(tracks, track)
		}
	}
	if len(tracks) > 0 {
		discography.Tracks = tracks
	}

	return &discography, nil
}

func readCSVFile(filename string) (*Discography, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per line

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var discography Discography
	for _, record := range records {
		if len(record) < 3 {
			continue // Skip rows with insufficient data
		}
		discography.Tracks = append(discography.Tracks, Track{
			Orchestra: record[0], // Orchestra
			Year:      record[1], // Year
			Name:      record[2], // Track name
		})
	}

	if Config.Verbose > 1 {
		fmt.Println("Parsed discography from CSV:")
		for _, track := range discography.Tracks {
			fmt.Printf("track: %+v\n", track)
		}
	}

	return &discography, nil
}

func readFile(filename string) (*Discography, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".xml":
		return readXMLFile(filename)
	case ".csv":
		return readCSVFile(filename)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}
