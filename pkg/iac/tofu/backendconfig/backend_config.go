package backendconfig

import (
	"fmt"

	"github.com/plantonhq/project-planton/pkg/iac/tofu/tofulabels"
	"github.com/plantonhq/project-planton/pkg/reflection/metadatareflect"
	"google.golang.org/protobuf/proto"
)

// TofuBackendConfig represents the Terraform/Tofu backend configuration
type TofuBackendConfig struct {
	// BackendType specifies the backend type (e.g., "s3", "gcs", "azurerm")
	BackendType string
	// BackendObject specifies the backend object path
	// For S3: "bucket-name/path/to/state"
	// For GCS: "bucket-name/path/to/state"
	// For Azure: "container-name/path/to/state"
	BackendObject string
}

// ExtractFromManifest extracts Terraform/Tofu backend configuration from manifest labels.
// The provisionerType should be "terraform" or "tofu" to determine which label prefix to use.
// It first checks for provisioner-specific labels (e.g., tofu.project-planton.org/backend.type),
// then falls back to legacy terraform.* labels for backward compatibility.
func ExtractFromManifest(manifest proto.Message, provisionerType string) (*TofuBackendConfig, error) {
	labels := metadatareflect.ExtractLabels(manifest)
	if labels == nil {
		return nil, fmt.Errorf("no labels found in manifest")
	}

	// Try provisioner-specific labels first
	typeLabelKey := tofulabels.BackendTypeLabelKey(provisionerType)
	objectLabelKey := tofulabels.BackendObjectLabelKey(provisionerType)

	backendType, hasType := labels[typeLabelKey]
	backendObject, hasObject := labels[objectLabelKey]

	// If provisioner-specific labels not found, fall back to legacy terraform.* labels
	// This ensures backward compatibility for existing manifests
	if !hasType && !hasObject {
		backendType, hasType = labels[tofulabels.LegacyBackendTypeLabelKey]
		backendObject, hasObject = labels[tofulabels.LegacyBackendObjectLabelKey]
		// Update label keys for error messages
		if hasType || hasObject {
			typeLabelKey = tofulabels.LegacyBackendTypeLabelKey
			objectLabelKey = tofulabels.LegacyBackendObjectLabelKey
		}
	}

	// Both labels are optional - return nil if neither is present
	if !hasType && !hasObject {
		return nil, nil
	}

	// If one is present, both must be present
	if !hasType || !hasObject {
		return nil, fmt.Errorf("both %s and %s must be specified together",
			typeLabelKey, objectLabelKey)
	}

	if backendType == "" || backendObject == "" {
		return nil, fmt.Errorf("backend labels cannot be empty")
	}

	// Validate supported backend types
	switch backendType {
	case "s3", "gcs", "azurerm", "local":
		// Supported backend types
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", backendType)
	}

	return &TofuBackendConfig{
		BackendType:   backendType,
		BackendObject: backendObject,
	}, nil
}
