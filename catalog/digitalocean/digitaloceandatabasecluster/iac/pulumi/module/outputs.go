package module

const (
	// OpClusterId is the UUID of the created database cluster.
	OpClusterId = "cluster_id"
	// OpConnectionUri is the full public connection URI (includes user/pass/db).
	OpConnectionUri = "connection_uri"
	// OpHost is the public hostname of the cluster.
	OpHost = "host"
	// OpPort is the port the cluster listens on.
	OpPort = "port"
	// OpDatabaseUser is the default user.
	OpDatabaseUser = "database_user"
	// OpDatabasePassword is that user's password.
	OpDatabasePassword = "database_password"
	// OpPrivateHost is the private-network hostname (same-VPC access).
	OpPrivateHost = "private_host"
	// OpPrivateUri is the private-network connection URI.
	OpPrivateUri = "private_uri"
	// OpDatabaseName is the name of the default database.
	OpDatabaseName = "database_name"
	// OpUiHost is the OpenSearch Dashboards hostname (OpenSearch only).
	OpUiHost = "ui_host"
	// OpUiPort is the OpenSearch Dashboards port (OpenSearch only).
	OpUiPort = "ui_port"
	// OpUiUri is the OpenSearch Dashboards connection URI (OpenSearch only).
	OpUiUri = "ui_uri"
	// OpUiDatabase is the OpenSearch Dashboards default database (OpenSearch only).
	OpUiDatabase = "ui_database"
	// OpUiUser is the OpenSearch Dashboards username (OpenSearch only).
	OpUiUser = "ui_user"
	// OpUiPassword is the OpenSearch Dashboards password (OpenSearch only).
	OpUiPassword = "ui_password"
)
