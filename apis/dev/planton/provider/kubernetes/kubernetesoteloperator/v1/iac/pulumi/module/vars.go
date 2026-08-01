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
	// index). Bumping it requires re-staging the ../crds files from the
	// new pin — the staged files ARE the chart's CRDs at this version.
	DefaultChartVersion string

	// CrdDirectory is where the module-owned CRD files are staged,
	// relative to the Pulumi project directory (the working directory at
	// program run time). The staged files are TOKENIZED renders of the
	// pinned chart's templated CRDs — see crds.go for the substitution
	// contract.
	CrdDirectory string

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
	HelmChartRepo:       "https://open-telemetry.github.io/opentelemetry-helm-charts",
	HelmChartName:       "opentelemetry-operator",
	DefaultChartVersion: "0.120.0",
	CrdDirectory:        "../crds",
	ManagerImagePath:    "open-telemetry/opentelemetry-operator/opentelemetry-operator",
	HelmTimeoutSeconds:  600,
}
