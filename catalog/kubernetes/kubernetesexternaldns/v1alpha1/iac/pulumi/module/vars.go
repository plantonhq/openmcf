package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// releases are cut separately from the controller (chart 1.21.x ships
	// controller v0.21.x).
	DefaultChartVersion string
}{
	HelmChartName:       "external-dns",
	HelmChartRepo:       "https://kubernetes-sigs.github.io/external-dns/",
	DefaultChartVersion: "1.21.1",
}
