package main

import (
	"strings"
	"testing"
)

// TestParseWashes pins the -washes syntax: separators no CSS color uses, a
// label color chosen from the paper, and a plain sentence for every way the
// flag can be malformed.
func TestParseWashes(t *testing.T) {
	ws, err := parseWashes("compute light=#ffffff/rgba(200,120,40,.11); compute dark=#121212/rgba(230,160,90,.18)")
	if err != nil {
		t.Fatal(err)
	}
	if len(ws) != 2 {
		t.Fatalf("expected 2 washes, got %d", len(ws))
	}
	if ws[0].name != "compute light" || ws[0].paper != "#ffffff" || ws[0].tint != "rgba(200,120,40,.11)" || ws[0].label != "#444" {
		t.Fatalf("light wash parsed wrong: %+v", ws[0])
	}
	if ws[1].name != "compute dark" || ws[1].paper != "#121212" || ws[1].tint != "rgba(230,160,90,.18)" || ws[1].label != "#bbb" {
		t.Fatalf("dark wash parsed wrong: %+v", ws[1])
	}

	for _, bad := range []string{"", ";", "light", "light=#fff", "=#fff/red", "light=/red", "light=#fff/"} {
		if _, err := parseWashes(bad); err == nil || !strings.Contains(err.Error(), "-washes") {
			t.Fatalf("%q must fail with a sentence naming the flag, got %v", bad, err)
		}
	}
}

func TestParseSizes(t *testing.T) {
	px, err := parseSizes("18, 26,34")
	if err != nil || len(px) != 3 || px[0] != 18 || px[2] != 34 {
		t.Fatalf("got %v, %v", px, err)
	}
	if _, err := parseSizes("18,x"); err == nil {
		t.Fatal("a non-numeric size must fail")
	}
}
