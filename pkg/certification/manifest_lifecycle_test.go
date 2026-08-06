package certification

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plantonhq/planton/internal/manifest"
)

// The manifest-lifecycle certification cases: the offline seams every real
// kind's manifest flows through (load strictly, enforce the envelope,
// validate the spec) must hold for the torture kind's canonical manifest and
// must reject the canonical corruptions with plain-language errors.

// corruptedManifest returns the canonical torture manifest with one exact
// string replaced -- cases corrupt one thing at a time.
func corruptedManifest(t *testing.T, old, new string) []byte {
	t.Helper()
	valid, err := os.ReadFile(TortureManifestPath(t))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(valid), old) {
		t.Fatalf("the canonical torture manifest no longer contains %q -- the corruption would be vacuous", old)
	}
	return []byte(strings.Replace(string(valid), old, new, 1))
}

// Case: the known-good torture manifest loads and validates end to end.
func TestCertify_ValidManifestPassesOfflineValidation(t *testing.T) {
	loaded, err := manifest.LoadManifest(TortureManifestPath(t))
	if err != nil {
		t.Fatalf("the canonical torture manifest failed to load: %v", err)
	}
	if err := manifest.ValidateLoaded(loaded); err != nil {
		t.Fatalf("the canonical torture manifest failed validation: %v", err)
	}
}

// Case: a wrong apiVersion is rejected before anything reaches a server,
// and the error names the exact fix (the plain-language error bar).
func TestCertify_WrongEnvelopeRejectedWithExactFix(t *testing.T) {
	corrupted := corruptedManifest(t, "_test.planton.dev/v1alpha1", "_test.planton.dev/v99")

	loaded, err := manifest.LoadManifestBytes(corrupted, "torture-wrong-envelope.yaml")
	if err == nil {
		err = manifest.ValidateLoaded(loaded)
	}
	if err == nil {
		t.Fatal("a manifest with a wrong apiVersion must be rejected")
	}
	msg := err.Error()
	if !strings.Contains(msg, "_test.planton.dev/v1alpha1") {
		t.Errorf("the rejection must name the correct apiVersion so the fix is in the error; got: %v", err)
	}
}

// Case: a field the schema does not know is rejected at load -- never
// silently dropped (the data-safety law this program was born from).
func TestCertify_UnknownFieldRejectedAtLoad(t *testing.T) {
	valid, err := os.ReadFile(TortureManifestPath(t))
	if err != nil {
		t.Fatal(err)
	}
	corrupted := append([]byte(nil), valid...)
	corrupted = append(corrupted, []byte("  fieldNobodyDeclared: some-value\n")...)

	if _, err := manifest.LoadManifestBytes(corrupted, "torture-unknown-field.yaml"); err == nil {
		t.Fatal("a manifest carrying an unknown field must fail to load, never lose the field silently")
	}
}

// Case: the torture kind keeps its full catalog shape. The kind only proves
// what it exercises -- if a pipeline artifact quietly disappears, every
// certification case that depends on it degrades to vacuous, so the shape
// itself is certified.
func TestCertify_TortureKindKeepsFullCatalogShape(t *testing.T) {
	dir := TortureKindDir(t)
	for _, required := range []string{
		"api.proto",
		"spec.proto",
		"stack_input.proto",
		"stack_outputs.proto",
		"README.md",
		"docs/README.md",
		filepath.Join("presets", "01-default.yaml"),
		filepath.Join("iac", "tf", "main.tf"),
		filepath.Join("iac", "tf", "variables.tf"),
		filepath.Join("iac", "tf", "outputs.tf"),
		filepath.Join("iac", "pulumi", "main.go"),
		filepath.Join("iac", "pulumi", "module", "main.go"),
	} {
		if _, err := os.Stat(filepath.Join(dir, required)); err != nil {
			t.Errorf("torture kind lost part of its catalog shape: %s (%v)", required, err)
		}
	}
}
