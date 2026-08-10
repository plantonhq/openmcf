package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSloName     = "slo_name"
	OpServiceName = "service_name"
)
