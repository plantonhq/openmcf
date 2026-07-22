package module

var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart-name drift installs different software per
	// engine; the repository serves the chart as "ingress-nginx").
	HelmChartName: "ingress-nginx",
	HelmChartRepo: "https://kubernetes.github.io/ingress-nginx",
	// Fallback when spec.chart_version is unset AND the platform's
	// defaulting middleware did not run. Keep aligned with the spec
	// default and the chart-repo index (chart 4.15.1 = controller v1.15.1).
	DefaultChartVersion: "4.15.1",
}
