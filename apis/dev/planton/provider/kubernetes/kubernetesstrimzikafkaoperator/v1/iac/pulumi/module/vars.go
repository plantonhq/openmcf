package module

var vars = struct {
	// HelmChartRepo is the Strimzi project's chart repository. The same
	// charts are also published OCI at quay.io/strimzi-helm; the HTTPS
	// repository is the canonical reference the modules pin.
	HelmChartRepo string
	// HelmChartName is the Strimzi cluster operator chart
	// ("strimzi-kafka-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart and operator versions move TOGETHER
	// for this chart (chart 1.1.0 ships operator 1.1.0); the SERVED index
	// (https://strimzi.io/charts/) governs — the Chart.yaml inside the
	// Strimzi source tree carries a build-time placeholder and never
	// reflects the served version.
	DefaultChartVersion string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://strimzi.io/charts/",
	HelmChartName:       "strimzi-kafka-operator",
	DefaultChartVersion: "1.1.0",
	HelmTimeoutSeconds:  600,
}
