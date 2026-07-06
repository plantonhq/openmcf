package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpDiscoveryAddress = "discovery_address"
	OpDiscoveryPort    = "discovery_port"
	OpInstanceUid      = "instance_uid"
	OpNodeSizeGb       = "node_size_gb"
	OpName             = "name"
	OpBackupCollection = "backup_collection"
)
