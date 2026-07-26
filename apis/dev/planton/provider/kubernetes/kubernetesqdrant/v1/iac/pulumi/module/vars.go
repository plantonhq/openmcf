package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. The
	// chart version tracks the Qdrant release it ships (1.18.2 = app
	// v1.18.2); the chart pin governs.
	DefaultChartVersion string
	// Fallback for spec.storage.size / spec.snapshots.size — mirrors of
	// the proto defaults.
	DefaultStorageSize   string
	DefaultSnapshotsSize string
	// Where the client-listener TLS Secret mounts and the config paths
	// the engine reads the certificate from (values.go renders both
	// sides from these constants so they can never drift apart).
	TlsMountPath string
}{
	HelmChartName:        "qdrant",
	HelmChartRepo:        "https://qdrant.github.io/qdrant-helm",
	DefaultChartVersion:  "1.18.2",
	DefaultStorageSize:   "10Gi",
	DefaultSnapshotsSize: "10Gi",
	TlsMountPath:         "/qdrant/tls",
}
