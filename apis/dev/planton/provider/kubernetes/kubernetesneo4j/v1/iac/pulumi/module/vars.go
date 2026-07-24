package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// versions track Neo4j calendar releases (2026.6.0 ships Neo4j
	// 2026.06); the chart pin governs.
	DefaultChartVersion string
	// Fallback for spec.data_volume.size — mirror of the proto field's
	// default option. The chart REQUIRES volumes.data.mode, so the module
	// always renders the data volume block, defaulted or not.
	DefaultDataVolumeSize string
}{
	HelmChartName:         "neo4j",
	HelmChartRepo:         "https://helm.neo4j.com/neo4j",
	DefaultChartVersion:   "2026.6.0",
	DefaultDataVolumeSize: "10Gi",
}
