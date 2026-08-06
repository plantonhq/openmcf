package module

var vars = struct {
	// HelmChartRepo is Percona's chart repository — it publishes the
	// operator charts for every Percona database product side by side.
	HelmChartRepo string
	// HelmChartName is the Percona Operator for MongoDB chart
	// ("psmdb-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart and operator versions move TOGETHER for
	// this chart (chart 1.22.0 ships operator 1.22.0); the chart pin
	// governs.
	DefaultChartVersion string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://percona.github.io/percona-helm-charts",
	HelmChartName:       "psmdb-operator",
	DefaultChartVersion: "1.22.0",
	HelmTimeoutSeconds:  600,
}
