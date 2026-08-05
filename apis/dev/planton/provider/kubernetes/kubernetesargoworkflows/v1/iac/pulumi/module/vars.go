package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 1.0.23 ships Argo Workflows v4.0.8; the chart pin governs.
	DefaultChartVersion string
	// Fallback for spec.workflow_service_account — the chart's own
	// default runner ServiceAccount name.
	DefaultWorkflowServiceAccount string
	// The Argo server's port (Service and container).
	ServerPort int
	// The chart's fullname budget: every child name is
	// `<fullname>-<component>` truncated at 63 characters, and the longest
	// component suffix is "-workflow-controller" (20 chars) — names past
	// 63-20=43 characters truncate SILENTLY and break the naming contract
	// the exported outputs are built on.
	FullnameBudget int
	// Fallbacks for the s3 credentials Secret's key names when the spec
	// leaves them unset — the KubernetesSeaweedFs generated `-s3-secret`
	// convention (mirror of the proto field defaults), so that Secret
	// composes with zero key configuration. The chart takes the key names
	// verbatim through its accessKeySecret/secretKeySecret selectors.
	S3AccessKeyIdKeyDefault     string
	S3SecretAccessKeyKeyDefault string
	// The chart's secret-selector key contracts for the gcs/azure
	// artifact-store credentials (its own documented shapes).
	GcsServiceAccountKey string
	AzureAccountKeyKey   string
	// The chart's archive credential-selector defaults (userNameSecret /
	// passwordSecret keys when the spec leaves them unset).
	ArchiveUsernameKey string
	ArchivePasswordKey string
	// The archive table name Argo Workflows documents as its default.
	ArchiveTableName string
}{
	HelmChartName:                 "argo-workflows",
	HelmChartRepo:                 "https://argoproj.github.io/argo-helm",
	DefaultChartVersion:           "1.0.23",
	DefaultWorkflowServiceAccount: "argo-workflow",
	ServerPort:                    2746,
	FullnameBudget:                43,
	S3AccessKeyIdKeyDefault:       "admin_access_key_id",
	S3SecretAccessKeyKeyDefault:   "admin_secret_access_key",
	GcsServiceAccountKey:          "serviceAccountKey",
	AzureAccountKeyKey:            "account-access-key",
	ArchiveUsernameKey:            "username",
	ArchivePasswordKey:            "password",
	ArchiveTableName:              "argo_workflows",
}
