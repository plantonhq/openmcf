package module

var vars = struct {
	// HelmChartRepo is the opensearch-project chart repository serving the
	// opensearch-operator chart.
	HelmChartRepo string
	// HelmChartName is the OpenSearch operator chart
	// ("opensearch-operator"). Also the input to the chart's fullname
	// helper — the deployment-name derivation in locals.go depends on it.
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. 2.8.0 is the newest SERVED chart whose
	// default manager image is a STABLE operator release — the newer
	// served charts (2.8.3+, 3.0.x) default to a PRERELEASE image and the
	// 3.x line migrates the CRDs to the opensearch.org API group.
	DefaultChartVersion string

	// CrdDirectory is where the module-owned CRD files are staged,
	// relative to the Pulumi project directory (the working directory at
	// program run time). The staged files are the chart 2.8.0 CRDs —
	// upgrade them together with DefaultChartVersion.
	CrdDirectory string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://opensearch-project.github.io/opensearch-k8s-operator/",
	HelmChartName:       "opensearch-operator",
	DefaultChartVersion: "2.8.0",
	CrdDirectory:        "../crds",
	HelmTimeoutSeconds:  600,
}
