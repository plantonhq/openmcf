package platformversion

import (
	"fmt"
	"strings"
	"testing"
)

func TestMinimumSupportedIsAReleaseVersion(t *testing.T) {
	if !releaseForm.MatchString(MinimumSupported) {
		t.Fatalf("MinimumSupported %q must itself be a full release version", MinimumSupported)
	}
	if strings.Contains(MinimumSupported, "-") {
		t.Fatalf("MinimumSupported %q must name a release, never a pre-release", MinimumSupported)
	}
}

// bump returns the floor with its patch number moved by delta, so the cases
// below stay true when the floor moves.
func bump(t *testing.T, delta int) string {
	t.Helper()
	var major, minor, patch int
	if _, err := fmt.Sscanf(MinimumSupported, "v%d.%d.%d", &major, &minor, &patch); err != nil {
		t.Fatalf("parsing MinimumSupported %q: %v", MinimumSupported, err)
	}
	return fmt.Sprintf("v%d.%d.%d", major, minor, patch+delta)
}

func TestCheck_Supported(t *testing.T) {
	for _, v := range []string{
		MinimumSupported,        // the floor itself
		bump(t, 1),              // the next patch
		"v0.1.0",                // a newer minor
		"v1.0.0",                // a newer major
		bump(t, 1) + "-rc.1",    // a pre-release of a NEWER version is above the floor
		MinimumSupported + "+b", // build metadata never affects precedence
	} {
		got := Check(v)
		if !got.Supported || got.Reason != ReasonSupported || got.Message != "" {
			t.Errorf("Check(%q) = %+v, want Supported with no message", v, got)
		}
	}
}

func TestCheck_BelowMinimum(t *testing.T) {
	for _, v := range []string{
		bump(t, -1),
		"v0.0.41-selfhosted-preview", // the retired preview line sorts below every real release
		"v0.0.40-selfhosted-preview",
		MinimumSupported + "-rc.1", // a pre-release of the floor version was cut before the floor
	} {
		got := Check(v)
		if got.Supported || got.Reason != ReasonBelowOperatorMinimum {
			t.Errorf("Check(%q) = %+v, want BelowOperatorMinimum", v, got)
			continue
		}
		for _, want := range []string{v, MinimumSupported, "Nothing running was changed", "Set spec.version to " + MinimumSupported} {
			if !strings.Contains(got.Message, want) {
				t.Errorf("Check(%q) message %q lacks %q", v, got.Message, want)
			}
		}
	}
}

func TestCheck_Unreadable(t *testing.T) {
	for _, v := range []string{
		"",
		"local",
		"latest",
		MinimumSupported[1:], // no v prefix
		"v1.2",               // shorthand the semver library would accept; the CRD does not
		MinimumSupported + " ",
	} {
		got := Check(v)
		if got.Supported || got.Reason != ReasonUnreadable {
			t.Errorf("Check(%q) = %+v, want Unreadable", v, got)
			continue
		}
		if !strings.Contains(got.Message, "image.tag") || !strings.Contains(got.Message, MinimumSupported) {
			t.Errorf("Check(%q) message must name the override and the floor: %q", v, got.Message)
		}
	}
}
