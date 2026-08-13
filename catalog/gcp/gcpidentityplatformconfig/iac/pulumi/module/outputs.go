package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpConfigName        = "config_name"
	OpApiKey            = "api_key"
	OpFirebaseSubdomain = "firebase_subdomain"
)
