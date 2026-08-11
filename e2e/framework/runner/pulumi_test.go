package runner

import (
	"strings"
	"testing"
)

// TestGenerateStackNameShortNamesUnchanged pins the common case: a name inside
// the cap keeps the exact e2e-{label}-{shortRunID} shape with no digest.
func TestGenerateStackNameShortNamesUnchanged(t *testing.T) {
	got := GenerateStackName("awsvpc-minimal", "87e9a9f8")
	want := "e2e-awsvpc-minimal-87e9a9f8"
	if got != want {
		t.Fatalf("short name changed: got %q, want %q", got, want)
	}
}

// TestGenerateStackNameLongLabelsStayUnique pins the live-caught failure mode:
// the three AwsSubnet install-profile instances (and the run id) exceed the
// length cap, and a plain prefix truncation collapsed all three into ONE stack
// name — each successive `pulumi up` then silently replaced the previous
// instance's cloud resource. The digest tail must keep them distinct while
// honoring the cap.
func TestGenerateStackNameLongLabelsStayUnique(t *testing.T) {
	labels := []string{
		"dep-awssubnet-planton-oss-e2e-awssubnet-prereq-a",
		"dep-awssubnet-planton-oss-e2e-awssubnet-prereq-b",
		"dep-awssubnet-planton-oss-e2e-awssubnet-prereq-c",
	}

	seen := make(map[string]string, len(labels))
	for _, label := range labels {
		name := GenerateStackName(label, "87e9a9f8")
		if len(name) > maxStackNameLen {
			t.Fatalf("stack name for %q exceeds cap: %q (%d chars)", label, name, len(name))
		}
		if prior, dup := seen[name]; dup {
			t.Fatalf("stack name collision: labels %q and %q both produced %q", prior, label, name)
		}
		seen[name] = label

		if !strings.HasPrefix(name, "e2e-dep-awssubnet-") {
			t.Fatalf("truncated name lost its readable prefix: %q", name)
		}
	}
}

// TestGenerateStackNameDeterministic pins that capping is a pure function of
// its inputs — dependency DESTROY recomputes nothing, it reuses the recorded
// deploy-time name, but determinism keeps re-runs and debugging sane.
func TestGenerateStackNameDeterministic(t *testing.T) {
	label := "dep-awssubnet-planton-oss-e2e-awssubnet-prereq-a"
	first := GenerateStackName(label, "87e9a9f8")
	second := GenerateStackName(label, "87e9a9f8")
	if first != second {
		t.Fatalf("same inputs produced different names: %q vs %q", first, second)
	}
	if other := GenerateStackName(label, "deadbeef"); other == first {
		t.Fatalf("different run ids produced the same name: %q", first)
	}
}
