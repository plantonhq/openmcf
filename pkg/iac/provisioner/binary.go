package provisioner

import (
	"fmt"
	"os/exec"
)

// HclBinary represents the HCL-based IaC binary to use (Tofu or Terraform).
// Both binaries are functionally equivalent - they execute the same HCL modules.
type HclBinary string

const (
	// HclBinaryTofu represents the OpenTofu CLI binary.
	HclBinaryTofu HclBinary = "tofu"

	// HclBinaryTerraform represents the Terraform CLI binary.
	HclBinaryTerraform HclBinary = "terraform"
)

// String returns the binary name as a string.
func (b HclBinary) String() string {
	return string(b)
}

// DisplayName returns a human-friendly display name for the binary.
func (b HclBinary) DisplayName() string {
	switch b {
	case HclBinaryTofu:
		return "OpenTofu"
	case HclBinaryTerraform:
		return "Terraform"
	default:
		return string(b)
	}
}

// CheckAvailable verifies that the binary is available in the system PATH.
// Returns nil if available, or an error with helpful installation guidance.
func (b HclBinary) CheckAvailable() error {
	_, err := exec.LookPath(string(b))
	if err != nil {
		return fmt.Errorf("%s CLI not found in PATH. Please install %s to continue", b.DisplayName(), b.DisplayName())
	}
	return nil
}

// HclBinaryFromProvisionerType converts a ProvisionerType to its corresponding HclBinary.
// Returns empty string for non-HCL provisioners (e.g., Pulumi).
func HclBinaryFromProvisionerType(pt ProvisionerType) HclBinary {
	switch pt {
	case ProvisionerTypeTofu:
		return HclBinaryTofu
	case ProvisionerTypeTerraform:
		return HclBinaryTerraform
	default:
		return ""
	}
}
