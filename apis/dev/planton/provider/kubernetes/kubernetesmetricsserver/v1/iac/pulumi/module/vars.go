package module

var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string
	// ReleaseName is FIXED: metrics-server registers the cluster-wide
	// v1beta1.metrics.k8s.io APIService, a singleton — one installation per
	// cluster is an upstream constraint, so the release name never derives
	// from metadata.name.
	ReleaseName string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart-name drift installs different software per
	// engine).
	HelmChartName: "metrics-server",
	HelmChartRepo: "https://kubernetes-sigs.github.io/metrics-server/",
	// Fallback when spec.chart_version is unset AND the platform's
	// defaulting middleware did not run. Keep aligned with the spec
	// default and the chart-repo index (chart 3.13.1 = metrics-server
	// 0.8.1).
	DefaultChartVersion: "3.13.1",
	ReleaseName:         "metrics-server",
}
