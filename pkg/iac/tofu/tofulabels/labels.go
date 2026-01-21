package tofulabels

import "fmt"

// BackendTypeLabelKey returns the backend type label key for the given provisioner.
// The provisioner should be "terraform" or "tofu".
// Example: BackendTypeLabelKey("terraform") returns "terraform.project-planton.org/backend.type"
func BackendTypeLabelKey(provisioner string) string {
	return fmt.Sprintf("%s.project-planton.org/backend.type", provisioner)
}

// BackendObjectLabelKey returns the backend object label key for the given provisioner.
// The provisioner should be "terraform" or "tofu".
// Example: BackendObjectLabelKey("tofu") returns "tofu.project-planton.org/backend.object"
func BackendObjectLabelKey(provisioner string) string {
	return fmt.Sprintf("%s.project-planton.org/backend.object", provisioner)
}

// Legacy constants for backward compatibility.
// These are kept to ensure existing manifests using terraform.* labels
// continue to work regardless of the provisioner being used.
const (
	// LegacyBackendTypeLabelKey is the legacy backend type label (terraform prefix)
	LegacyBackendTypeLabelKey = "terraform.project-planton.org/backend.type"

	// LegacyBackendObjectLabelKey is the legacy backend object label (terraform prefix)
	// For S3: "bucket-name/path/to/state"
	// For GCS: "bucket-name/path/to/state"
	// For Azure: "container-name/path/to/state"
	LegacyBackendObjectLabelKey = "terraform.project-planton.org/backend.object"
)
