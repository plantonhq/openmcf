package module

const (
	// OpReplicaId is the replica's own UUID.
	OpReplicaId = "replica_id"
	// OpClusterId is the UUID of the primary cluster the replica follows.
	OpClusterId = "cluster_id"
	// OpReplicaName is the replica's name (its API identity within the cluster).
	OpReplicaName = "replica_name"
	// OpHost is the public hostname of the replica endpoint.
	OpHost = "host"
	// OpPrivateHost is the private-network hostname (same-VPC access).
	OpPrivateHost = "private_host"
	// OpPort is the port the replica listens on.
	OpPort = "port"
	// OpDatabase is the default database served by the replica.
	OpDatabase = "database"
	// OpUser is the replica's default username.
	OpUser = "user"
	// OpPassword is that user's password (secret).
	OpPassword = "password"
	// OpUri is the full public connection URI (secret; includes credentials).
	OpUri = "uri"
	// OpPrivateUri is the private-network connection URI (secret).
	OpPrivateUri = "private_uri"
)
