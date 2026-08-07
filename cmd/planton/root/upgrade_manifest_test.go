package root

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func tortureConversionsDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller location")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..",
		"catalog", "_test", "testcloudresourcegeneric", "conversions")
}

// The offline upgrade end to end THROUGH THE EMBED: an old-version manifest
// converts to the served version, losses are reported, and the output is a
// valid manifest.
func TestUpgradeManifestFile(t *testing.T) {
	input := filepath.Join(tortureConversionsDir(t), "testdata", "full-shape", "input.yaml")

	out, losses, err := upgradeManifestFile(input)
	if err != nil {
		t.Fatalf("upgrade failed: %v", err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(out, &doc); err != nil {
		t.Fatalf("output is not valid YAML: %v", err)
	}
	if doc["apiVersion"] != "_test.planton.dev/v1alpha2" {
		t.Errorf("output apiVersion = %v, want the served version", doc["apiVersion"])
	}
	if len(losses) != 1 || losses[0].Path != "spec.stringNoDefault" {
		t.Errorf("expected the one declared loss, got %+v", losses)
	}
}

// A manifest already at the served version is refused with a friendly
// explanation, not converted into itself.
func TestUpgradeManifestFile_AlreadyServed(t *testing.T) {
	current := filepath.Join(tortureConversionsDir(t), "testdata", "full-shape", "expected.yaml")

	_, _, err := upgradeManifestFile(current)
	if err == nil {
		t.Fatal("upgrading an already-current manifest must be refused")
	}
	if !strings.Contains(err.Error(), "already at") {
		t.Errorf("the refusal must say the manifest is already current; got: %v", err)
	}
}
