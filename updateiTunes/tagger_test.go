package main

import (
	"io"
	"log"
	"os"
	"testing"

	id3v2 "github.com/bogem/id3v2"
)

func TestUpdateTagsWithRealMP3(t *testing.T) {
	// Step 1: Open reference MP3 (you must provide this in your testdata dir)
	srcName := "testdata/test.mp3"
	srcFile, err := os.Open(srcName)
	if err != nil {
		t.Fatalf("failed to open test.mp3: %v", err)
	}
	defer srcFile.Close()

	// Step 2: Create temporary copy
	tmpFile, err := os.CreateTemp("", "test_copy_*.mp3")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up after test

	// Copy data from test.mp3 to temp file
	log.Printf("Copy %s to %s, sleep 1sec", srcName, tmpFile.Name())
	if _, err := io.Copy(tmpFile, srcFile); err != nil {
		t.Fatalf("failed to copy MP3 data to temp file: %v", err)
	}
	tmpFile.Close() // close so we can reopen it for updating

	// Step 3: Setup Track data
	track := Track{
		Name:     "UnitTest Track",
		Vocal:    "Test Vocalist",
		Year:     "1935-11-23",
		Genre:    "Tango",
		Composer: "Composer Example",
		Author:   "Author Example",
		Label:    "Label Example",
	}

	// Step 4: Update tags
	err = UpdateTags("Test Orchestra", tmpFile.Name(), &track, false, 3)
	if err != nil {
		t.Fatalf("UpdateTags failed: %v", err)
	}

	// Step 5: Reopen and check updated tags
	tag, err := id3v2.Open(tmpFile.Name(), id3v2.Options{Parse: true})
	defer tag.Close()
	if err != nil {
		t.Fatalf("id3v2.Open fails with %v", err)
	}

	// Step 6: Assertions
	if got := tag.Title(); got != track.Name {
		t.Errorf("expected Title %q, got %q", track.Name, got)
	}
	if got := tag.Artist(); got != "Test Orchestra" {
		t.Errorf("expected Artist %q, got %q", "Test Orchestra", got)
	}
	if got := tag.Genre(); got != track.Genre {
		t.Errorf("expected Genre %q, got %q", track.Genre, got)
	}
	if got := tag.Year(); got != "1935" {
		t.Errorf("expected Year '1935', got %q", got)
	}
}
