package module

var vars = struct {
	// DefaultImageRepo is the ClickHouse server image used when
	// spec.image.repo is empty.
	DefaultImageRepo string

	// KeeperImageRepo is the image for the managed Keeper ensemble. The
	// spec exposes no Keeper image override; the tag is pinned to
	// spec.version (Keeper images are published in lockstep with server
	// releases) instead of the operator's own fallback — `latest` —
	// which would float across pod restarts.
	KeeperImageRepo string

	// DefaultClusterName is the logical cluster name when
	// spec.cluster_name is empty (the spec's documented default).
	DefaultClusterName string

	// KeeperClusterName is the CHK's internal cluster name — a fixed
	// module constant, never user-visible in connection strings.
	KeeperClusterName string

	// DefaultKeeperReplicas / DefaultKeeperDiskSize back the managed
	// Keeper when coordination.keeper leaves them unset.
	DefaultKeeperReplicas int
	DefaultKeeperDiskSize string

	// HttpPort / TcpPort are ClickHouse's fixed default interface ports —
	// the operator's generated Services listen on them and the spec does
	// not remap them.
	HttpPort int
	TcpPort  int

	// AuthSecretSuffix builds the module-managed Secret name
	// `<metadata.name><suffix>` that carries one key per provisioned user.
	AuthSecretSuffix string
}{
	DefaultImageRepo:      "clickhouse/clickhouse-server",
	KeeperImageRepo:       "clickhouse/clickhouse-keeper",
	DefaultClusterName:    "main",
	KeeperClusterName:     "keeper",
	DefaultKeeperReplicas: 3,
	DefaultKeeperDiskSize: "10Gi",
	HttpPort:              8123,
	TcpPort:               9000,
	AuthSecretSuffix:      "-clickhouse-auth",
}
