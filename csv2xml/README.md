# CSV to XML Converter

## Overview
This program converts a given CSV file into an XML file based on user-defined parameters. It allows specifying column names and an orchestra name for the XML output.

## Usage
Run the program with the following command:

```sh
./csv2xml -csvFile file.csv -xmlFile file.xml -columns "name,year,genre" -orchestra "Orchestra Name"
```

### Arguments:
- `-csvFile`: Path to the input CSV file (Required)
- `-xmlFile`: Path to the output XML file (Required)
- `-columns`: Column names for the CSV fields, specified in a comma-separated string (Default: `name,year,genre`)
- `-orchestra`: Orchestra name to be included in the XML output (Default: `Unknown Orchestra`)

## Example

### Input CSV (`file.csv`):
```csv
track_name1,year1,tango
track_name2,year2,vals
```

### Command:
```sh
./csv2xml -csvFile file.csv -xmlFile file.xml -columns "name,year,genre" -orchestra "Orchestra Name"
```

### Generated XML (`file.xml`):
```xml
<?xml version="1.0" encoding="UTF-8"?>
<discography orchestra="Orchestra Name">
  <track name="track_name1" year="year1" genre="tango"/>
  <track name="track_name2" year="year2" genre="vals"/>
</discography>
```

## Dependencies
- Go (1.16 or later)

## Compilation
To build the program, run:
```sh
go build -o csv2xml main.go
```

## License
This project is licensed under the MIT License.
