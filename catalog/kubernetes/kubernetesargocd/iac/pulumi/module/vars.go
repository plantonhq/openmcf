package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 10.2.1 ships Argo CD v3.4.5; the chart pin governs.
	DefaultChartVersion string
	// Fixed by the APPLICATION (util/settings in the argo-cd source): the
	// Secret Argo CD generates its initial admin password into at first
	// start — the name never varies with the release name, which is why
	// one generated-password instance exists per namespace.
	InitialAdminSecretName string
	// The chart's fullname budget: every child name is
	// `<fullname>-<component>` truncated at 63 characters, and the longest
	// component suffix is "-applicationset-controller" (26 chars) — names
	// past 63-26=37 characters truncate SILENTLY and break the naming
	// contract the exported outputs are built on.
	FullnameBudget int
	// The external-Redis credential contract the chart documents on
	// externalRedis.existingSecret: key `redis-password` (and
	// `redis-username` when the user is not `default`).
	ExternalRedisPasswordKey string
}{
	HelmChartName:            "argo-cd",
	HelmChartRepo:            "https://argoproj.github.io/argo-helm",
	DefaultChartVersion:      "10.2.1",
	InitialAdminSecretName:   "argocd-initial-admin-secret",
	FullnameBudget:           37,
	ExternalRedisPasswordKey: "redis-password",
}
