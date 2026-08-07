package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSelfLink        = "self_link"
	OpSslPolicyName   = "ssl_policy_name"
	OpEnabledFeatures = "enabled_features"
	OpRegion          = "region"
)
