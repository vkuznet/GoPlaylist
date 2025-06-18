# 🎶 updateiTunes

`updateiTunes` is a command-line utility written in Go that syncs your iTunes
music library with detailed metadata from tango (or other genre) discography
XML files. It scans your local music collection, matches tracks by name (with
strict or fuzzy logic), and updates ID3v2 tags with metadata such as composer,
vocalist, year, label, and more.

It’s ideal for tango DJs and music collectors who manage large collections and
want high-quality, consistent metadata based on historical discographies.

---

## Features

* Parse one or more XML discography files with tango metadata
* Match MP3 files using **strict** or **fuzzy** name comparison
* Update ID3v2 tags with fields like Composer, Vocalist, Label, etc.
* Dry-run mode for safe preview
* Optional title correction from discography
* Verbose logging for debugging and transparency

---

## Getting Started

### 1. Build the tool

```sh
go build -o updateiTunes
```

### 2. Prepare your files

* XML files should describe tracks in the following format:

```xml
<discography>
  <track name="Malena" vocal="Raúl Berón" year="1942" genre="Tango" composer="Luciano Demare" author="Homero Manzi" label="RCA Victor"/>
  ...
</discography>
```

* Your music directory should contain `.mp3` files organized in subfolders, e.g.:

```
~/Music/iTunes/iTunes Music/Lucio Demare/Malena.mp3
```

### 3. Run the tool

```sh
./updateiTunes \
  -orchestra "Lucio Demare" \
  -xml "./discographies/Lucio Demare*.xml" \
  -musicDir ./music \
  -matchMode fuzzy \
  -dryRun \
  -fixTitle \
  -verbose 1
```

### CLI Options

| Flag         | Description                                                              |
| ------------ | ------------------------------------------------------------------------ |
| `-xml`       | Glob pattern for XML files (`"*.xml"` by default)                        |
| `-musicDir`  | Path to your iTunes music folder                                         |
| `-orchestra` | Orchestra name to filter subdirectories (e.g., "Carlos Di Sarli")        |
| `-matchMode` | `strict` or `fuzzy` matching based on filename vs. track name            |
| `-dryRun`    | If set, tags won’t be written to files (safe preview mode)               |
| `-fixTitle`  | Replace MP3 title with the XML track name                                |
| `-verbose`   | Increase output verbosity (0 = silent, 1 = summary, 2+ = detailed debug) |

---

## How Matching Works

* **Strict mode** requires exact lowercase match between MP3 filename and discography title.
* **Fuzzy mode** allows approximate matches using `lithammer/fuzzysearch` (normalization and ranking).

---

## Tags Updated

The following ID3v2 frames are updated for each matched MP3:

| Frame ID | Description              | Source Field (XML)      |
| -------- | ------------------------ | ----------------------- |
| `TIT2`   | Title                    | `name`                  |
| `TPE1`   | Lead Performer / Artist  | `-orchestra` flag       |
| `TPE2`   | Orchestra / Album Artist | `vocal`                 |
| `TCON`   | Genre                    | `genre`                 |
| `TYER`   | Year                     | `year` (first 4 digits) |
| `TCOM`   | Composer                 | `composer`              |
| `TEXT`   | Author / Lyricist        | `author`                |
| `TPUB`   | Publisher / Label        | `label`                 |

💡 You can easily extend this tool to write additional tags such as `TRCK`, `TALB`, or user-defined `TXXX` frames if needed.

---

## Example Use Case

You're a tango DJ managing a digital archive. You’ve carefully compiled XML files listing the original composer, vocalist, and publisher of historical recordings. With `updateiTunes`, you can:

* Match those entries to actual MP3 files
* Standardize and enrich tags across your library
* Preview the impact safely before writing changes

---

## Example Discography XML

```xml
<discography>
  <track name="La Cumparsita" vocal="Roberto Díaz" year="1944" genre="Tango" composer="Gerardo Matos Rodríguez" author="Pascual Contursi" label="Odeon"/>
</discography>
```

---

## Notes on ID3v2

* `TYER` (Year) is used for compatibility, though ID3v2.4 prefers `TDRC`.
* `TPE2` is often treated by iTunes as "Album Artist."
* Many apps recognize only a subset of ID3v2 frames.

---

## References

* [ID3v2.3 Tag Spec](https://id3.org/id3v2.3.0)
* [lithammer/fuzzysearch](https://github.com/lithammer/fuzzysearch) — fuzzy match logic
* [bogem/id3v2](https://github.com/bogem/id3v2) — ID3 tag writing in Go
