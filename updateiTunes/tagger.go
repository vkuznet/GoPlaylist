package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2"
)

func UpdateTags(orchestra, filePath string, track *Track, dryRun bool, verbose int) error {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	defer tag.Close()
	if err != nil {
		if verbose > 2 {
			log.Printf("warning: could not open tag in file %s: %v", filePath, err)
		}
		// Attempt to rewrite the tag from scratch
		return writeNewTag(orchestra, filePath, track, dryRun, verbose)
	}
	defer tag.Close()

	return writeTag(tag, orchestra, filePath, track, dryRun, verbose)
}

func writeNewTag(orchestra, filePath string, track *Track, dryRun bool, verbose int) error {
	tag := id3v2.NewEmptyTag()
	return writeTag(tag, orchestra, filePath, track, dryRun, verbose)
}

func writeTag(tag *id3v2.Tag, orchestra, filePath string, track *Track, dryRun bool, verbose int) error {
	if verbose > 1 {
		log.Printf("Updating %s\n", filePath)
		log.Printf("    new tags: Title=%s Artist=%s Genre=%s Year=%s Album_Artist=%s Composer=%s Author=%s Label=%s",
			track.Name, orchestra, track.Genre, track.Year, track.Vocal, track.Composer, track.Author, track.Label)
	}

	tag.SetTitle(track.Name)
	tag.SetArtist(orchestra)
	tag.SetGenre(track.Genre)

	if len(track.Year) >= 4 {
		tag.SetYear(track.Year[:4])
	}

	tag.AddTextFrame("TPE2", tag.DefaultEncoding(), track.Vocal)    // Album Artist
	tag.AddTextFrame("TCOM", tag.DefaultEncoding(), track.Composer) // Composer
	tag.AddTextFrame("TEXT", tag.DefaultEncoding(), track.Author)   // Author
	tag.AddTextFrame("TPUB", tag.DefaultEncoding(), track.Label)

	if !dryRun {
		if err := tag.Save(); err != nil {
			return fmt.Errorf("error saving tags: file %s, error %w", filePath, err)
		}
		/*
			// Overwrite file with new tag
			mp3File, err := os.OpenFile(filePath, os.O_RDWR, 0666)
			if err != nil {
				return fmt.Errorf("cannot open file to write new tags: %s, error: %w", filePath, err)
			}
			defer mp3File.Close()

			if _, err := tag.WriteTo(mp3File); err != nil {
				return fmt.Errorf("failed to write new tags to file: %s, error: %w", filePath, err)
			}
		*/
	}
	return nil
}
