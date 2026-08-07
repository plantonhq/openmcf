package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpServiceAccountId = "service_account_id"
	OpRole             = "role"
	OpMember           = "member"
	OpEtag             = "etag"
)
