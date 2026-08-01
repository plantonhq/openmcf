package module

var vars = struct {
	// HelmChartRepo is the ray-project chart repository serving the
	// kuberay-operator chart.
	HelmChartRepo string
	// HelmChartName is the KubeRay operator chart ("kuberay-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart 1.6.2 pairs with operator image
	// quay.io/kuberay/operator:v1.6.2. The chart's crds/-directory ray.io
	// CRDs are install-once — a bump here never upgrades them (apply the
	// new release's CRD files manually when a bump changes them).
	DefaultChartVersion string

	// OperatorImagePath is the operator image's repository path WITHOUT
	// the registry — the half image_registry replaces (the air-gap mirror
	// seam; the chart default registry is quay.io). Keep in lockstep with
	// the Terraform module's operator_image.
	OperatorImagePath string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers the image pull plus the server-side install of the
	// three ~1MB ray.io CRDs; atomic rolls back on expiry so a wedged
	// install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://ray-project.github.io/kuberay-helm",
	HelmChartName:       "kuberay-operator",
	DefaultChartVersion: "1.6.2",
	OperatorImagePath:   "kuberay/operator",
	HelmTimeoutSeconds:  600,
}
