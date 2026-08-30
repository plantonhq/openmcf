package backendconfig

import (
	"strings"
	"testing"
)

// The remote (TFE-protocol) backend carries none of the bucket-shaped fields: its
// addressing rides raw --backend-config flags or the ambient TF_TOKEN_<hostname>
// environment variable, and its workspaces block lives in the backend declaration.
// Validation must accept the type and say exactly that, instead of demanding buckets.
func TestValidate_RemoteIsValidWithGuidanceWarning(t *testing.T) {
	result := Validate(&TofuBackendConfig{BackendType: "remote"})
	if !result.Valid {
		t.Fatalf("Validate(remote) = invalid, want valid; missing: %+v", result.MissingFields)
	}
	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "TF_TOKEN_") {
		t.Fatalf("Validate(remote) warnings = %v, want the addressing guidance", result.Warnings)
	}
}

// Unknown types fail loudly and the message enumerates what IS supported -- the next
// backend type added without a validation arm shows up here.
func TestValidate_UnknownTypeNamesTheSupportedSet(t *testing.T) {
	result := Validate(&TofuBackendConfig{BackendType: "consul"})
	if result.Valid {
		t.Fatal("Validate(consul) = valid, want invalid")
	}
	if len(result.MissingFields) != 1 || !strings.Contains(result.MissingFields[0].Description, "remote") {
		t.Fatalf("unknown-type message must enumerate supported types incl. remote, got: %+v", result.MissingFields)
	}
}

// Anchor the bucket-shaped contract: s3 without its required fields reports each one.
func TestValidate_S3RequiresBucketKeyRegion(t *testing.T) {
	result := Validate(&TofuBackendConfig{BackendType: "s3"})
	if result.Valid {
		t.Fatal("Validate(empty s3) = valid, want invalid")
	}
	if len(result.MissingFields) != 3 {
		t.Fatalf("Validate(empty s3) missing = %d fields, want 3 (bucket, key, region): %+v", len(result.MissingFields), result.MissingFields)
	}
}

func TestValidate_LocalAndEmptyHaveNoRequirements(t *testing.T) {
	for _, backendType := range []string{"local", ""} {
		if result := Validate(&TofuBackendConfig{BackendType: backendType}); !result.Valid {
			t.Errorf("Validate(%q) = invalid, want valid", backendType)
		}
	}
}
