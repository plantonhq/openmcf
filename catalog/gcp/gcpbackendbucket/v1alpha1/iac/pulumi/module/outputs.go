package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpSelfLink          = "self_link"
	OpBackendBucketName = "backend_bucket_name"
	OpBucketName        = "bucket_name"
)
