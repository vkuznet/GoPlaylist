package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseXML(t *testing.T) {
	xmlData := `
<discography orchestra="Carlos Gardel">
    <track name="Tomo y obligo" vocal="Charlo" year="1931-10-09" genre="Tango" composer="Carlos Gardel" author="Manuel Romero"/>
</discography>`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.xml")
	if err := os.WriteFile(tmpFile, []byte(xmlData), 0644); err != nil {
		t.Fatalf("failed to write test XML: %v", err)
	}

	tracks, err := ParseXML(tmpFile)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}
	if len(tracks) != 1 {
		t.Fatalf("expected 1 track, got %d", len(tracks))
	}
	if tracks[0].Name != "Tomo y obligo" {
		t.Errorf("expected track name 'Tomo y obligo', got %s", tracks[0].Name)
	}
}

