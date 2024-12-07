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
	discography, err := readFile(csvFilename, "", "")
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
	discography, err = readFile(xmlFilename, "", "")
	if err != nil {
		t.Fatalf("Failed to read XML file: %v", err)
	}

	//     if discography.Orchestra != "Ricardo Tanturi" {
	//         t.Errorf("Expected orchestra 'Ricardo Tanturi', got '%s'", discography.Orchestra)
	//     }

	// Test readFile with XML
	discography, err = readFile(xmlFilename, "", "")
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
