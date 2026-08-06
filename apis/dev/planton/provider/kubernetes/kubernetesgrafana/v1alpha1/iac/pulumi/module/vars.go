package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
//
// KNOW THIS about the repo URL: the grafana chart's canonical home is the
// grafana-community index — the old https://grafana.github.io/helm-charts
// stopped serving new versions at chart 10.5.x, and kube-prometheus-stack's
// own dependency block points at the community repo. Never "fix" this back
// to the grafana.github.io URL.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 12.8.0 ships Grafana 13.1.1; the chart pin governs.
	DefaultChartVersion string
	// Fallback for spec.storage.size — mirror of the proto default.
	DefaultStorageSize string
	// Grafana's admin-credential Secret keys — the chart's own key names,
	// used both for the chart-generated Secret and as the defaults for an
	// existing Secret's key overrides.
	AdminUserKey     string
	AdminPasswordKey string
	// Service port the chart exposes (targets container 3000).
	ServicePort int
}{
	HelmChartName:       "grafana",
	HelmChartRepo:       "https://grafana-community.github.io/helm-charts",
	DefaultChartVersion: "12.8.0",
	DefaultStorageSize:  "10Gi",
	AdminUserKey:        "admin-user",
	AdminPasswordKey:    "admin-password",
	ServicePort:         80,
}
