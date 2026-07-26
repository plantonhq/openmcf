package module

var vars = struct {
	// HelmChartName is the official kyverno chart.
	HelmChartName string

	// HelmChartRepo is the served chart index.
	HelmChartRepo string

	// DefaultChartVersion is the fallback when spec.chart_version is
	// unset AND the platform's defaulting middleware did not run. Keep
	// aligned with the spec default. The chart's appVersion pins every
	// controller image (chart 3.8.2 = Kyverno v1.18.2).
	DefaultChartVersion string

	// FullnameMaxLen is the fail-loud budget for metadata.name.
	// Chart-truth: controller Deployments derive from the CHART name
	// (constant "kyverno-..."), but the webhook Service ("<fullname>-svc"),
	// the runtime ConfigMap (the fullname) and the pre-delete hook Job
	// ("<fullname>-hook-pre-delete", the longest suffix at 16 chars) derive
	// from the fullname — past 47 the hook Job's name silently truncates at
	// the Kubernetes 63-char limit.
	FullnameMaxLen int
}{
	HelmChartName:       "kyverno",
	HelmChartRepo:       "https://kyverno.github.io/kyverno",
	DefaultChartVersion: "3.8.2",
	FullnameMaxLen:      47,
}
