package conversion

import (
	"os"
	"path/filepath"
	"testing"

	"io/fs"

	"github.com/plantonhq/planton/pkg/conversion/embedded"
)

// The embedded mirror must be byte-identical to the authored specs beside
// the protos -- a stale mirror would ship a released CLI whose offline
// upgrades disagree with the catalog. Regenerate with
// `make generate-conversion-registry`.
func TestEmbeddedSpecsMirrorTheAuthoredSpecs(t *testing.T) {
	base := providerBaseDir(t)
	authored, err := SpecFiles(os.DirFS(base))
	if err != nil {
		t.Fatal(err)
	}
	if len(authored) == 0 {
		t.Fatal("no authored conversion specs found -- the drift gate is broken")
	}

	specsFS, err := embedded.SpecsFS()
	if err != nil {
		t.Fatal(err)
	}
	mirrored, err := SpecFiles(specsFS)
	if err != nil {
		t.Fatal(err)
	}

	authoredSet := map[string]bool{}
	for _, file := range authored {
		authoredSet[file] = true
		want, err := os.ReadFile(filepath.Join(base, file))
		if err != nil {
			t.Fatal(err)
		}
		got, err := fs.ReadFile(specsFS, file)
		if err != nil {
			t.Errorf("%s is authored but missing from the embedded mirror -- run `make generate-conversion-registry`", file)
			continue
		}
		if string(got) != string(want) {
			t.Errorf("%s differs between the authored spec and the embedded mirror -- run `make generate-conversion-registry`", file)
		}
	}
	for _, file := range mirrored {
		if !authoredSet[file] {
			t.Errorf("%s exists in the embedded mirror but is no longer authored -- run `make generate-conversion-registry`", file)
		}
	}
}
