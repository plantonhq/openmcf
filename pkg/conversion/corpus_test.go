package conversion

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"sigs.k8s.io/yaml"

	// The corpus converts between BOTH torture-kind versions; the old
	// version's package is linked explicitly so its descriptors resolve (only
	// the served version reaches the binary through the kind registry).
	_ "github.com/plantonhq/planton/apis/dev/planton/provider/_test/testcloudresourcegeneric/v1alpha1"
)

// The golden corpus: every conversion spec in the catalog must carry fixture
// pairs (conversions/testdata/<case>/{input,expected}.yaml), and every
// fixture must satisfy two laws:
//
//  1. upgrade(input) == expected, byte-for-byte at the document level.
//  2. downgrade(upgrade(input)) == input minus the upgrade's DECLARED losses
//     -- the round-trip law: loss is legal only when the spec says so.
//
// Both the Go engine here and the platform's Java engine gate on this same
// corpus; an engine disagreement is a CI failure, not a runtime surprise.

func providerBaseDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller location")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "apis", "dev", "planton", "provider")
}

func loadDoc(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing %s: %v", path, err)
	}
	return doc
}

func TestGoldenCorpus(t *testing.T) {
	base := providerBaseDir(t)
	fsys := os.DirFS(base)
	specFiles, err := SpecFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if len(specFiles) == 0 {
		t.Fatal("no conversion specs found -- the corpus gate is broken")
	}

	for _, specFile := range specFiles {
		spec, err := LoadSpec(filepath.Join(base, specFile))
		if err != nil {
			t.Errorf("%s: %v", specFile, err)
			continue
		}

		fixturesDir := filepath.Join(base, filepath.Dir(specFile), "testdata")
		entries, err := os.ReadDir(fixturesDir)
		if err != nil || len(entries) == 0 {
			t.Errorf("%s has no golden fixtures -- every conversion spec ships its corpus (conversions/testdata/<case>/{input,expected}.yaml)", specFile)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			caseName := entry.Name()
			t.Run(spec.Kind+"/"+spec.From+"_to_"+spec.To+"/"+caseName, func(t *testing.T) {
				input := loadDoc(t, filepath.Join(fixturesDir, caseName, "input.yaml"))
				expected := loadDoc(t, filepath.Join(fixturesDir, caseName, "expected.yaml"))

				upgraded, losses, err := Apply(spec, Upgrade, input)
				if err != nil {
					t.Fatalf("upgrade failed: %v", err)
				}
				if !reflect.DeepEqual(upgraded, expected) {
					t.Errorf("upgrade(input) != expected\ngot:      %#v\nexpected: %#v", upgraded, expected)
				}

				// Round-trip law: the downgrade restores the input except for
				// losses the UPGRADE declared.
				roundTripped, _, err := Apply(spec, Downgrade, upgraded)
				if err != nil {
					t.Fatalf("downgrade failed: %v", err)
				}
				pruned := deepCopy(input).(map[string]any)
				for _, loss := range losses {
					deletePath(pruned, loss.Path)
				}
				if !reflect.DeepEqual(roundTripped, pruned) {
					t.Errorf("round-trip broke: downgrade(upgrade(input)) != input minus declared losses\ngot:      %#v\nexpected: %#v", roundTripped, pruned)
				}
			})
		}
	}
}
