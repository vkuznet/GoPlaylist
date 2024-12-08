package main

import (
	"encoding/xml"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestXMLParsing(t *testing.T) {
	file := "testplaylist.xml"
	discography, err := readXMLFile(file)
	if err != nil {
		t.Error(err)
	}
	orchestra := getOrchestra("bla", discography)
	if orchestra != "bla" {
		t.Errorf("wrong orchestra %s, discography %+v", orchestra, discography)
	}

	// loop over tracks to see orchestra
	for _, track := range discography.Tracks {
		if strings.Contains(track.Name, "Sin") {
			if track.Orchestra != "Orquesta Tipica Victor" {
				t.Errorf("wrong orchestra %+v", track)
			}
		}
	}

}

func TestWriteAndReadXMLFile(t *testing.T) {
	xmlFilename := "test.xml"

	expectedDiscography := &Discography{
		Orchestra: "Ricardo Tanturi",
		Tracks: []Track{
			{Name: "Una noche más", Orchestra: "Ricardo Tanturi", Year: "1941"},
			{Name: "En el salón", Orchestra: "Ricardo Tanturi", Year: "1943"},
		},
	}

	// Write XML file
	file, err := os.Create(xmlFilename)
	if err != nil {
		t.Fatalf("Failed to create XML file: %v", err)
	}
	defer os.Remove(xmlFilename) // Cleanup after the test
	defer file.Close()

	// Write the XML content to the file
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(expectedDiscography); err != nil {
		t.Fatalf("Failed to write XML to file: %v", err)
	}

	// Read XML file back
	discography, err := readXMLFile(xmlFilename)
	if err != nil {
		t.Fatalf("Failed to read XML file: %v", err)
	}

	// Compare the expected and actual results
	if !reflect.DeepEqual(discography, expectedDiscography) {
		t.Errorf("Expected %+v, got %+v", expectedDiscography, discography)
	}
}

// Test writing and reading CSV
func TestWriteAndReadCSV(t *testing.T) {
	csvFilename := "test.csv"
	expectedDiscography := &Discography{
		Tracks: []Track{
			{Name: "Una noche más", Orchestra: "Ricardo Tanturi", Year: "1941"},
			{Name: "En el salón", Orchestra: "Ricardo Tanturi", Year: "1943"},
		},
	}

	// Write CSV
	csvFile, err := os.Create(csvFilename)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}
	defer os.Remove(csvFilename)
	defer csvFile.Close()

	for _, track := range expectedDiscography.Tracks {
		_, err := csvFile.WriteString(track.Orchestra + "," + track.Year + "," + track.Name + "\n")
		if err != nil {
			t.Fatalf("Failed to write to CSV file: %v", err)
		}
	}

	csvFile.Close()

	// Read CSV
	discography, err := readCSVFile(csvFilename)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if !reflect.DeepEqual(discography, expectedDiscography) {
		t.Errorf("Expected %+v, got %+v", expectedDiscography, discography)
	}
}

// Test generalized readFile function
func TestReadFile(t *testing.T) {
	xmlFilename := "test.xml"
	csvFilename := "test.csv"

	// Create CSV for test
	expectedDiscography := &Discography{
		Tracks: []Track{
			{Name: "Una noche más", Orchestra: "Ricardo Tanturi", Year: "1941"},
			{Name: "En el salón", Orchestra: "Ricardo Tanturi", Year: "1943"},
		},
	}

	csvFile, err := os.Create(csvFilename)
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}
	defer os.Remove(csvFilename)
	defer csvFile.Close()

	for _, track := range expectedDiscography.Tracks {
		_, err := csvFile.WriteString(track.Orchestra + "," + track.Year + "," + track.Name + "\n")
		if err != nil {
			t.Fatalf("Failed to write to CSV file: %v", err)
		}
	}

	csvFile.Close()

	// Test readFile with CSV
	discography, err := readFile(csvFilename, "", "", nil)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if !reflect.DeepEqual(discography, expectedDiscography) {
		t.Errorf("Fail to process csv file: expected %+v, got %+v", expectedDiscography, discography)
	}

	// Write XML file
	file, err := os.Create(xmlFilename)
	if err != nil {
		t.Fatalf("Failed to create XML file: %v", err)
	}
	defer os.Remove(xmlFilename) // Cleanup after the test
	defer file.Close()

	// Write the XML content to the file
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(expectedDiscography); err != nil {
		t.Fatalf("Failed to write XML to file: %v", err)
	}

	// Test readFile with XML
	discography, err = readFile(xmlFilename, "", "", nil)
	if err != nil {
		t.Fatalf("Failed to read XML file: %v", err)
	}

	//     if discography.Orchestra != "Ricardo Tanturi" {
	//         t.Errorf("Expected orchestra 'Ricardo Tanturi', got '%s'", discography.Orchestra)
	//     }

	// Test readFile with XML
	discography, err = readFile(xmlFilename, "", "", nil)
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if !reflect.DeepEqual(discography, expectedDiscography) {
		t.Errorf("fail to process xml file: expected %+v, got %+v", expectedDiscography, discography)
	}
}

