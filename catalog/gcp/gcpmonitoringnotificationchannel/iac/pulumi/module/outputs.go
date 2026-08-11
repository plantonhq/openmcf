package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpChannelName        = "channel_name"
	OpVerificationStatus = "verification_status"
)
