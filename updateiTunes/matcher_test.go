package main

import (
	"testing"
)

func TestStrictMatch(t *testing.T) {
	xmlTracks := []Track{
		{Name: "La Muchachada Del Centro", Vocal: "Ernesto Fama", Year: "1932"},
	}
	files := []string{"La Muchachada Del Centro.mp3"}

	results := MatchTracks(files, xmlTracks, Strict, 0)
	if len(results) != 1 {
		t.Errorf("expected 1 match, got %d", len(results))
	}
}

func TestFuzzyMatch(t *testing.T) {
	xmlTracks := []Track{
		{Name: "La Muchachada Del Centro", Vocal: "Ernesto Fama"},
	}
	files := []string{"La Muchachada Del Centro - Francisco Canaro.mp3"}

	results := MatchTracks(files, xmlTracks, Fuzzy, 0)
	if len(results) != 1 {
		t.Errorf("expected 1 fuzzy match, got %d", len(results))
	}
}
