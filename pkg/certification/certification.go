// Package certification is the catalog's launch-certification suite.
//
// The suite drives the synthetic torture kinds (the _test provider) through
// every capability the platform promises about API versioning and data
// safety, using the same seams real kinds flow through. The torture kinds
// exist so these proofs are hermetic and permanent: no cloud, no
// credentials, no dependence on production kind shapes that refactor over
// time.
//
// The suite grows with the machinery -- a capability lands together with its
// certification cases, never later:
//
//   - manifest lifecycle: envelope enforcement, strict loading, validation
//     (this package's initial cases)
//   - version conversion: golden-corpus round-trips once the conversion
//     engine exists
//   - bundle conformance: catalog-bundle round-trips once the bundle exists
//
// Public launch is gated on the full suite passing -- treat a failure here
// as a release blocker, not a flaky test.
package certification

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TortureKindDir returns the torture kind's version directory, resolved from
// this file's location so cases work from any test working directory.
func TortureKindDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller location")
	}
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	dir := filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "_test", "testcloudresourcegeneric", "v1alpha1")
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("torture kind directory missing: %v", err)
	}
	return dir
}

// TortureManifestPath returns the canonical known-good torture manifest (the
// kind's default preset -- one source of truth for "a valid manifest").
func TortureManifestPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(TortureKindDir(t), "presets", "01-default.yaml")
}
