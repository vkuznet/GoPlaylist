package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// toSentenceCase converts a string to sentence case: first letter uppercase, rest lowercase.
func toSentenceCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(strings.ToLower(s))  // Convert to lowercase first
	runes[0] = unicode.ToUpper(runes[0]) // Capitalize first letter
	return string(runes)
}

func main() {
	// Define command-line flags
	csvFile := flag.String("csvFile", "", "Input CSV file path")
	xmlFile := flag.String("xmlFile", "", "Output XML file path")
	columns := flag.String("columns", "", "Comma-separated list of column names")
	orchestra := flag.String("orchestra", "", "Orchestra name")
	flag.Parse()

	// Validate inputs
	if *csvFile == "" || *xmlFile == "" || *columns == "" {
		fmt.Println("Usage: ./csv2xml -csvFile file.csv -xmlFile file.xml -columns \"col1,col2,col3\" -orchestra Orchestra")
		os.Exit(1)
	}

	// Read CSV file
	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		os.Exit(1)
	}

	// Process column headers
	columnNames := strings.Split(*columns, ",")
	if len(columnNames) != len(rows[0]) {
		fmt.Println("Error: Number of provided column names does not match CSV column count.")
		os.Exit(1)
	}

	// Open XML file for writing
	xmlFileHandle, err := os.Create(*xmlFile)
	if err != nil {
		fmt.Println("Error creating XML file:", err)
		os.Exit(1)
	}
	defer xmlFileHandle.Close()

	// Start writing XML
	fmt.Fprintln(xmlFileHandle, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(xmlFileHandle, `<discography orchestra="%s">`, *orchestra)

	// Convert CSV rows into XML
	for _, row := range rows {
		fmt.Fprint(xmlFileHandle, "\n  <track")
		for i, col := range row {
			fmt.Fprintf(xmlFileHandle, ` %s="%s"`, columnNames[i], toSentenceCase(col))
		}
		fmt.Fprint(xmlFileHandle, " />")
	}

	// Close XML structure
	fmt.Fprintln(xmlFileHandle, "\n</discography>")

	fmt.Println("XML file generated successfully:", *xmlFile)
}
