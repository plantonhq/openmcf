package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpTriggerName                   = "trigger_name"
	OpPartnerChannelActivationToken = "partner_channel_activation_token"
	OpTriggerId                     = "trigger_id"
)
