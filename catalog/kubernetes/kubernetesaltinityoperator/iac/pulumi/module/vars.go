package module

var vars = struct {
	// HelmChartRepo is the Altinity-served repository for the
	// altinity-clickhouse-operator chart.
	HelmChartRepo string
	// HelmChartName is the Altinity ClickHouse operator chart
	// ("altinity-clickhouse-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart versions track operator releases
	// one-to-one (chart 0.27.2 = operator image 0.27.2).
	DefaultChartVersion string

	// DefaultCredentialsUsername mirrors the spec default for
	// operator_credentials.username — resolved here so the rendered
	// secret block is complete whether or not the platform's defaulting
	// middleware ran.
	DefaultCredentialsUsername string

	// MetricsPort is the metrics-exporter sidecar's per-cluster metrics
	// port on the chart's "<fullname>-metrics" Service (the ch-metrics
	// endpoint; the operator's own op-metrics rides 9999).
	MetricsPort int

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters plus the pre-install CRD
	// hook job; atomic rolls back on expiry so a wedged install never
	// lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:              "https://docs.altinity.com/clickhouse-operator/",
	HelmChartName:              "altinity-clickhouse-operator",
	DefaultChartVersion:        "0.27.2",
	DefaultCredentialsUsername: "clickhouse_operator",
	MetricsPort:                8888,
	HelmTimeoutSeconds:         600,
}
