package module

var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string
	// ReleaseName is FIXED: Cilium is the node dataplane — the agent
	// DaemonSet, operator, and generated CNI configuration are cluster
	// singletons, so one dataplane per cluster is an upstream constraint
	// and the release name never derives from metadata.name.
	ReleaseName string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart-name drift installs different software per
	// engine).
	HelmChartName: "cilium",
	HelmChartRepo: "https://helm.cilium.io",
	// Fallback when spec.chart_version is unset AND the platform's
	// defaulting middleware did not run. Keep aligned with the spec
	// default and the chart-repo index (Cilium chart and app versions
	// move together: chart 1.19.6 = Cilium 1.19.6).
	DefaultChartVersion: "1.19.6",
	ReleaseName:         "cilium",
}
