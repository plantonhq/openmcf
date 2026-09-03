package module

var vars = struct {
	// HelmChartRepo is the Apache Solr project's chart repository.
	HelmChartRepo string
	// HelmChartName is the Solr operator chart ("solr-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart versions carry NO `v` prefix (the
	// operator image artifacts DO — chart 0.9.1 ships operator v0.9.1);
	// the SERVED index (https://solr.apache.org/charts) governs. The CRDs
	// are derived from whatever version is pinned, so a bump here changes
	// exactly one line.
	DefaultChartVersion string

	// CrdRenderOverride is the one values document merged LAST when the
	// pinned chart is rendered to derive its CRDs: the bundled
	// zookeeper-operator subchart's CRD switch turned on, so the
	// ZookeeperCluster CRD it templates joins the three solr.apache.org
	// CRDs from the chart's crds/ directory. The release itself installs
	// with the switch off and skip_crds set, so this never reaches the
	// cluster through Helm. Twin of the Terraform module's
	// helm_crds_args.render_override.
	CrdRenderOverride string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls on cold clusters; atomic rolls back on
	// expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://solr.apache.org/charts",
	HelmChartName:       "solr-operator",
	DefaultChartVersion: "0.9.1",
	CrdRenderOverride:   "zookeeper-operator:\n  crd:\n    create: true\n",
	HelmTimeoutSeconds:  600,
}
