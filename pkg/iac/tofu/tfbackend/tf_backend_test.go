package tfbackend

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plantonhq/planton/shared/iac/terraform"
)

// The no-body form must stay byte-identical to the historical output: every existing
// caller writes `backend "<type>" {}` and downstream tooling parses that shape.
func TestWriteBackendFile_NoBodyIsByteIdenticalToHistoricalForm(t *testing.T) {
	dir := t.TempDir()
	if err := WriteBackendFile(dir, terraform.TerraformBackendType_s3); err != nil {
		t.Fatalf("WriteBackendFile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "backend.tf"))
	if err != nil {
		t.Fatalf("reading backend.tf: %v", err)
	}
	want := `terraform {
  backend "s3" {}
}
`
	if string(got) != want {
		t.Fatalf("backend.tf = %q, want %q", got, want)
	}
}

// The remote backend's workspaces block is part of the DECLARATION -- an HCL block
// cannot ride -backend-config flags. This is the exact shape proven against a live
// TFE-protocol service.
func TestWriteBackendFile_RemoteWithWorkspacesBody(t *testing.T) {
	dir := t.TempDir()
	if err := WriteBackendFile(dir, terraform.TerraformBackendType_remote, WorkspacesNameBody("acme-prod-bucket-x")...); err != nil {
		t.Fatalf("WriteBackendFile: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "backend.tf"))
	if err != nil {
		t.Fatalf("reading backend.tf: %v", err)
	}
	want := `terraform {
  backend "remote" {
    workspaces {
      name = "acme-prod-bucket-x"
    }
  }
}
`
	if string(got) != want {
		t.Fatalf("backend.tf = %q, want %q", got, want)
	}
}

// The enum name IS the HCL backend name; the parser must cover every declared value
// and reject unknown words with the unspecified sentinel.
func TestBackendTypeFromString_CoversEveryEnumValue(t *testing.T) {
	for value, name := range terraform.TerraformBackendType_name {
		backendType := terraform.TerraformBackendType(value)
		if backendType == terraform.TerraformBackendType_terraform_backend_type_unspecified {
			continue
		}
		if got := BackendTypeFromString(name); got != backendType {
			t.Errorf("BackendTypeFromString(%q) = %v, want %v -- a new enum value needs its parser arm", name, got, backendType)
		}
	}
	if got := BackendTypeFromString("consul"); got != terraform.TerraformBackendType_terraform_backend_type_unspecified {
		t.Errorf("BackendTypeFromString(consul) = %v, want unspecified", got)
	}
}
