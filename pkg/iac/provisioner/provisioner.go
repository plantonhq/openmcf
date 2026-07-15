package provisioner

import (
	"fmt"
	"strings"

	"github.com/plantonhq/planton/pkg/iac/provisionerannotationkeys"
	"github.com/plantonhq/planton/pkg/reflection/metadatareflect"
	"google.golang.org/protobuf/proto"
)

// ProvisionerType represents the IaC provisioner type
type ProvisionerType int

const (
	ProvisionerTypeUnspecified ProvisionerType = iota
	ProvisionerTypePulumi
	ProvisionerTypeTofu
	ProvisionerTypeTerraform
)

// String returns the string representation of the provisioner type
func (p ProvisionerType) String() string {
	switch p {
	case ProvisionerTypePulumi:
		return "pulumi"
	case ProvisionerTypeTofu:
		return "tofu"
	case ProvisionerTypeTerraform:
		return "terraform"
	default:
		return "unspecified"
	}
}

// ExtractFromManifest extracts the provisioner type from manifest annotations
// Returns:
//   - ProvisionerType and nil error if the annotation exists and is valid
//   - ProvisionerTypeUnspecified and nil error if the annotation is missing (needs user prompt)
//   - ProvisionerTypeUnspecified and error if the annotation value is invalid
func ExtractFromManifest(manifest proto.Message) (ProvisionerType, error) {
	annotations := metadatareflect.ExtractAnnotations(manifest)
	if annotations == nil {
		return ProvisionerTypeUnspecified, nil
	}

	provisioner, ok := annotations[provisionerannotationkeys.ProvisionerAnnotationKey]
	if !ok || provisioner == "" {
		// Annotation not present - return unspecified (caller should prompt user)
		return ProvisionerTypeUnspecified, nil
	}

	// Case-insensitive matching
	provisionerLower := strings.ToLower(strings.TrimSpace(provisioner))

	switch provisionerLower {
	case "pulumi":
		return ProvisionerTypePulumi, nil
	case "tofu":
		return ProvisionerTypeTofu, nil
	case "terraform":
		return ProvisionerTypeTerraform, nil
	default:
		return ProvisionerTypeUnspecified, fmt.Errorf("invalid provisioner value '%s': must be one of 'pulumi', 'tofu', or 'terraform'", provisioner)
	}
}

// FromString converts a string to ProvisionerType (case-insensitive)
func FromString(s string) (ProvisionerType, error) {
	sLower := strings.ToLower(strings.TrimSpace(s))
	switch sLower {
	case "pulumi":
		return ProvisionerTypePulumi, nil
	case "tofu":
		return ProvisionerTypeTofu, nil
	case "terraform":
		return ProvisionerTypeTerraform, nil
	default:
		return ProvisionerTypeUnspecified, fmt.Errorf("invalid provisioner '%s': must be one of 'pulumi', 'tofu', or 'terraform'", s)
	}
}
