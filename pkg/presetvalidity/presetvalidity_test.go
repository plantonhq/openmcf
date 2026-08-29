package presetvalidity

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot resolves the repository root from this file's location so the
// gate works from any test working directory (including the Bazel sandbox,
// where the catalog source tree is absent -- callers skip there).
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller location")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..")
}

// TestPresetValidityGate is the CI guardrail: the live walk over every
// catalog preset must not introduce a violation outside the baseline or
// leave a stale baseline entry. On failure, fix the preset (the detail names
// the exact rejection a user copying it would hit) or -- for a preset whose
// repair is deliberately routed to its provider's sweep batch -- regenerate
// the baseline with PLANTON_REGEN_PRESET_VALIDITY_BASELINE=1 and justify the
// growth in review.
func TestPresetValidityGate(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "catalog")); err != nil {
		t.Skip("catalog source tree not present (bazel sandbox); runs under go test and the lint.preset-validity lane")
	}

	violations, err := Check(root)
	if err != nil {
		t.Fatalf("preset walk: %v", err)
	}

	_, thisFile, _, _ := runtime.Caller(0)
	baselinePath := filepath.Join(filepath.Dir(thisFile), "baseline.yaml")
	if os.Getenv("PLANTON_REGEN_PRESET_VALIDITY_BASELINE") == "1" {
		if err := WriteBaseline(baselinePath, violations); err != nil {
			t.Fatalf("write baseline: %v", err)
		}
		t.Logf("baseline regenerated -- review the diff before committing")
		return
	}

	baseline, err := LoadBaseline(baselinePath)
	if err != nil {
		t.Fatalf("load baseline: %v", err)
	}
	res := Gate(violations, baseline)
	for _, v := range res.NewViolations {
		t.Errorf("preset shipped invalid: %s -- %s", v.ID(), v.Detail)
	}
	for _, id := range res.StaleEntries {
		t.Errorf("stale baseline entry (no longer a violation): %s -- remove it from baseline.yaml", id)
	}
}

// torturePreset reads the torture kind's default preset -- the repository's
// canonical known-good manifest -- as the hermetic green against which each
// red below is one deliberate defect. A gate that cannot fail teaches false
// confidence.
func torturePreset(t *testing.T) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "catalog", "_test", "testcloudresourcegeneric", "presets", "01-default.yaml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("torture preset not present: %v", err)
	}
	return string(content)
}

// TestCheckPreset_HermeticFixture proves the rule fires against deliberately
// broken presets and stays quiet on the canonical valid one.
func TestCheckPreset_HermeticFixture(t *testing.T) {
	green := torturePreset(t)

	if vs := CheckPreset("catalog/_test/x/presets/01-green.yaml", []byte(green)); len(vs) != 0 {
		t.Errorf("the canonical torture preset must pass, got %v", vs)
	}

	// An angle-bracket placeholder where the schema demands a real shape:
	// the exact defect class this gate exists to make unshippable.
	brokenValue := strings.Replace(green, "int32Field: 7", `int32Field: "<replace-me>"`, 1)
	if vs := CheckPreset("catalog/_test/x/presets/02-broken.yaml", []byte(brokenValue)); len(vs) != 1 || vs[0].Rule != RuleInvalidPreset {
		t.Errorf("expected one %s violation for a type-broken placeholder, got %v", RuleInvalidPreset, vs)
	}

	// A wrong envelope: the manifest declares an apiVersion its kind does
	// not serve. The validator's envelope contract catches it here.
	brokenEnvelope := strings.Replace(green, "apiVersion: _test.planton.dev/v1alpha2", "apiVersion: _test.planton.dev/v999", 1)
	if vs := CheckPreset("catalog/_test/x/presets/03-envelope.yaml", []byte(brokenEnvelope)); len(vs) != 1 || vs[0].Rule != RuleInvalidPreset {
		t.Errorf("expected one %s violation for a wrong envelope, got %v", RuleInvalidPreset, vs)
	}

	// A required field removed: the load succeeds, validation rejects.
	brokenRequired := strings.Replace(green, "  requiredRef:\n    value: literal-required-value\n", "", 1)
	if vs := CheckPreset("catalog/_test/x/presets/04-required.yaml", []byte(brokenRequired)); len(vs) != 1 || vs[0].Rule != RuleInvalidPreset {
		t.Errorf("expected one %s violation for a missing required field, got %v", RuleInvalidPreset, vs)
	}
}

// TestGate mirrors the anatomy and catalogpage gate semantics: new drift
// fails, baselined drift passes, a fixed entry left in the baseline fails
// as stale.
func TestGate(t *testing.T) {
	v := Violation{Path: "catalog/aws/x/presets/01-a.yaml", Rule: RuleInvalidPreset}
	if res := Gate([]Violation{v}, map[string]bool{}); res.OK() || len(res.NewViolations) != 1 {
		t.Errorf("expected new drift to be detected, got %+v", res)
	}
	if res := Gate([]Violation{v}, map[string]bool{v.ID(): true}); !res.OK() {
		t.Errorf("expected baselined drift to pass, got %+v", res)
	}
	if res := Gate(nil, map[string]bool{"catalog/aws/gone/presets/01-a.yaml:invalid-preset": true}); res.OK() || len(res.StaleEntries) != 1 {
		t.Errorf("expected a stale entry to be detected, got %+v", res)
	}
	if res := Gate([]Violation{v, v}, map[string]bool{v.ID(): true}); !res.OK() {
		t.Errorf("expected duplicate-rule collapse, got %+v", res)
	}
}
