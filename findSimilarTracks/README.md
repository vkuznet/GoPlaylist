# README: findSimilarTracks

## Description
`findSimilarTracks` is a Go program designed to process multiple XML files containing information about musical tracks. The program identifies tracks with the same name but different metadata (e.g., `year` or `orchestra`), and generates a consolidated XML file containing these "similar" tracks. 

The output XML file contains the filtered tracks sorted by their names.

---

## Input XML File Example
Here is an example of an input XML file (`input1.xml`):
```xml
<?xml version="1.0" encoding="UTF-8"?>
<discography orchestra="Orchestra1" source="tango.info">
    <track name="Name1" vocal="Vocal" year="1942-10-27" genre="tango" composer="Composer1" author="Author1" duration="2:42" popularity="4"/>
    <track name="Name2" vocal="Vocal" year="1942-11-04" genre="tango" composer="Composer2" author="Author2" duration="2:48" popularity="3"/>
</discography>
```

And another input file (`input2.xml`):
```xml
<?xml version="1.0" encoding="UTF-8"?>
<discography orchestra="Orchestra2" source="tango.info">
    <track name="Name1" vocal="Vocal" year="1950-08-12" genre="tango" composer="Composer1" author="Author1" duration="3:00" popularity="5"/>
    <track name="Name3" year="1924" genre="tango" composer="Composer3" author="Author3" duration="2:36" popularity="5"/>
</discography>
```

---

## Output XML Example
When the program processes the input XML files, it generates an output file (`output.xml`) containing only the similar tracks:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<tracks>
  <track name="Name1" year="1942-10-27" orchestra="Orchestra1" genre="tango"/>
  <track name="Name1" year="1950-08-12" orchestra="Orchestra2" genre="tango"/>
</tracks>
```

---

## How It Works
1. The program takes an input file pattern (e.g., `/path/*.xml`) and an output file name (e.g., `output.xml`).
2. It scans all XML files matching the pattern, recursively if necessary, and reads their contents.
3. It identifies tracks with the same name but different `year` or `orchestra` attributes.
4. The program removes duplicate entries and sorts the output tracks alphabetically by name.
5. The result is saved in the specified output file.

---

## Usage
To compile and run the program:
1. **Build the executable**:
   ```bash
   go build -o findSimilarTracks
   ```

2. **Run the program**:
   ```bash
   ./findSimilarTracks -input "/path/*.xml" -output output.xml

   # if you want final statistics about tracks use -stats flag
   ./findSimilarTracks -input "/path/*.xml" -output output.xml -stats
   ```

   Replace `/path/*.xml` with the directory and pattern of your input files.

3. **Output**:
   The filtered and sorted tracks will be saved to `output.xml`.

---

## Notes
- Input files must conform to the provided XML format.
- The program handles cases where the orchestra is missing by defaulting to the discography-level `orchestra` attribute.
- If no similar tracks are found, the output file will be empty.

---

Enjoy using `findSimilarTracks` to manage and process your music discography data!
