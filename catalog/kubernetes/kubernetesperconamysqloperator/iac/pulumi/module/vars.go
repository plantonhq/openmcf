package module

var vars = struct {
	// HelmChartRepo is Percona's chart repository — it publishes the
	// operator charts for every Percona database product side by side.
	HelmChartRepo string
	// HelmChartName is the Percona Operator for MySQL (Percona XtraDB
	// Cluster) chart ("pxc-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart and operator versions move TOGETHER for
	// this chart (chart 1.20.0 ships operator 1.20.0); the chart pin
	// governs.
	DefaultChartVersion string

	// DefaultImageRepository mirrors the chart's own
	// operatorImageRepository default. Needed for the tag-without-
	// repository image override: the chart's repository value is always
	// suffixed with the chart's app version, so pinning a DIFFERENT tag
	// requires the chart's full-image override ("<repository>:<tag>"),
	// and the repository half must come from somewhere when the spec
	// leaves it empty.
	DefaultImageRepository string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:          "https://percona.github.io/percona-helm-charts",
	HelmChartName:          "pxc-operator",
	DefaultChartVersion:    "1.20.0",
	DefaultImageRepository: "percona/percona-xtradb-cluster-operator",
	HelmTimeoutSeconds:     600,
}
