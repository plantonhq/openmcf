package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce.
	DefaultChartVersion string
}{
	HelmChartName:       "cert-manager",
	HelmChartRepo:       "https://charts.jetstack.io",
	DefaultChartVersion: "v1.20.3",
}
