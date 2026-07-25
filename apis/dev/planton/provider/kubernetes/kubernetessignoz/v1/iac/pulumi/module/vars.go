package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name
// drift deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. The
	// signoz chart's version tracks the SigNoz application version in
	// lockstep (chart 0.133.0 = app v0.133.0).
	DefaultChartVersion string
	// Fallbacks mirroring the proto defaults.
	DefaultClickHouseDiskSize string
	DefaultServerDiskSize     string
	// The bundled arm's ClickHouse username (the chart's own value; the
	// password is module-generated — the chart's publicly-documented
	// default password never ships).
	BundledClickHouseUser string
	// The bundled installation's client port (ClickHouse native
	// protocol) — the exported clickhouse_endpoint's port.
	ClickHouseTcpPort int
	// Server + collector service ports (the chart's defaults; the
	// exported endpoints are built from them).
	ServerHttpPort int
	OtlpGrpcPort   int
	OtlpHttpPort   int
	// The chart wraps the ClickHouseInstallation's operator-generated
	// StatefulSet names in ~27 characters of scaffolding
	// (chi-<name>-clickhouse-cluster-0-0) within Kubernetes'
	// 63-character cap — a resource name longer than this corrupts the
	// naming contract the outputs promise.
	MaxNameLength int
}{
	HelmChartName:             "signoz",
	HelmChartRepo:             "https://charts.signoz.io",
	DefaultChartVersion:       "0.133.0",
	DefaultClickHouseDiskSize: "20Gi",
	DefaultServerDiskSize:     "1Gi",
	BundledClickHouseUser:     "admin",
	ClickHouseTcpPort:         9000,
	ServerHttpPort:            8080,
	OtlpGrpcPort:              4317,
	OtlpHttpPort:              4318,
	MaxNameLength:             30,
}
