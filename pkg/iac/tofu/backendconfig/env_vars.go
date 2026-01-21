package backendconfig

import "os"

// Environment variable names for backend configuration.
// These provide an alternative to CLI flags and manifest labels,
// useful for CI/CD pipelines and 12-factor app patterns.
const (
	// EnvBackendType specifies the backend type (s3, gcs, azurerm, local)
	EnvBackendType = "PROJECT_PLANTON_BACKEND_TYPE"

	// EnvBackendBucket specifies the state bucket/container name
	EnvBackendBucket = "PROJECT_PLANTON_BACKEND_BUCKET"

	// EnvBackendRegion specifies the AWS region (use "auto" for S3-compatible backends)
	EnvBackendRegion = "PROJECT_PLANTON_BACKEND_REGION"

	// EnvBackendEndpoint specifies a custom S3-compatible endpoint URL (R2, MinIO, etc.)
	EnvBackendEndpoint = "PROJECT_PLANTON_BACKEND_ENDPOINT"
)

// Note: Backend key is intentionally NOT configurable via environment variable.
// State paths should be explicit and traceable via manifest labels or CLI flags.

// EnvBackendConfig holds backend configuration read from environment variables.
type EnvBackendConfig struct {
	BackendType     string
	BackendBucket   string
	BackendRegion   string
	BackendEndpoint string
}

// ReadFromEnv reads backend configuration from environment variables.
// Returns an EnvBackendConfig with values from PROJECT_PLANTON_BACKEND_* variables.
// Empty strings are returned for unset variables.
func ReadFromEnv() EnvBackendConfig {
	return EnvBackendConfig{
		BackendType:     os.Getenv(EnvBackendType),
		BackendBucket:   os.Getenv(EnvBackendBucket),
		BackendRegion:   os.Getenv(EnvBackendRegion),
		BackendEndpoint: os.Getenv(EnvBackendEndpoint),
	}
}

// HasAnyValues returns true if any environment variable was set.
func (e EnvBackendConfig) HasAnyValues() bool {
	return e.BackendType != "" || e.BackendBucket != "" ||
		e.BackendRegion != "" || e.BackendEndpoint != ""
}
