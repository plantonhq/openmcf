package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 4.40.0 ships appVersion "4.40"; the chart pin governs.
	DefaultChartVersion string
	// Per-tier data-volume size fallbacks — mirrors of the proto
	// documentation. Master and filer hold metadata (small); the volume
	// tier holds the object bytes and is the one to size deliberately.
	DefaultMasterVolumeSize string
	DefaultVolumeVolumeSize string
	DefaultFilerVolumeSize  string
	DefaultAdminVolumeSize  string
}{
	HelmChartName:           "seaweedfs",
	HelmChartRepo:           "https://seaweedfs.github.io/seaweedfs/helm",
	DefaultChartVersion:     "4.40.0",
	DefaultMasterVolumeSize: "5Gi",
	DefaultVolumeVolumeSize: "30Gi",
	DefaultFilerVolumeSize:  "10Gi",
	DefaultAdminVolumeSize:  "10Gi",
}
