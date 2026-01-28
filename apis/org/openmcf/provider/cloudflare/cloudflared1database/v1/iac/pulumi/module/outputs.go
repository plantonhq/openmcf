package module

const (
	// OpDatabaseId is the exported stack output containing the database UUID.
	OpDatabaseId = "database_id"
	// OpDatabaseName is the exported stack output containing the database name.
	OpDatabaseName = "database_name"
	// OpConnectionString is the exported stack output containing the connection string.
	// Pulumi does not expose this today, so the module currently exports an empty string.
	OpConnectionString = "connection_string"
)
