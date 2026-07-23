package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart and
	// app versions move separately (chart 0.11.0 ships Valkey 9.1.1); the
	// chart pin governs.
	DefaultChartVersion string
}{
	HelmChartName:       "valkey",
	HelmChartRepo:       "https://valkey.io/valkey-helm/",
	DefaultChartVersion: "0.11.0",
}
