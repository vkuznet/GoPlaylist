package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Dict struct {
	Key     string `xml:"key"`
	Dict    *Dict  `xml:"dict"`
	Integer int    `xml:"integer"`
	String  string `xml:"string"`
}

func main() {
	// Define command-line flags
	xmlInput := flag.String("xmlInput", "", "Path to the XML input file")
	csvOutput := flag.String("csvOutput", "", "Path to the CSV output file")
	flag.Parse()

	// Check if the required flags are provided
	if *xmlInput == "" || *csvOutput == "" {
		fmt.Println("Usage: go run main.go -xmlInput <file.xml> -csvOutput <file.csv>")
		os.Exit(1)
	}

	// Open the XML file
	xmlFile, err := os.Open(*xmlInput)
	if err != nil {
		fmt.Printf("Error opening XML file: %v\n", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	// Create CSV file
	csvFile, err := os.Create(*csvOutput)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Parse the XML data
	decoder := xml.NewDecoder(xmlFile)

	var currentKey string
	var trackID, year int
	var trackName, artist, genre string

	playlist := make(map[string]any)

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		//         fmt.Printf("### token %+v, type %T\n", token, token)

		switch elem := token.(type) {
		case xml.StartElement:
			if elem.Name.Local == "key" {
				var key string
				decoder.DecodeElement(&key, &elem)
				currentKey = key
				if key == "Playlists" {
					break // we are done with writing tracks to CSV file
				}
			} else if elem.Name.Local == "integer" && currentKey == "Track ID" {
				decoder.DecodeElement(&trackID, &elem)
			} else if elem.Name.Local == "string" && currentKey == "Name" {
				decoder.DecodeElement(&trackName, &elem)
			} else if elem.Name.Local == "string" && currentKey == "Artist" {
				decoder.DecodeElement(&artist, &elem)
			} else if elem.Name.Local == "string" && currentKey == "Genre" {
				decoder.DecodeElement(&genre, &elem)
			} else if elem.Name.Local == "integer" && currentKey == "Year" {
				decoder.DecodeElement(&year, &elem)
			}
		}
		if trackID > 0 && trackName != "" && artist != "" && genre != "" && year > 0 {
			record := []string{
				artist,
				fmt.Sprintf("%d", year),
				trackName,
				genre,
			}
			line := strings.Join(record, ",")
			if _, ok := playlist[line]; !ok {
				writer.Write(record)
				playlist[line] = true
				trackID = 0
				year = 0
				trackName = ""
				artist = ""
				genre = ""
				continue
			}
		}
	}
}
