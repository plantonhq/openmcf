package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpCryptoKeyId = "crypto_key_id"
	OpRole        = "role"
	OpMember      = "member"
	OpEtag        = "etag"
)
