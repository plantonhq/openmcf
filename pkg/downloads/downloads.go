package downloads

import (
	"fmt"
	"strings"
)

const (
	// BaseURL is the base URL for Planton artifact downloads hosted on Cloudflare R2.
	// All non-CLI release artifacts (Pulumi binaries, Terraform modules, content zips)
	// are published here. CLI binaries remain on GitHub Releases via GoReleaser.
	BaseURL = "https://downloads.planton.dev/releases"
)

// This package is pure URL grammar: every function composes the R2 key shapes
// that the release workflows upload (release.terraform-modules.yaml,
// release.pulumi-modules.yaml). The two sides MUST stay in lockstep — a shape
// change here or in CI without the matching change on the other side does not
// fail loudly; it turns every module fetch into a 404 that downstream code
// masks with a slow git-clone fallback. The unit tests in this package pin the
// shapes, and the release lanes probe one uploaded artifact per run, so drift
// surfaces in CI rather than in the field.
//
// Module artifacts are keyed per served version of a component
// ({component}/{versionDir}) because a component can ship modules for more
// than one served version of its API. Callers derive the version directory
// from the kind registry (crkreflect.ComponentVersionDir), never a literal.

// BuildPulumiDownloadURL constructs the R2 download URL for a Pulumi component binary.
//
// URL format: https://downloads.planton.dev/releases/{version}/modules/pulumi/{component}/{versionDir}_{platform}.gz
// (windows artifacts carry the executable suffix: {versionDir}_{platform}.exe.gz)
//
// Examples (on darwin/arm64):
//
//	BuildPulumiDownloadURL("AwsEcsService", "v1alpha1", "v0.3.50", "darwin_arm64")
//	  -> https://downloads.planton.dev/releases/v0.3.50/modules/pulumi/awsecsservice/v1alpha1_darwin_arm64.gz
func BuildPulumiDownloadURL(component, versionDir, releaseVersion, platform string) string {
	// The release lane gzips "{component}.exe" on windows, so the remote
	// artifact name carries ".exe.gz" there; local cache naming re-adds the
	// ".exe" independently (see pulumibinary.BuildBinaryName).
	suffix := ".gz"
	if strings.HasPrefix(platform, "windows") {
		suffix = ".exe.gz"
	}
	artifact := fmt.Sprintf("%s_%s%s", versionDir, platform, suffix)
	return fmt.Sprintf("%s/%s/modules/pulumi/%s/%s", BaseURL, releaseVersion, strings.ToLower(component), artifact)
}

// BuildTerraformDownloadURL constructs the R2 download URL for a Terraform module zip.
//
// URL format: https://downloads.planton.dev/releases/{version}/modules/terraform/{component}/{versionDir}.zip
//
// Examples:
//
//	BuildTerraformDownloadURL("AwsEcsService", "v1alpha1", "v0.3.50")
//	  -> https://downloads.planton.dev/releases/v0.3.50/modules/terraform/awsecsservice/v1alpha1.zip
func BuildTerraformDownloadURL(component, versionDir, releaseVersion string) string {
	return fmt.Sprintf("%s/%s/modules/terraform/%s/%s.zip", BaseURL, releaseVersion, strings.ToLower(component), versionDir)
}

// BuildDefinitionsDownloadURL constructs the R2 download URL for a file in
// a definitions release (the skill/agent content packaged by
// pkg/skills/defspack: per-skill zips, agent instruction files, and the
// definitions-manifest.json that indexes them with checksums).
//
// URL format: https://downloads.planton.dev/releases/{version}/definitions/{file}
//
// Examples:
//
//	BuildDefinitionsDownloadURL("v0.4.0", "definitions-manifest.json")
//	  -> https://downloads.planton.dev/releases/v0.4.0/definitions/definitions-manifest.json
//	BuildDefinitionsDownloadURL("v0.4.0", "skill-planton.zip")
//	  -> https://downloads.planton.dev/releases/v0.4.0/definitions/skill-planton.zip
func BuildDefinitionsDownloadURL(releaseVersion, file string) string {
	return fmt.Sprintf("%s/%s/definitions/%s", BaseURL, releaseVersion, file)
}
