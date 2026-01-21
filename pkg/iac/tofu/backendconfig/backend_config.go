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
	// BackendBucket specifies the bucket or container name for remote backends
	BackendBucket string
	// BackendKey specifies the state file path within the bucket
	BackendKey string
	// BackendRegion specifies the region for S3 backends
	BackendRegion string
	// BackendEndpoint specifies a custom S3-compatible endpoint (for R2, MinIO, etc.)
	BackendEndpoint string
	// S3Compatible indicates this is an S3-compatible backend requiring skip flags
	S3Compatible bool
}

// IsS3Compatible returns true if this is an S3-compatible backend (R2, MinIO, etc.)
// Detection signals: explicit endpoint is set OR region is "auto"
func (c *TofuBackendConfig) IsS3Compatible() bool {
	return c.BackendEndpoint != "" || c.BackendRegion == "auto"
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
	bucketLabelKey := tofulabels.BackendBucketLabelKey(provisionerType)
	keyLabelKey := tofulabels.BackendKeyLabelKey(provisionerType)
	regionLabelKey := tofulabels.BackendRegionLabelKey(provisionerType)
	endpointLabelKey := tofulabels.BackendEndpointLabelKey(provisionerType)

	backendType, hasType := labels[typeLabelKey]
	backendBucket, hasBucket := labels[bucketLabelKey]
	backendKey, hasKey := labels[keyLabelKey]
	backendRegion, _ := labels[regionLabelKey]
	backendEndpoint, _ := labels[endpointLabelKey]

	// If provisioner-specific labels not found, fall back to legacy terraform.* labels
	// This ensures backward compatibility for existing manifests
	if !hasType && !hasBucket && !hasKey {
		backendType, hasType = labels[tofulabels.LegacyBackendTypeLabelKey]
		backendBucket, hasBucket = labels[tofulabels.LegacyBackendBucketLabelKey]
		// Try backend.key first, then fall back to deprecated backend.object
		backendKey, hasKey = labels[tofulabels.LegacyBackendKeyLabelKey]
		if !hasKey {
			backendKey, hasKey = labels[tofulabels.LegacyBackendObjectLabelKey]
		}
		backendRegion, _ = labels[tofulabels.LegacyBackendRegionLabelKey]
		backendEndpoint, _ = labels[tofulabels.LegacyBackendEndpointLabelKey]
		// Update label keys for error messages
		if hasType || hasBucket || hasKey {
			typeLabelKey = tofulabels.LegacyBackendTypeLabelKey
			bucketLabelKey = tofulabels.LegacyBackendBucketLabelKey
			keyLabelKey = tofulabels.LegacyBackendKeyLabelKey
		}
	}

	// All labels are optional - return nil if none are present
	if !hasType && !hasBucket && !hasKey {
		return nil, nil
	}

	// If type is present, key must also be present
	if hasType && !hasKey {
		return nil, fmt.Errorf("both %s and %s must be specified together",
			typeLabelKey, keyLabelKey)
	}

	if backendType == "" {
		return nil, fmt.Errorf("backend type label cannot be empty")
	}

	// Validate supported backend types
	switch backendType {
	case "s3", "gcs", "azurerm", "local":
		// Supported backend types
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", backendType)
	}

	// For remote backends (non-local), bucket is required
	if backendType != "local" && !hasBucket {
		return nil, fmt.Errorf("%s is required for %s backend", bucketLabelKey, backendType)
	}

	config := &TofuBackendConfig{
		BackendType:     backendType,
		BackendBucket:   backendBucket,
		BackendKey:      backendKey,
		BackendRegion:   backendRegion,
		BackendEndpoint: backendEndpoint,
	}
	// Compute S3-compatible flag based on endpoint or region=auto
	config.S3Compatible = config.IsS3Compatible()

	return config, nil
}
