package module

// Module-owned constants pinned to the Percona Operator for MySQL based on
// Percona XtraDB Cluster, v1.20.0 release train. Users never configure
// these: the spec chooses versions ONLY through image_name (the MySQL
// version seam) — everything else here is the module's pin, upgraded by
// shipping a new module build, never by user input. The Terraform twin
// carries the same literals in locals.tf — keep them in lockstep.
var vars = struct {
	// CrVersion is spec.crVersion — tells the operator which API/behavior
	// contract the resource was authored against.
	CrVersion string

	// PxcDefaultImage is the database image used when spec.image_name is
	// empty (MySQL 8.4, the operator release's companion build).
	PxcDefaultImage string

	// HaproxyImage / ProxysqlImage / LogCollectorImage / BackupImage are
	// the sidecar/proxy/backup companion images of the pinned operator
	// release.
	HaproxyImage      string
	ProxysqlImage     string
	LogCollectorImage string
	BackupImage       string

	// UpgradeOptions posture: the version-service endpoint is upstream's
	// default; apply=disabled turns AUTOMATIC upgrades off — version
	// changes happen only through explicit image changes (SmartUpdate
	// still rolls them safely).
	VersionServiceEndpoint string
	UpgradeApply           string

	// PxcPodsFinalizer makes deletion drain database pods in reverse
	// order so the last node down holds the newest data — the upstream
	// default finalizer posture.
	PxcPodsFinalizer string

	// DefaultUpdateStrategy is SmartUpdate: the operator orders restarts
	// safely (reverse ordinal, primary last).
	DefaultUpdateStrategy string

	// DefaultInstances / DefaultProxyReplicas are the production-shape
	// defaults (Galera quorum; proxy availability).
	DefaultInstances     int
	DefaultProxyReplicas int

	// DefaultPitrTimeBetweenUploads is the binlog-upload cadence (seconds)
	// when the spec leaves it unset — upstream's default recovery-point
	// objective.
	DefaultPitrTimeBetweenUploads int

	// CertManagerIssuerGroup / DefaultIssuerKind complete the issuerConf
	// reference when spec.tls.issuer is set.
	CertManagerIssuerGroup string
	DefaultIssuerKind      string
}{
	CrVersion:                     "1.20.0",
	PxcDefaultImage:               "percona/percona-xtradb-cluster:8.4.8-8.1",
	HaproxyImage:                  "percona/haproxy:2.8.18-1",
	ProxysqlImage:                 "percona/proxysql2:2.7.3-1.3",
	LogCollectorImage:             "percona/fluentbit:5.0.6-1",
	BackupImage:                   "percona/percona-xtrabackup:8.4.0-5.1",
	VersionServiceEndpoint:        "https://check.percona.com",
	UpgradeApply:                  "disabled",
	PxcPodsFinalizer:              "percona.com/delete-pxc-pods-in-order",
	DefaultUpdateStrategy:         "SmartUpdate",
	DefaultInstances:              3,
	DefaultProxyReplicas:          3,
	DefaultPitrTimeBetweenUploads: 60,
	CertManagerIssuerGroup:        "cert-manager.io",
	DefaultIssuerKind:             "ClusterIssuer",
}
