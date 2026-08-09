package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSecretName        = "secret_name"
	OpSecretId          = "secret_id"
	OpLatestVersionName = "latest_version_name"
)
