package main

import "testing"

// TestGetArtist
func TestArtist(t *testing.T) {
	ainput := "Orquesta Tipica Victor (dir. Adolfo Carabelli)"
	expect := "Orquesta Tipica Victor dir Adolfo Carabelli"
	artist := getArtist(ainput, nil)
	if artist != expect {
		t.Errorf("Fail to parse input '%s': expect='%s' received='%s' ", ainput, expect, artist)
	}
	// test case when discogrpahy exist
	discography := &Discography{Orchestra: "OTV"}
	artist = getArtist(ainput, discography)
	if artist != "OTV" {
		t.Errorf("Fail to parse input '%s': expect='%s' received='%s' ", ainput, "OTV", artist)
	}

}
