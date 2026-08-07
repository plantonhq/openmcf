package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
// The private key is write-only in GCP and deliberately never exported.
const (
	OpSelfLink        = "self_link"
	OpCertificateName = "certificate_name"
	OpCertificateId   = "certificate_id"
	OpExpireTime      = "expire_time"
	OpRegion          = "region"
)