func TestSortBy(t *testing.T) {
	// Create a sample Discography object
	discography := Discography{
		Orchestra: "Orchestra1",
		Tracks: []Track{
			{Orchestra: "Orchestra2", Year: "2020", Name: "Track C", Artist: "Artist A", Genre: "Genre X", Vocal: "Vocal X"},
			{Orchestra: "Orchestra3", Year: "2019", Name: "Track A", Artist: "Artist B", Genre: "Genre Z", Vocal: "Vocal Y"},
			{Orchestra: "Orchestra1", Year: "2021", Name: "Track B", Artist: "Artist A", Genre: "Genre Y", Vocal: "Vocal Z"},
		},
	}

	tests := []struct {
		attr        string
		order       string
		expected    []Track
		description string
	}{
		{
			attr:  "year",
			order: "ascending",
			expected: []Track{
				{Orchestra: "Orchestra3", Year: "2019", Name: "Track A", Artist: "Artist B", Genre: "Genre Z", Vocal: "Vocal Y"},
				{Orchestra: "Orchestra2", Year: "2020", Name: "Track C", Artist: "Artist A", Genre: "Genre X", Vocal: "Vocal X"},
				{Orchestra: "Orchestra1", Year: "2021", Name: "Track B", Artist: "Artist A", Genre: "Genre Y", Vocal: "Vocal Z"},
			},
			description: "Sort by year in ascending order",
		},
		{
			attr:  "name",
			order: "descending",
			expected: []Track{
				{Orchestra: "Orchestra2", Year: "2020", Name: "Track C", Artist: "Artist A", Genre: "Genre X", Vocal: "Vocal X"},
				{Orchestra: "Orchestra1", Year: "2021", Name: "Track B", Artist: "Artist A", Genre: "Genre Y", Vocal: "Vocal Z"},
				{Orchestra: "Orchestra3", Year: "2019", Name: "Track A", Artist: "Artist B", Genre: "Genre Z", Vocal: "Vocal Y"},
			},
			description: "Sort by name in descending order",
		},
		{
			attr:  "artist",
			order: "ascending",
			expected: []Track{
				{Orchestra: "Orchestra2", Year: "2020", Name: "Track C", Artist: "Artist A", Genre: "Genre X", Vocal: "Vocal X"},
				{Orchestra: "Orchestra1", Year: "2021", Name: "Track B", Artist: "Artist A", Genre: "Genre Y", Vocal: "Vocal Z"},
				{Orchestra: "Orchestra3", Year: "2019", Name: "Track A", Artist: "Artist B", Genre: "Genre Z", Vocal: "Vocal Y"},
			},
			description: "Sort by artist in ascending order",
		},
		{
			attr:  "genre",
			order: "ascending",
			expected: []Track{
				{Orchestra: "Orchestra2", Year: "2020", Name: "Track C", Artist: "Artist A", Genre: "Genre X", Vocal: "Vocal X"},
				{Orchestra: "Orchestra1", Year: "2021", Name: "Track B", Artist: "Artist A", Genre: "Genre Y", Vocal: "Vocal Z"},
				{Orchestra: "Orchestra3", Year: "2019", Name: "Track A", Artist: "Artist B", Genre: "Genre Z", Vocal: "Vocal Y"},
			},
			description: "Sort by genre in ascending order",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Make a copy of the original discography to sort
			d := discography
			d.sortBy(test.attr, test.order)

			// Verify the result matches the expected order
			if !reflect.DeepEqual(d.Tracks, test.expected) {
				t.Errorf("Failed %s\nGot: %+v\nExpected: %+v", test.description, d.Tracks, test.expected)
			}
		})
	}
}

