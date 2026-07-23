package module

var vars = struct {
	// HelmChartRepo serves BOTH charts (cloudnative-pg and
	// plugin-barman-cloud) — the upstream project publishes them from one
	// repository.
	HelmChartRepo string
	// HelmChartName is the operator chart ("cloudnative-pg").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default and the chart-repo index (chart 0.29.0 ships
	// operator 1.30.0 — chart and app versions move SEPARATELY; the chart
	// pin governs).
	DefaultChartVersion string
	// ReleaseName is FIXED: the operator registers cluster-scoped CRDs and
	// mutating/validating webhooks whose service name is baked into the
	// chart ("cnpg-webhook-service" — it is embedded in the webhook
	// certificate and cannot be configured), so a second installation would
	// fight over both. One operator per cluster is an upstream constraint —
	// the release name never derives from metadata.name.
	ReleaseName string

	// PluginChartName is the Barman Cloud plugin chart
	// ("plugin-barman-cloud") — the object-store backup path for every
	// KubernetesPostgres the operator manages. Installed as a SEPARATE
	// release in the SAME namespace: upstream forbids folding the plugin
	// into the operator's release (Helm ownership of shared resources
	// would conflict).
	PluginChartName string
	// PluginReleaseName is FIXED for the same singleton reason as the
	// operator: the plugin's gRPC service name ("barman-cloud") is baked
	// into its TLS certificate and cannot be configured.
	PluginReleaseName string
	// DefaultPluginChartVersion is the fallback when
	// spec.barman_cloud_plugin.chart_version is unset (chart 0.7.0 ships
	// plugin v0.13.0).
	DefaultPluginChartVersion string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of both
	// releases. 600s covers image pulls on cold clusters; atomic rolls
	// back on expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:             "https://cloudnative-pg.github.io/charts",
	HelmChartName:             "cloudnative-pg",
	DefaultChartVersion:       "0.29.0",
	ReleaseName:               "cnpg",
	PluginChartName:           "plugin-barman-cloud",
	PluginReleaseName:         "plugin-barman-cloud",
	DefaultPluginChartVersion: "0.7.0",
	HelmTimeoutSeconds:        600,
}
