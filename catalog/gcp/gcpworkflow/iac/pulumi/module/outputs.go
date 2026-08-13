package module

// Output keys must match the field names in outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpWorkflowId   = "workflow_id"
	OpWorkflowName = "workflow_name"
	OpRevisionId   = "revision_id"
	OpState        = "state"
)
