package tofuzip

import (
	"strings"
	"testing"

	"github.com/plantonhq/planton/internal/cli/version"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/downloads"
)

// Pins the wiring, not the shape (pkg/downloads owns the shape): the URL this
// package builds must carry the version segment the kind registry declares —
// the drift that once broke every released download was exactly this segment
// going stale against CI.
func TestBuildDownloadURLDerivesVersionSegmentFromRegistry(t *testing.T) {
	const component = "AwsS3Bucket"
	const release = "v0.3.50"

	versionDir, err := crkreflect.ComponentVersionDir(component)
	if err != nil {
		t.Fatalf("registry has no version for %s: %v", component, err)
	}

	got, err := BuildDownloadURL(component, release)
	if err != nil {
		t.Fatalf("BuildDownloadURL() error: %v", err)
	}

	want := downloads.BuildTerraformDownloadURL(component, versionDir, release)
	if got != want {
		t.Errorf("BuildDownloadURL() = %q, want %q", got, want)
	}
	if !strings.Contains(got, "/"+versionDir+".zip") {
		t.Errorf("BuildDownloadURL() = %q does not carry the registry version segment %q", got, versionDir)
	}
}

func TestBuildDownloadURLUnknownKindFailsPlainly(t *testing.T) {
	_, err := BuildDownloadURL("NoSuchKind", "v0.3.50")
	if err == nil {
		t.Fatal("BuildDownloadURL() with an unknown kind must fail instead of composing a URL that 404s")
	}
	if !strings.Contains(err.Error(), "NoSuchKind") {
		t.Errorf("error %q does not name the unresolvable kind", err.Error())
	}
}

// Dev builds have no released artifacts, so the zip fast path must stay off
// and resolution must fall through to staging.
func TestCanUseZipModeIsOffForDevBuilds(t *testing.T) {
	saved := version.Version
	defer func() { version.Version = saved }()

	version.Version = ""
	if CanUseZipMode() {
		t.Error("CanUseZipMode() must be false for an unset version")
	}
	version.Version = version.DefaultVersion
	if CanUseZipMode() {
		t.Error("CanUseZipMode() must be false for the dev default version")
	}
	version.Version = "v0.3.50"
	if !CanUseZipMode() {
		t.Error("CanUseZipMode() must be true for a released version")
	}
}
