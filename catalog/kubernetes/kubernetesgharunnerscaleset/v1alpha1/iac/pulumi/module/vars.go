package module

var vars = struct {
	// HelmOciRepo is the OCI registry path holding both ARC charts.
	// Pulumi's helm.v3.Release does not resolve oci:// through
	// RepositoryOpts the way the Terraform provider does — the chart
	// reference must be the JOINED "<repo>/<chart>" string (see main.go).
	HelmOciRepo string

	// HelmChartName is the runner scale set chart.
	HelmChartName string

	// DefaultChartVersion is the fallback when spec.chart_version is
	// unset AND the platform's defaulting middleware did not run. Keep
	// aligned with the spec default — and EQUAL to the controller
	// kind's: GitHub supports controller and scale set charts only at
	// matching versions.
	DefaultChartVersion string

	// GithubAuthSecretSuffix names the module-materialized credential
	// Secret (`<name>-github-auth`) for the declared PAT / GitHub App
	// arms. The existing-Secret arm references the user's own Secret
	// instead.
	GithubAuthSecretSuffix string

	// ScaleSetNameBudget is the chart's own hard cap on the runner scale
	// set name (the chart template fails past 45 characters — a GitHub
	// registration limit).
	ScaleSetNameBudget int

	// DefaultRunnerImage is the chart's default runner container image.
	DefaultRunnerImage string

	// RunnerCommand is the runner container's fixed entrypoint (the
	// chart's own default container renders it; any container override
	// must re-state it because Helm values LISTS replace, never merge).
	RunnerCommand string
}{
	HelmOciRepo:            "oci://ghcr.io/actions/actions-runner-controller-charts",
	HelmChartName:          "gha-runner-scale-set",
	DefaultChartVersion:    "0.14.2",
	GithubAuthSecretSuffix: "-github-auth",
	ScaleSetNameBudget:     45,
	DefaultRunnerImage:     "ghcr.io/actions/actions-runner:latest",
	RunnerCommand:          "/home/runner/run.sh",
}
