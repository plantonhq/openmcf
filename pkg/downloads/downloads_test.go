package downloads

import "testing"

// These tests pin the R2 key shapes to what the release workflows upload
// (release.terraform-modules.yaml, release.pulumi-modules.yaml). The expected
// strings are deliberate literals: if either side changes its shape, exactly
// one of the two must fail — that failure is the contract doing its job.
// Update the workflows and these literals together, never one alone.

func TestBuildTerraformDownloadURL(t *testing.T) {
	tests := []struct {
		name       string
		component  string
		versionDir string
		release    string
		want       string
	}{
		{
			name:       "canonical component",
			component:  "AwsEcsService",
			versionDir: "v1alpha1",
			release:    "v0.3.50",
			want:       "https://downloads.planton.dev/releases/v0.3.50/modules/terraform/awsecsservice/v1alpha1.zip",
		},
		{
			name:       "graduated version segment",
			component:  "AwsS3Bucket",
			versionDir: "v1",
			release:    "v0.4.0",
			want:       "https://downloads.planton.dev/releases/v0.4.0/modules/terraform/awss3bucket/v1.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildTerraformDownloadURL(tt.component, tt.versionDir, tt.release); got != tt.want {
				t.Errorf("BuildTerraformDownloadURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildPulumiDownloadURL(t *testing.T) {
	tests := []struct {
		name       string
		component  string
		versionDir string
		release    string
		platform   string
		want       string
	}{
		{
			name:       "darwin arm64",
			component:  "AwsEcsService",
			versionDir: "v1alpha1",
			release:    "v0.3.50",
			platform:   "darwin_arm64",
			want:       "https://downloads.planton.dev/releases/v0.3.50/modules/pulumi/awsecsservice/v1alpha1_darwin_arm64.gz",
		},
		{
			name:       "linux amd64",
			component:  "AwsS3Bucket",
			versionDir: "v1alpha1",
			release:    "v0.3.50",
			platform:   "linux_amd64",
			want:       "https://downloads.planton.dev/releases/v0.3.50/modules/pulumi/awss3bucket/v1alpha1_linux_amd64.gz",
		},
		{
			// The release lane gzips "{component}.exe" on windows, so the
			// remote artifact carries ".exe.gz" — pin it so the windows fast
			// path cannot regress to the extensionless shape.
			name:       "windows carries the executable suffix",
			component:  "AwsEcsService",
			versionDir: "v1alpha1",
			release:    "v0.3.50",
			platform:   "windows_amd64",
			want:       "https://downloads.planton.dev/releases/v0.3.50/modules/pulumi/awsecsservice/v1alpha1_windows_amd64.exe.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildPulumiDownloadURL(tt.component, tt.versionDir, tt.release, tt.platform); got != tt.want {
				t.Errorf("BuildPulumiDownloadURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildDefinitionsDownloadURL(t *testing.T) {
	want := "https://downloads.planton.dev/releases/v0.4.0/definitions/definitions-manifest.json"
	if got := BuildDefinitionsDownloadURL("v0.4.0", "definitions-manifest.json"); got != want {
		t.Errorf("BuildDefinitionsDownloadURL() = %q, want %q", got, want)
	}
}
