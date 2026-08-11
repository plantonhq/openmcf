package runner

import (
	"strings"
	"testing"
)

// Short compositions pass through untouched -- the cap must not perturb the
// names every existing green lane already uses.
func TestGenerateStackNameShortPassthrough(t *testing.T) {
	got := GenerateStackName("dep-azureresourcegroup-rg", "0a1b2c3d")
	want := "e2e-dep-azureresourcegroup-rg-0a1b2c3d"
	if got != want {
		t.Fatalf("short name changed: got %q want %q", got, want)
	}
}

// The run id contributes at most 8 characters regardless of its length.
func TestGenerateStackNameRunIDShortened(t *testing.T) {
	got := GenerateStackName("dep-x", "0123456789abcdef")
	want := "e2e-dep-x-01234567"
	if got != want {
		t.Fatalf("run id not shortened: got %q want %q", got, want)
	}
}

// Long compositions respect the cap.
func TestGenerateStackNameCapped(t *testing.T) {
	got := GenerateStackName(strings.Repeat("a", 80), "0a1b2c3d")
	if len(got) != maxStackNameLen {
		t.Fatalf("capped name length = %d, want %d (%q)", len(got), maxStackNameLen, got)
	}
}

// The regression that motivated the hash tail: an install profile's
// documents share a long common prefix and differ only in suffixes past the
// cap (manifest-name tails, per-document indexes). A blind prefix cut
// collapsed all four subnet fixture stacks into ONE shared stack, so each
// deploy silently replaced the previous document's resources. Distinct full
// compositions must yield distinct capped names.
func TestGenerateStackNameLongLabelsStayDistinct(t *testing.T) {
	runID := "0a1b2c3d"
	labels := []string{
		"dep-azuresubnet-planton-oss-e2e-azure-fixture-subnet",
		"dep-azuresubnet-planton-oss-e2e-azure-fixture-firewall-subnet-1",
		"dep-azuresubnet-planton-oss-e2e-azure-fixture-gateway-subnet-2",
		"dep-azuresubnet-planton-oss-e2e-azure-fixture-pls-subnet-3",
	}
	seen := map[string]string{}
	for _, label := range labels {
		name := GenerateStackName(label, runID)
		if len(name) > maxStackNameLen {
			t.Fatalf("name for %q exceeds cap: %q (%d chars)", label, name, len(name))
		}
		if prev, dup := seen[name]; dup {
			t.Fatalf("labels %q and %q collide on stack name %q", prev, label, name)
		}
		seen[name] = label
	}
}

// The derivation is deterministic: dependency stacks are keyed by run id so
// every scenario in a run reuses the same stack -- a capped name must come
// out identical on every call.
func TestGenerateStackNameDeterministic(t *testing.T) {
	label := strings.Repeat("dep-long-label-", 6)
	first := GenerateStackName(label, "0a1b2c3d")
	second := GenerateStackName(label, "0a1b2c3d")
	if first != second {
		t.Fatalf("derivation not deterministic: %q vs %q", first, second)
	}
}
