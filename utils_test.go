package main

import (
	"testing"
)

// TestGetorchestra
func Testorchestra(t *testing.T) {
	ainput := "Orquesta Tipica Victor (dir. Adolfo Carabelli)"
	expect := "Orquesta Tipica Victor dir Adolfo Carabelli"
	orchestra := getOrchestra(ainput, nil)
	if orchestra != expect {
		t.Errorf("Fail to parse input '%s': expect='%s' received='%s' ", ainput, expect, orchestra)
	}
	// test case when discography exist
	discography := &Discography{Orchestra: "OTV"}
	orchestra = getOrchestra(ainput, discography)
	if orchestra != "OTV" {
		t.Errorf("Fail to parse input '%s': expect='%s' received='%s' ", ainput, "OTV", orchestra)
	}

}

func TestInUtils(t *testing.T) {
	trk1 := Track{Name: "Name", Year: "Year", Orchestra: "Orchestra1"}
	trk2 := Track{Name: "Name", Year: "Year", Orchestra: "Orchestra2"}
	trk3 := Track{Name: "Name", Year: "Year", Orchestra: "Orchestra3"}
	trackList := []Track{trk1, trk2, trk3}
	for _, trk := range trackList {
		res := inList(trk, trackList)
		if !res {
			t.Errorf("unable to find item '%s' in a list %s", trk.String(), trackList)
		}
	}
}
