package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpTableId      = "table_id"
	OpTableName    = "table_name"
	OpInstanceName = "instance_name"
)
