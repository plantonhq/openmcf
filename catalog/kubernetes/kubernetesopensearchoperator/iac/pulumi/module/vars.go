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
	// 3.x line migrates the CRDs to the opensearch.org API group. The CRDs
	// are derived from whatever version is pinned, so a bump here changes
	// exactly one line.
	DefaultChartVersion string

	// CrdRenderOverride is the one values document merged LAST when the
	// pinned chart is rendered to derive its CRDs: the chart's CRD switch
	// turned on. The release itself installs with the switch off and
	// skip_crds set, so this never reaches the cluster through Helm. Twin
	// of the Terraform module's helm_crds_args.render_override.
	CrdRenderOverride string

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
	CrdRenderOverride:   "installCRDs: true\n",
	HelmTimeoutSeconds:  600,
}
