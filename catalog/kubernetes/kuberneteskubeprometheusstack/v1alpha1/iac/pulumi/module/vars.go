package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 87.19.1 pairs with Prometheus Operator v0.92.1; the chart pin
	// governs.
	DefaultChartVersion string
	// Fallbacks for the storage arms — mirrors of the proto defaults.
	DefaultPrometheusDiskSize   string
	DefaultAlertmanagerDiskSize string
	DefaultGrafanaStorageSize   string
	// The chart SILENTLY truncates fullnameOverride at this many
	// characters (its own headroom for the longest child name it
	// derives). The modules pin the fullname to the resource name and
	// FAIL LOUDLY on longer names instead of letting the chart truncate —
	// a truncated fullname breaks the naming contract every exported
	// output is built on.
	FullnameBudget int
	// Suffix key names of the module-owned Secret that carries declared
	// remote-write basic-auth usernames. The Prometheus CRD reads BOTH
	// basic-auth halves from Secrets; usernames are not secrets, so the
	// spec accepts them as plain strings and the module materializes this
	// Secret (the declared-credentials pattern) rather than pushing a
	// pre-created Secret onto the user for a non-secret value.
	RemoteWriteAuthSecretSuffix string
}{
	HelmChartName:               "kube-prometheus-stack",
	HelmChartRepo:               "https://prometheus-community.github.io/helm-charts",
	DefaultChartVersion:         "87.19.1",
	DefaultPrometheusDiskSize:   "50Gi",
	DefaultAlertmanagerDiskSize: "2Gi",
	DefaultGrafanaStorageSize:   "10Gi",
	FullnameBudget:              26,
	RemoteWriteAuthSecretSuffix: "-remote-write-auth",
}
