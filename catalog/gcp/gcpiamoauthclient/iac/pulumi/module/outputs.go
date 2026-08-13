package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpClientId     = "client_id"
	OpClientName   = "client_name"
	OpState        = "state"
	OpClientSecret = "client_secret"
)
