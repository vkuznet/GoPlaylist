# updateiTunes
This package provides utility to update iTunes and fix tracks according to
discographies. Here is simple note how to use it:

```
# compile the code 
go build

# run code with your discography file
# here we can use dryRun (to avoid overwriting files)
# different verbose levels provide different details
updateiTunes -orchestra "Lucio Demare" -xml "./discographies/Lucio Demare.xml" -musicDir ./music -matchMode fuzzy -dryRun -verbose 2
```

---
Here is a list of common ID3 tags which may be used, see `tagger.go` codebase

### 🎵 **Common ID3v2 Text Frame Identifiers**

| **Frame ID** | **Meaning**                                  | **Example / Notes**                     |
| ------------ | -------------------------------------------- | --------------------------------------- |
| **TIT2**     | Title/songname/content description           | Track title                             |
| **TPE1**     | Lead performer(s)/Soloist(s)                 | Artist                                  |
| **TPE2**     | Band/orchestra/accompaniment                 | Album artist                            |
| **TALB**     | Album/Movie/Show title                       | Album name                              |
| **TRCK**     | Track number/Position in set                 | `1/10`, `5`                             |
| **TYER**     | Year (ID3v2.3 only)                          | `1995` (replaced by `TDRC` in ID3v2.4)  |
| **TDRC**     | Recording time (ID3v2.4)                     | `2001-01-01` or `2001`                  |
| **TCON**     | Content type                                 | Genre (e.g. `Tango`, `Jazz`)            |
| **TCOM**     | Composer                                     | `Carlos Gardel`                         |
| **TEXT**     | Lyricist/Text writer                         | `Manuel Romero`                         |
| **TPUB**     | Publisher                                    | Label name                              |
| **TENC**     | Encoded by                                   | Software/Person who encoded the track   |
| **TIT1**     | Content group description                    | For classical: "Symphony No. 5"         |
| **TIT3**     | Subtitle/Description refinement              | For example: "Live at Luna Park"        |
| **TKEY**     | Initial key                                  | Musical key: `C#m`, `F#`, etc.          |
| **TLEN**     | Length (milliseconds)                        | Duration                                |
| **TFLT**     | File type                                    | `MPG/3` (MP3), `WAV`, etc.              |
| **TOPE**     | Original artist                              | Useful in cover songs                   |
| **TOAL**     | Original album                               | Album for the original version          |
| **TSRC**     | ISRC (International Standard Recording Code) | Example: `USRC17607839`                 |
| **TSSE**     | Software/Hardware encoding settings          | Useful for debugging encoding issues    |
| **TXXX**     | User-defined text                            | Free-form text with a description field |

---

### 📌 Notes

* **`TXXX`** is useful when no predefined frame fits your need:

  ```go
  tag.AddUserDefinedTextFrame("TXXX", "DJ Notes", "Golden Era classic")
  ```

* **`TYER`** is deprecated in ID3v2.4, replaced by `TDRC`, which supports full dates.

* Many media players (e.g., iTunes, VLC) mostly support a subset: `TIT2`, `TPE1`, `TALB`, `TRCK`, `TYER`, `TCON`, `TCOM`.

---

### Tango discographies use case:

We likely want to use:

* `TIT2` (Title)
* `TPE1` (Lead Performer, e.g., Vocalist)
* `TPE2` (Orchestra / Album Artist)
* `TALB` (Album)
* `TRCK` (Track number)
* `TYER` / `TDRC` (Year)
* `TCON` (Genre: Tango, Milonga, Vals)
* `TCOM` (Composer)
* `TEXT` (Author/Lyricist)
* `TPUB` (Publisher/Label)
