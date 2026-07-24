package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the Terraform
	// module's coalesce. Verified against the SERVED repository index at
	// https://ui.charts.kafbat.io/index.yaml: kafka-ui 1.6.4 (appVersion
	// v1.5.0). Chart and app versions move separately; the chart pin
	// governs.
	DefaultChartVersion string
}{
	HelmChartName:       "kafka-ui",
	HelmChartRepo:       "https://ui.charts.kafbat.io",
	DefaultChartVersion: "1.6.4",
}
