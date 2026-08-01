package module

var vars = struct {
	// HelmChartRepo is the official ASF chart repository serving the
	// spark-kubernetes-operator chart.
	HelmChartRepo string
	// HelmChartName is the Apache Spark Kubernetes operator chart
	// ("spark-kubernetes-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. 1.8.0 is the newest SERVED chart (= operator
	// appVersion 1.0.0, verified against the repository index). The
	// spark.apache.org CRDs ship from the chart's crds/ directory: Helm
	// installs them once and NEVER upgrades them — bumping this version
	// does not touch the CRDs (apply the new release's CRD files manually
	// when a bump changes them).
	DefaultChartVersion string

	// OperatorImagePath is the operator image's repository path WITHOUT
	// the registry — the half image_registry replaces (the air-gap mirror
	// seam). Keep in lockstep with the Terraform module's operator_image.
	OperatorImagePath string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls plus the operator JVM's 30s-initial-delay
	// startup probe; atomic rolls back on expiry so a wedged install
	// never lingers half-deployed.
	HelmTimeoutSeconds int
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepo:       "https://apache.github.io/spark-kubernetes-operator",
	HelmChartName:       "spark-kubernetes-operator",
	DefaultChartVersion: "1.8.0",
	OperatorImagePath:   "apache/spark-kubernetes-operator",
	HelmTimeoutSeconds:  600,
}
