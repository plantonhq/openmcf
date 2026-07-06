package infrachart

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestDiscoverCharts_WithoutAllReturnsArgsVerbatim(t *testing.T) {
	dirs, err := DiscoverCharts([]string{"a", "b"}, false)
	if err != nil {
		t.Fatalf("DiscoverCharts: %v", err)
	}
	if !reflect.DeepEqual(dirs, []string{"a", "b"}) {
		t.Fatalf("got %v, want args verbatim", dirs)
	}
}

func TestDiscoverCharts_WalksRootsForCharts(t *testing.T) {
	dirs, err := DiscoverCharts([]string{"testdata"}, true)
	if err != nil {
		t.Fatalf("DiscoverCharts: %v", err)
	}
	want := []string{
		filepath.Join("testdata", "broken-chart"),
		filepath.Join("testdata", "valid-chart"),
	}
	if !reflect.DeepEqual(dirs, want) {
		t.Fatalf("got %v, want %v (sorted, one entry per chart)", dirs, want)
	}
}
