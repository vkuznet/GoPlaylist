package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Track represents track
type Track struct {
	Orchestra string `xml:"orchestra,attr"`
	Year      string `xml:"year,attr"`
	Name      string `xml:"name,attr"`
	Artist    string `xml:"artist,attr,omitempty"`
	Genre     string `xml:"genre,attr,omitempty"`
	Vocal     string `xml:"vocal,attr,omitempty"`
}

// String provides string representation of the track
func (t *Track) String() string {
	return fmt.Sprintf("%s,%s,%s,%s", t.Orchestra, t.Year, t.Name, t.Artist)
}

// helper function to construct track from its string representation
func constructTrack(t string) Track {
	arr := strings.Split(t, ",")
	// NOTE: array should match String() method
	track := Track{
		Orchestra: arr[0],
		Year:      arr[1],
		Name:      arr[2],
		Artist:    arr[3],
	}
	return track
}

// Discography represents discography object
type Discography struct {
	Orchestra string  `xml:"orchestra,attr"`
	Tracks    []Track `xml:"track"`
}

// helper function to read XML (discography) file
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

// helper function to read CSV (discography) file
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
		// NOTE: record should match String() method of Track object above
		track := Track{
			Orchestra: record[0],
			Year:      record[1],
			Name:      record[2],
		}
		if len(record) == 4 {
			track.Artist = record[3]
		} else if len(record) == 5 {
			track.Genre = record[4]
		} else if len(record) == 6 {
			track.Vocal = record[5]
		}
		discography.Tracks = append(discography.Tracks, track)
	}

	if Config.Verbose > 1 {
		fmt.Println("Parsed discography from CSV:")
		for _, track := range discography.Tracks {
			fmt.Printf("track: %+v\n", track)
		}
	}

	return &discography, nil
}

// helper function to read (discography) file
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
