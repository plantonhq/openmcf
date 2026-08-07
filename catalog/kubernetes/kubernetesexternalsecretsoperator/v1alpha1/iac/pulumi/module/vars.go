package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart and
	// operator versions are aligned upstream.
	DefaultChartVersion string
}{
	HelmChartName:       "external-secrets",
	HelmChartRepo:       "https://charts.external-secrets.io",
	DefaultChartVersion: "2.8.0",
}
