package module

var vars = struct {
	// HelmChartName is the official gatekeeper chart.
	HelmChartName string

	// HelmChartRepo is the served chart index.
	HelmChartRepo string

	// DefaultChartVersion is the fallback when spec.chart_version is
	// unset AND the platform's defaulting middleware did not run. Keep
	// aligned with the spec default. Chart and app versions move in
	// lockstep (chart 3.23.0 = Gatekeeper v3.23.0).
	DefaultChartVersion string

	// WebhookServiceName and WebhookCertSecretName are CHART-FIXED names
	// (the chart hardcodes them — no fullname derivation): the Service
	// the webhook configurations point at, and the Secret the embedded
	// cert-controller rotates. Exported as outputs.
	WebhookServiceName    string
	WebhookCertSecretName string
}{
	HelmChartName:         "gatekeeper",
	HelmChartRepo:         "https://open-policy-agent.github.io/gatekeeper/charts",
	DefaultChartVersion:   "3.23.0",
	WebhookServiceName:    "gatekeeper-webhook-service",
	WebhookCertSecretName: "gatekeeper-webhook-server-cert",
}
