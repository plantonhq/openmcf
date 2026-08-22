package module

const (
	// OpClusterId is the UUID of the PostgreSQL cluster the pool runs on.
	OpClusterId = "cluster_id"
	// OpPoolName is the pool's name (clients connect to it as a database name).
	OpPoolName = "pool_name"
	// OpHost is the public hostname of the pool endpoint.
	OpHost = "host"
	// OpPrivateHost is the private-network hostname (same-VPC access).
	OpPrivateHost = "private_host"
	// OpPort is the port the pool listens on.
	OpPort = "port"
	// OpUri is the full public connection URI (secret; includes credentials).
	OpUri = "uri"
	// OpPrivateUri is the private-network connection URI (secret).
	OpPrivateUri = "private_uri"
	// OpPassword is the pool user's password (secret; empty for inbound-user pools).
	OpPassword = "password"
)
