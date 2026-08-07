package module

var vars = struct {
	// HelmOciRepo is the OCI registry path holding both ARC charts.
	// Pulumi's helm.v3.Release does not resolve oci:// through
	// RepositoryOpts the way the Terraform provider does — the chart
	// reference must be the JOINED "<repo>/<chart>" string (see main.go).
	HelmOciRepo string

	// HelmChartName is the controller chart.
	HelmChartName string

	// DefaultChartVersion is the fallback when spec.chart_version is
	// unset AND the platform's defaulting middleware did not run. Keep
	// aligned with the spec default. Chart and controller image move in
	// lockstep (chart 0.14.2 = appVersion 0.14.2).
	DefaultChartVersion string
}{
	HelmOciRepo:         "oci://ghcr.io/actions/actions-runner-controller-charts",
	HelmChartName:       "gha-runner-scale-set-controller",
	DefaultChartVersion: "0.14.2",
}
