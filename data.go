package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Discography struct {
	XMLName   xml.Name `xml:"discography"`
	Orchestra string   `xml:"orchestra,attr"`
	Tracks    []Track  `xml:"track"`
}

type Track struct {
	Name       string `xml:"name,attr"`
	Vocal      string `xml:"vocal,attr"`
	Year       string `xml:"year,attr"`
	Genre      string `xml:"genre,attr"`
	Duration   string `xml:"duration,attr"`
	Popularity int    `xml:"popularity,attr"`
	Orchestra  string `xml:"orchestra,attr"`
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

	if Config.Verbose > 1 {
		fmt.Println("orchestra:", discography.Orchestra)
		for _, track := range discography.Tracks {
			fmt.Printf("track: %+v\n", track)
		}
	}

	return &discography, nil
}
