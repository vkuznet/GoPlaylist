package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type MatchMode string

const (
	Strict MatchMode = "strict"
	Fuzzy  MatchMode = "fuzzy"
)

type MatchResult struct {
	FilePath string
	Track    *Track
}

func MatchTracks(files []string, discography []Track, mode MatchMode, verbose int) []MatchResult {
	var results []MatchResult

	for _, file := range files {
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		name = strings.ToLower(name)

		var bestMatch *Track
		for _, track := range discography {
			matchKey := strings.ToLower(track.Name)
			if mode == Strict && name == matchKey {
				if verbose > 1 {
					log.Printf("strict mode: file '%s' match key '%s' track %+v\n", name, matchKey, track)
				}
				bestMatch = &track
				break
			} else if mode == Fuzzy {
				mfold := fuzzy.MatchNormalizedFold(matchKey, name)
				rfold := fuzzy.RankMatchNormalizedFold(matchKey, name)
				matchStr := strings.Contains(name, matchKey)
				if bestMatch == nil {
					if mfold && rfold < 30 || mfold && matchStr {
						if verbose > 1 {
							log.Printf("fuzzy mode: file '%s' match fold %v rank fold %v track %+v\n", name, mfold, rfold, track)
						}
						bestMatch = &track
					}
				}
			}
		}

		if bestMatch != nil {
			results = append(results, MatchResult{FilePath: file, Track: bestMatch})
		}
	}

	return results
}
