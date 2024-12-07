package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	record := strings.Split(t, ",")
	// NOTE: record should match String() method
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
	return track
}

// Discography represents discography object
type Discography struct {
	Orchestra string  `xml:"orchestra,attr"`
	Tracks    []Track `xml:"track"`
}

// SortBy sorts the tracks of the Discography by the specified attribute and order.
// Supported attributes: "orchestra", "year", "name", "artist", "genre", "vocal"
// Order can be "ascending" or "descending".
func (d *Discography) sortBy(attr, order string) {
	sort.Slice(d.Tracks, func(i, j int) bool {
		var less bool
		switch strings.ToLower(attr) {
		case "orchestra":
			less = d.Tracks[i].Orchestra < d.Tracks[j].Orchestra
		case "year":
			less = d.Tracks[i].Year < d.Tracks[j].Year
		case "name":
			less = d.Tracks[i].Name < d.Tracks[j].Name
		case "artist":
			less = d.Tracks[i].Artist < d.Tracks[j].Artist
		case "genre":
			less = d.Tracks[i].Genre < d.Tracks[j].Genre
		case "vocal":
			less = d.Tracks[i].Vocal < d.Tracks[j].Vocal
		default:
			log.Printf("Invalid attribute '%s' for sorting. No action taken.", attr)
			return false
		}

		// Reverse the sorting order if "descending"
		if strings.ToLower(order) == "descending" || strings.ToLower(order) == "desc" {
			return !less
		}
		return less
	})
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
func readFile(filename, sortBy, sortOrder string) (*Discography, error) {
	ext := filepath.Ext(filename)
	var discography *Discography
	var err error
	switch ext {
	case ".xml":
		discography, err = readXMLFile(filename)
	case ".csv":
		discography, err = readCSVFile(filename)
	default:
		err = fmt.Errorf("unsupported file format: %s", ext)
	}
	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "ascending"
		}
		discography.sortBy(sortBy, sortOrder)
	}
	return discography, err
}
