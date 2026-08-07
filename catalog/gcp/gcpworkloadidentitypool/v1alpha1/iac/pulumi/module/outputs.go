package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpName                   = "name"
	OpWorkloadIdentityPoolId = "workload_identity_pool_id"
	OpState                  = "state"
)
