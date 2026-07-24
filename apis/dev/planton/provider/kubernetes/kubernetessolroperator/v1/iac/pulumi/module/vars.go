package module

var vars = struct {
	// HelmChartRepo is the Apache Solr project's chart repository.
	HelmChartRepo string
	// HelmChartName is the Solr operator chart ("solr-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart versions carry NO `v` prefix (the
	// operator image/CRD artifacts DO — chart 0.9.1 ships operator
	// v0.9.1); the SERVED index (https://solr.apache.org/charts) governs.
	DefaultChartVersion string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int

	// CrdsDir is where the module-owned CRD files are staged, relative to
	// the Pulumi project directory (the program's working directory —
	// Pulumi always runs the program with cwd = the directory holding
	// Pulumi.yaml, so ../crds is <kind>/v1/iac/crds). The solr-operator
	// chart ships NO CRDs (they are separate release artifacts); the
	// module owns all four: the three solr.apache.org CRDs plus the
	// ZookeeperCluster CRD of the bundled zookeeper-operator dependency.
	CrdsDir string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://solr.apache.org/charts",
	HelmChartName:       "solr-operator",
	DefaultChartVersion: "0.9.1",
	HelmTimeoutSeconds:  600,
	CrdsDir:             "../crds",
}
