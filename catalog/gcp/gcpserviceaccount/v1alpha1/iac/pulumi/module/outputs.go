package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpEmail     = "email"
	OpMember    = "member"
	OpUniqueId  = "unique_id"
	OpName      = "name"
	OpKeyBase64 = "key_base64"
)
