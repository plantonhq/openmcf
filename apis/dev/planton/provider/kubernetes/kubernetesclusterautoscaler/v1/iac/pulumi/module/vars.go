package module

var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string
	// ReleaseName is FIXED: the autoscaler leader-elects and owns the
	// cluster-wide scaling decision — a second installation would fight
	// the first over every scale-up, so one installation per cluster is
	// the operating model and the release name never derives from
	// metadata.name.
	ReleaseName string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart-name drift installs different software per
	// engine).
	HelmChartName: "cluster-autoscaler",
	HelmChartRepo: "https://kubernetes.github.io/autoscaler",
	// Fallback when spec.chart_version is unset AND the platform's
	// defaulting middleware did not run. Keep aligned with the spec
	// default and the chart-repo index (chart 9.59.0 ships autoscaler
	// 1.35.0 — chart and app versions move SEPARATELY; the chart pin
	// governs).
	DefaultChartVersion: "9.59.0",
	ReleaseName:         "cluster-autoscaler",
}
