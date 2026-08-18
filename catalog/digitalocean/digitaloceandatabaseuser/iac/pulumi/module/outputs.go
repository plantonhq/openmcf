package module

const (
	// OpClusterId is the UUID of the cluster the user belongs to.
	OpClusterId = "cluster_id"
	// OpUserName is the user's name (its API identity within the cluster).
	OpUserName = "user_name"
	// OpRole is the role DigitalOcean assigned (normally "normal").
	OpRole = "role"
	// OpPassword is the server-generated password (secret).
	OpPassword = "password"
	// OpAccessCert is the Kafka mTLS access certificate (secret; Kafka only).
	OpAccessCert = "access_cert"
	// OpAccessKey is the Kafka mTLS access key (secret; Kafka only).
	OpAccessKey = "access_key"
)
