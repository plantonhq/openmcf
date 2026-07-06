package module

// Output keys must match the field names in stack_outputs.proto — the outputs
// transformer maps raw engine outputs onto the proto by name.
const (
	OpDatabaseId             = "database_id"
	OpDatabaseName           = "database_name"
	OpUid                    = "uid"
	OpCreateTime             = "create_time"
	OpEarliestVersionTime    = "earliest_version_time"
	OpVersionRetentionPeriod = "version_retention_period"
	OpKeyPrefix              = "key_prefix"
	OpUpdateTime             = "update_time"
)
