package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
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
func (d *Discography) sortBy(keys, order string) {
	sortKeys := strings.Split(keys, ",")
	if len(sortKeys) < 1 || len(sortKeys) > 2 {
		fmt.Println("Invalid number of keys. Provide one or two keys.")
		return
	}

	primaryKey := sortKeys[0]
	var secondaryKey string
	if len(sortKeys) == 2 {
		secondaryKey = sortKeys[1]
	}

	sort.SliceStable(d.Tracks, func(i, j int) bool {
		less := func(val1, val2 string) bool {
			if order == "descending" {
				return val1 > val2
			}
			return val1 < val2
		}

		// Primary sort key
		switch primaryKey {
		case "orchestra":
			if d.Tracks[i].Orchestra != d.Tracks[j].Orchestra {
				return less(d.Tracks[i].Orchestra, d.Tracks[j].Orchestra)
			}
		case "year":
			if d.Tracks[i].Year != d.Tracks[j].Year {
				return less(d.Tracks[i].Year, d.Tracks[j].Year)
			}
		case "name":
			if d.Tracks[i].Name != d.Tracks[j].Name {
				return less(d.Tracks[i].Name, d.Tracks[j].Name)
			}
		case "artist":
			if d.Tracks[i].Artist != d.Tracks[j].Artist {
				return less(d.Tracks[i].Artist, d.Tracks[j].Artist)
			}
		case "genre":
			if d.Tracks[i].Genre != d.Tracks[j].Genre {
				return less(d.Tracks[i].Genre, d.Tracks[j].Genre)
			}
		case "vocal":
			if d.Tracks[i].Vocal != d.Tracks[j].Vocal {
				return less(d.Tracks[i].Vocal, d.Tracks[j].Vocal)
			}
		default:
			fmt.Printf("Unsupported sort key: %s\n", primaryKey)
			return false
		}

		// Secondary sort key
		if secondaryKey != "" {
			switch secondaryKey {
			case "orchestra":
				return less(d.Tracks[i].Orchestra, d.Tracks[j].Orchestra)
			case "year":
				return less(d.Tracks[i].Year, d.Tracks[j].Year)
			case "name":
				return less(d.Tracks[i].Name, d.Tracks[j].Name)
			case "artist":
				return less(d.Tracks[i].Artist, d.Tracks[j].Artist)
			case "genre":
				return less(d.Tracks[i].Genre, d.Tracks[j].Genre)
			case "vocal":
				return less(d.Tracks[i].Vocal, d.Tracks[j].Vocal)
			default:
				fmt.Printf("Unsupported sort key: %s\n", secondaryKey)
				return false
			}
		}

		return false
	})
}

func (d *Discography) filterBy(filters map[string]string) {
	var filteredTracks []Track

	for _, track := range d.Tracks {
		match := true
		var matches []bool

		// Check all filters against the track
		var keys []string
		for key, value := range filters {
			keys = append(keys, key)
			switch key {
			case "orchestra":
				if strings.ToLower(track.Orchestra) != strings.ToLower(value) {
					match = false
				}
			case "year":
				// Try compiling as a regex
				re, err := regexp.Compile("^" + value + "$") // Ensure full match
				if err != nil {
					// If it's not a regex, do an exact match
					if strings.ToLower(track.Year) != strings.ToLower(value) {
						match = false
					}
				} else {
					// Use regex match
					if !re.MatchString(track.Year) {
						match = false
					}
				}
			case "name":
				if strings.ToLower(track.Name) != strings.ToLower(value) {
					match = false
				}
			case "artist":
				if strings.ToLower(track.Artist) != strings.ToLower(value) {
					match = false
				}
			case "genre":
				if strings.ToLower(track.Genre) != strings.ToLower(value) {
					match = false
				}
			case "vocal":
				if strings.ToLower(track.Vocal) != strings.ToLower(value) {
					match = false
				}
			default:
				fmt.Printf("Unsupported filter key: %s\n", key)
				match = false
			}

			if match {
				matches = append(matches, match)
			}

			// Stop checking further if any condition fails
			//             if !match {
			//                 break
			//             }
		}

		// Add the track to the filtered list if it matches all conditions
		if len(matches) == len(keys) {
			filteredTracks = append(filteredTracks, track)
		}
	}

	// Replace the original tracks with the filtered ones
	d.Tracks = filteredTracks
}

// Ensure unique tracks in the Discography struct
func (d *Discography) removeDuplicateTracks() {
	seen := make(map[string]bool)
	uniqueTracks := []Track{}

	for _, track := range d.Tracks {
		// Create a unique key for each track
		year := strings.Split(track.Year, "-")[0]
		key := fmt.Sprintf("%s|%s|%s|%s",
			track.Orchestra, year, track.Name, track.Genre)

		if !seen[key] {
			seen[key] = true
			uniqueTracks = append(uniqueTracks, track)
		}
	}

	d.Tracks = uniqueTracks
}

// helper function to read XML (discography) file
func readXMLFile(filename string) (*Discography, error) {
	// Match files using the provided filename pattern
	files, err := filepath.Glob(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to match files with pattern %s: %v", filename, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files matched the pattern %s", filename)
	}

	combinedDiscography := &Discography{}

	for _, fileName := range files {
		xmlFile, err := os.Open(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", fileName, err)
		}
		defer xmlFile.Close()

		byteValue, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", fileName, err)
		}

		var discography Discography
		err = xml.Unmarshal(byteValue, &discography)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal XML in file %s: %v", fileName, err)
		}

		// Patch orchestra in tracks if necessary
		for _, track := range discography.Tracks {
			if track.Orchestra == "" && discography.Orchestra != "" {
				track.Orchestra = discography.Orchestra
			}
			combinedDiscography.Tracks = append(combinedDiscography.Tracks, track)
		}

		// Set the orchestra attribute if it's not already set
		if combinedDiscography.Orchestra == "" && discography.Orchestra != "" {
			combinedDiscography.Orchestra = discography.Orchestra
		}
	}

	return combinedDiscography, nil
}

// helper function to read CSV (discography) file
func readCSVFile(filename string) (*Discography, error) {
	// Match files using the provided filename pattern
	files, err := filepath.Glob(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to match files with pattern %s: %v", filename, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files matched the pattern %s", filename)
	}

	combinedDiscography := &Discography{}

	for _, fileName := range files {
		csvFile, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
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

		// Patch orchestra in tracks if necessary
		for _, track := range discography.Tracks {
			if track.Orchestra == "" && discography.Orchestra != "" {
				track.Orchestra = discography.Orchestra
			}
			combinedDiscography.Tracks = append(combinedDiscography.Tracks, track)
		}

		// Set the orchestra attribute if it's not already set
		if combinedDiscography.Orchestra == "" && discography.Orchestra != "" {
			combinedDiscography.Orchestra = discography.Orchestra
		}

	}
	return combinedDiscography, nil
}

// helper function to read (discography) file
func readFile(filename, sortBy, sortOrder string, filters map[string]string) (*Discography, error) {
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
	if len(filters) > 0 {
		discography.filterBy(filters)
	}
	return discography, err
}
