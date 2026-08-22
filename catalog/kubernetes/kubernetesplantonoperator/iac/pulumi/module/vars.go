package module

var vars = struct {
	// HelmOciRepo is the OCI registry path holding the planton-operator
	// chart. Pulumi's helm.v3.Release does not resolve oci:// through
	// RepositoryOpts the way the Terraform provider does — the chart
	// reference must be the JOINED "<repo>/<chart>" string (see main.go).
	HelmOciRepo string
	// HelmChartName is the operator chart ("planton-operator").
	HelmChartName string
	// DefaultChartVersion is the chart this catalog release was validated
	// against — the version installed when spec.chart_version is unset.
	// The staged CRD at ../crds is extracted from EXACTLY this published
	// chart package, so the CRD and the operator it schedules always
	// match; bumping this pin requires re-staging the CRD file in the
	// same change.
	DefaultChartVersion string
	// MinChartVersion is the schema-contract floor: charts below 0.7.0
	// ship operators whose reconcilers predate the PlantonPlatform
	// schema the staged CRD advertises — the API server would ACCEPT
	// fields the running operator silently ignores, the worst failure
	// shape. Refused loudly in both engines instead.
	MinChartVersion string
	// ReleaseName is FIXED: the operator enforces one installation per
	// cluster itself at startup (a label-matched Deployment scan that
	// refuses to start beside a sibling), so the release name never
	// derives from metadata.name — a second differently-named release
	// would only produce a crash-looping manager, never a second
	// operator.
	ReleaseName string
	// CrdDirectory holds the module-owned plantonplatforms.planton.ai
	// CRD, staged from the published chart at DefaultChartVersion
	// (relative to the Pulumi project dir). The chart's own crds/
	// install is always skipped — module ownership is what makes
	// chart_version upgrades carry the CRD and keep-on-uninstall
	// guaranteed rather than incidental.
	CrdDirectory string
	// ExpectedCrdCount fails the install loudly when the staged files
	// did not travel with the module (an empty directory would silently
	// apply nothing). Exactly one CRD at chart 0.7.1 — update together
	// with DefaultChartVersion when re-staging.
	ExpectedCrdCount int
	// HelmTimeoutSeconds bounds the atomic install/upgrade. 600s covers
	// image pulls on cold clusters; atomic rolls back on expiry so a
	// wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmOciRepo:         "oci://ghcr.io/plantonhq/charts",
	HelmChartName:       "planton-operator",
	DefaultChartVersion: "0.7.1",
	MinChartVersion:     "0.7.0",
	ReleaseName:         "planton-operator",
	CrdDirectory:        "../crds",
	ExpectedCrdCount:    1,
	HelmTimeoutSeconds:  600,
}