func TestSortByMultipleKeys(t *testing.T) {
	d := &Discography{
		Tracks: []Track{
			{Name: "Song A", Year: "2000", Orchestra: "Orch2"},
			{Name: "Song B", Year: "1999", Orchestra: "Orch1"},
			{Name: "Song C", Year: "2000", Orchestra: "Orch1"},
			{Name: "Song D", Year: "1998", Orchestra: "Orch2"},
		},
	}

	// Sort by year in ascending order
	d.sortBy("year", "ascending")

	expectedOrder := []string{"Song D", "Song B", "Song A", "Song C"}
	for i, track := range d.Tracks {
		if track.Name != expectedOrder[i] {
			t.Errorf("Expected %s, but got %s", expectedOrder[i], track.Name)
		}
	}

	// Sort by orchestra and then by year
	d.sortBy("orchestra,year", "ascending")

	expectedOrder = []string{"Song B", "Song C", "Song D", "Song A"}
	for i, track := range d.Tracks {
		if track.Name != expectedOrder[i] {
			t.Errorf("Expected %s, but got %s", expectedOrder[i], track.Name)
		}
	}
}

func TestFilterBy(t *testing.T) {
	d := &Discography{
		Tracks: []Track{
			{Name: "Song A", Year: "2000", Orchestra: "Orch1", Genre: "Classical"},
			{Name: "Song B", Year: "1999", Orchestra: "Orch1", Genre: "Jazz"},
			{Name: "Song C", Year: "2000", Orchestra: "Orch2", Genre: "Classical"},
			{Name: "Song D", Year: "1998", Orchestra: "Orch2", Genre: "Jazz"},
		},
	}

	// Filter by "orchestra" = "Orch1"
	d.filterBy(map[string]string{
		"orchestra": "Orch1",
	})

	expectedTracks := []Track{
		{Name: "Song A", Year: "2000", Orchestra: "Orch1", Genre: "Classical"},
		{Name: "Song B", Year: "1999", Orchestra: "Orch1", Genre: "Jazz"},
	}

	if len(d.Tracks) != len(expectedTracks) {
		t.Fatalf("Expected %d tracks, but got %d", len(expectedTracks), len(d.Tracks))
	}

	for i, track := range d.Tracks {
		if track.Name != expectedTracks[i].Name || track.Orchestra != expectedTracks[i].Orchestra {
			t.Errorf("Track mismatch. Expected %+v, but got %+v", expectedTracks[i], track)
		}
	}

	// Filter by "orchestra" = "Orch1" and "genre" = "Classical"
	d = &Discography{
		Tracks: []Track{
			{Name: "Song A", Year: "2000", Orchestra: "Orch1", Genre: "Classical"},
			{Name: "Song B", Year: "1999", Orchestra: "Orch1", Genre: "Jazz"},
			{Name: "Song C", Year: "2000", Orchestra: "Orch2", Genre: "Classical"},
			{Name: "Song D", Year: "1998", Orchestra: "Orch2", Genre: "Jazz"},
		},
	}
	d.filterBy(map[string]string{
		"orchestra": "Orch1",
		"genre":     "Classical",
	})

	expectedTracks = []Track{
		{Name: "Song A", Year: "2000", Orchestra: "Orch1", Genre: "Classical"},
	}

	if len(d.Tracks) != len(expectedTracks) {
		t.Fatalf("Expected %d tracks, but got %d", len(expectedTracks), len(d.Tracks))
	}

	for i, track := range d.Tracks {
		if track.Name != expectedTracks[i].Name || track.Orchestra != expectedTracks[i].Orchestra || track.Genre != expectedTracks[i].Genre {
			t.Errorf("Track mismatch. Expected %+v, but got %+v", expectedTracks[i], track)
		}
	}
}
