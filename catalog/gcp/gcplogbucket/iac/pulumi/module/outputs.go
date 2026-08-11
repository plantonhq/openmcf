package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpBucketName      = "bucket_name"
	OpLinkedDatasetId = "linked_dataset_id"
)
