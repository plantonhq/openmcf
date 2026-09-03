package module

var vars = struct {
	// HelmChartRepo is the open-telemetry chart repository serving the
	// opentelemetry-operator chart.
	HelmChartRepo string
	// HelmChartName is the OpenTelemetry operator chart
	// ("opentelemetry-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. 0.120.0 is the newest SERVED stable chart
	// (= operator appVersion 0.156.0, verified against the repository
	// index). The CRDs are derived from whatever version is pinned, so a
	// bump here changes exactly one line.
	DefaultChartVersion string

	// CrdRenderOverride is the one values document merged LAST when the
	// pinned chart is rendered to derive its CRDs: the chart's CRD switch
	// turned on. The release itself installs with the switch off and
	// skip_crds set, so this never reaches the cluster through Helm. Twin
	// of the Terraform module's helm_crds_render_override.
	CrdRenderOverride string

	// CrdRenderApiVersions are the API versions the chart gates its CRD
	// templates on (.Capabilities.APIVersions): the collector CRD's
	// cert-manager.io/inject-ca-from annotation renders only when
	// cert-manager.io/v1 is declared served. Twin of the Terraform
	// module's helm_crds_api_versions.
	CrdRenderApiVersions []string

	// ManagerImagePath is the manager image's repository path WITHOUT the
	// registry — the half image_registry replaces (the air-gap mirror
	// seam). Keep in lockstep with the Terraform module's manager_image.
	ManagerImagePath string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls plus the cert-manager webhook-certificate
	// issuance the manager pod mounts; atomic rolls back on expiry so a
	// wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:        "https://open-telemetry.github.io/opentelemetry-helm-charts",
	HelmChartName:        "opentelemetry-operator",
	DefaultChartVersion:  "0.120.0",
	CrdRenderOverride:    "crds:\n  create: true\n",
	CrdRenderApiVersions: []string{"cert-manager.io/v1"},
	ManagerImagePath:     "open-telemetry/opentelemetry-operator/opentelemetry-operator",
	HelmTimeoutSeconds:   600,
}
