package backendconfig

import (
	"github.com/plantonhq/planton/pkg/iac/tofu/tofuannotationkeys"
	"github.com/plantonhq/planton/pkg/reflection/metadatareflect"
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

// ExtractFromManifest extracts Terraform/Tofu backend configuration from manifest annotations.
// The provisionerType should be "terraform" or "tofu" to determine which annotation prefix to use
// (e.g., tofu.planton.dev/backend.type vs terraform.planton.dev/backend.type).
func ExtractFromManifest(manifest proto.Message, provisionerType string) (*TofuBackendConfig, error) {
	annotations := metadatareflect.ExtractAnnotations(manifest)
	if annotations == nil {
		return nil, nil
	}

	backendType, hasType := annotations[tofuannotationkeys.BackendTypeAnnotationKey(provisionerType)]
	backendBucket, hasBucket := annotations[tofuannotationkeys.BackendBucketAnnotationKey(provisionerType)]
	backendKey, hasKey := annotations[tofuannotationkeys.BackendKeyAnnotationKey(provisionerType)]
	backendRegion := annotations[tofuannotationkeys.BackendRegionAnnotationKey(provisionerType)]
	backendEndpoint := annotations[tofuannotationkeys.BackendEndpointAnnotationKey(provisionerType)]

	// Return nil if no backend annotations are present
	if !hasType && !hasBucket && !hasKey {
		return nil, nil
	}

	// Extract whatever annotations are present - validation happens later via Validate()
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
