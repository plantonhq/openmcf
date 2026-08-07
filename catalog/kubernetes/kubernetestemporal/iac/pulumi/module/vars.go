package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 1.6.0 ships Temporal 1.31.2; the chart pin governs.
	DefaultChartVersion string
	// The chart's fullname budget: child names are
	// `<fullname>-<component>`, and the chart's componentname helper
	// TRUNCATES THE FULLNAME (not the component) to fit 63 characters —
	// the longest component suffix is "internal-frontend" (17 chars), so
	// names past 62-17=45 characters silently lose fullname characters
	// and break the naming contract the exported outputs are built on.
	FullnameBudget int
	// Service ports (the chart's defaults — the typed spec does not
	// remap them; helm_values can).
	FrontendGrpcPort int
	FrontendHttpPort int
	WebPort          int
	// SQL driver plugin names the schema jobs and server config key off
	// (values.yaml documents postgres12/postgres12_pgx/mysql8; the
	// pgx driver is the maintained PostgreSQL path).
	PostgresPlugin string
	MysqlPlugin    string
	// Dynamic-config keys, verified against the server source at the pin
	// (common/dynamicconfig/constants.go).
	DcHistorySizeLimitError  string
	DcHistorySizeLimitWarn   string
	DcHistoryCountLimitError string
	DcHistoryCountLimitWarn  string
	DcBlobSizeLimitError     string
	DcBlobSizeLimitWarn      string
	// Helm wait budget: the pre-install schema Jobs run BEFORE the
	// release resources, then four server Deployments + the UI must roll
	// out; a cold install against a fresh database fits comfortably.
	HelmTimeoutSeconds int
}{
	HelmChartName:            "temporal",
	HelmChartRepo:            "https://go.temporal.io/helm-charts",
	DefaultChartVersion:      "1.6.0",
	FullnameBudget:           45,
	FrontendGrpcPort:         7233,
	FrontendHttpPort:         7243,
	WebPort:                  8080,
	PostgresPlugin:           "postgres12",
	MysqlPlugin:              "mysql8",
	DcHistorySizeLimitError:  "limit.historySize.error",
	DcHistorySizeLimitWarn:   "limit.historySize.warn",
	DcHistoryCountLimitError: "limit.historyCount.error",
	DcHistoryCountLimitWarn:  "limit.historyCount.warn",
	DcBlobSizeLimitError:     "limit.blobSize.error",
	DcBlobSizeLimitWarn:      "limit.blobSize.warn",
	HelmTimeoutSeconds:       900,
}
