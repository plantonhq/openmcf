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
	// Fallback mirroring the proto default.
	DefaultServerDiskSize string
	// The ClickHouse native-protocol default port — the exported
	// clickhouse_endpoint's port when the connection declares none.
	ClickHouseTcpPort int
	// Server + collector service ports (the chart's defaults; the
	// exported endpoints are built from them).
	ServerHttpPort int
	OtlpGrpcPort   int
	OtlpHttpPort   int
	// The longest fullname-derived child is the collector Deployment
	// (`<name>-otel-collector`, a 15-character suffix) whose pod names
	// add a 16-character replica-set + pod suffix inside Kubernetes'
	// 63-character cap: 63 - 15 - 16 = 32. (The schema migrator's name
	// is a FIXED chart string, not fullname-derived.) A resource name
	// longer than this corrupts the naming contract the outputs promise.
	MaxNameLength int
}{
	HelmChartName:         "signoz",
	HelmChartRepo:         "https://charts.signoz.io",
	DefaultChartVersion:   "0.133.0",
	DefaultServerDiskSize: "1Gi",
	ClickHouseTcpPort:     9000,
	ServerHttpPort:        8080,
	OtlpGrpcPort:          4317,
	OtlpHttpPort:          4318,
	MaxNameLength:         32,
}
