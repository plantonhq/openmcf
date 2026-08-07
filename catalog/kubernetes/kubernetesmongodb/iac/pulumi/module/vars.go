package module

var vars = struct {
	// Module-owned constants, pinned to the operator release the module
	// targets (KubernetesPerconaMongoOperator at v1.22.0).
	CRVersion              string
	DefaultImage           string
	BackupImage            string
	LogcollectorImage      string
	VersionServiceEndpoint string

	// The operator deletes member pods in order on CR deletion — the safe
	// teardown for a replica set (primary last).
	Finalizer string

	// AdminPasswordKey is the key in the operator-managed system-users
	// Secret holding the database-admin password (paired username key:
	// MONGODB_DATABASE_ADMIN_USER).
	AdminPasswordKey string

	MongoDBPort int
}{
	CRVersion:              "1.22.0",
	DefaultImage:           "percona/percona-server-mongodb:8.0.19-7",
	BackupImage:            "percona/percona-backup-mongodb:2.12.0",
	LogcollectorImage:      "percona/fluentbit:4.0.1-2",
	VersionServiceEndpoint: "https://check.percona.com",
	Finalizer:              "percona.com/delete-psmdb-pods-in-order",
	AdminPasswordKey:       "MONGODB_DATABASE_ADMIN_PASSWORD",
	MongoDBPort:            27017,
}
