package module

import "fmt"

var vars = struct {
	// HelmChartRepoTemplate is the fmt template for the versioned ASF
	// downloads directory serving the chart — the chart is served PER
	// VERSION, so the repository URL itself carries the resolved chart
	// version (built in locals.go with fmt.Sprintf; the Terraform twin
	// interpolates the same URL inline). Mechanically chosen over a
	// fixed HelmChartRepo because no version-independent repository
	// exists for this chart.
	HelmChartRepoTemplate string
	// HelmChartName is the Flink Kubernetes operator chart
	// ("flink-kubernetes-operator").
	HelmChartName string
	// DefaultChartVersion is the fallback when spec.chart_version is unset
	// AND the platform's defaulting middleware did not run. Keep aligned
	// with the spec default. Chart version = operator version = the image
	// tag the module pins. The flink.apache.org CRDs ship from the
	// chart's crds/ directory: Helm installs them once and NEVER upgrades
	// them — bumping this version does not touch the CRDs (apply the new
	// release's CRD files manually when a bump changes them).
	DefaultChartVersion string

	// OperatorImagePath is the operator image's repository path WITHOUT
	// the registry — the half image_registry replaces (the air-gap
	// mirror seam). Keep in lockstep with the Terraform module's
	// operator_image.
	OperatorImagePath string

	// HelmTimeoutSeconds bounds the atomic install/upgrade of the release.
	// 600s covers image pulls plus (webhook arm) the cert-manager
	// certificate issuance the webhook container mounts; atomic rolls
	// back on expiry so a wedged install never lingers half-deployed.
	HelmTimeoutSeconds int

	// WebhookServiceName is the operator's webhook Service name —
	// CHART-FIXED (templates/webhook/service.yaml hardcodes it), not
	// fullname-derived.
	WebhookServiceName string
}{
	// Chart identity — MUST be identical in the Terraform module's locals
	// (cross-engine chart drift installs different software per engine).
	HelmChartRepoTemplate: "https://downloads.apache.org/flink/flink-kubernetes-operator-%s/",
	HelmChartName:         "flink-kubernetes-operator",
	DefaultChartVersion:   "1.15.0",
	OperatorImagePath:     "apache/flink-kubernetes-operator",
	HelmTimeoutSeconds:    600,
	WebhookServiceName:    "flink-operator-webhook-service",
}

// helmChartRepo builds the versioned repository URL for the resolved
// chart version.
func helmChartRepo(chartVersion string) string {
	return fmt.Sprintf(vars.HelmChartRepoTemplate, chartVersion)
}
