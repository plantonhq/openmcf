package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 1.22.0 ships Airflow 3.2.2; the chart pin governs.
	DefaultChartVersion string
	// Fallback when spec.airflow_version is unset — the chart pin's
	// appVersion. The module keeps `airflowVersion` (template gating)
	// and `defaultAirflowTag` (image tag) in lockstep from this value.
	DefaultAirflowVersion string
	// Fallback when spec.executor is unset — mirror of the proto
	// default. This kind defaults to the Kubernetes-native executor
	// (no broker dependency); the chart's own default is
	// CeleryExecutor.
	DefaultExecutor string
	// The chart's fullname budget: at the default useStandardNaming
	// false the fullname IS the release name (= metadata.name), and
	// child names append fixed suffixes to it. The longest
	// always-renderable suffix is "-run-airflow-migrations" (23
	// characters — the migration Job), so names past 63-23=40
	// characters produce children past the Kubernetes 63-character
	// name limit and the deploy fails midway with API rejections.
	FullnameBudget int
	// The API server's HTTP port (the chart's ports.apiServer).
	ApiServerPort int
	// The bundled Redis client port and PgBouncer listen port (the
	// chart's ports.redisDB / ports.pgbouncer).
	RedisPort     int
	PgbouncerPort int
	// Module-owned Secret name suffixes. The connection Secrets carry
	// the chart's `connection`-key contract; the key Secrets carry the
	// chart's own per-Secret key names (fernet-key, api-secret-key,
	// webserver-secret-key, jwt-secret).
	MetadataConnSecretSuffix      string
	ResultBackendConnSecretSuffix string
	BrokerUrlSecretSuffix         string
	RedisPasswordSecretSuffix     string
	LogReadConnSecretSuffix       string
	PgbouncerConfigSecretSuffix   string
	PgbouncerStatsSecretSuffix    string
	AdminAuthSecretSuffix         string
	FernetKeySecretSuffix         string
	ApiSecretKeySecretSuffix      string
	WebserverSecretKeySuffix      string
	JwtSecretSuffix               string
	// Chart contract keys inside the Secrets. KedaConnectionKey is the
	// metadata Secret's second key — the direct-database form the
	// chart's KEDA autoscalers read (env KEDA_DB_CONN) on mysql and
	// pgbouncer-bypass postures.
	ConnectionKey         string
	KedaConnectionKey     string
	FernetKeyKey          string
	FernetKeyStdB64Key    string
	ApiSecretKeyKey       string
	WebserverSecretKeyKey string
	JwtSecretKey          string
	AdminPasswordKey      string
	// Helm wait budget: the post-install migration + create-user hook
	// Jobs run inside the wait, then the API server, scheduler,
	// dag-processor, triggerer (and workers on Celery) must roll out —
	// a cold install against a fresh database migrates dozens of
	// schema revisions first.
	HelmTimeoutSeconds int
}{
	HelmChartName:                 "airflow",
	HelmChartRepo:                 "https://airflow.apache.org",
	DefaultChartVersion:           "1.22.0",
	DefaultAirflowVersion:         "3.2.2",
	DefaultExecutor:               "KubernetesExecutor",
	FullnameBudget:                40,
	ApiServerPort:                 8080,
	RedisPort:                     6379,
	PgbouncerPort:                 6543,
	MetadataConnSecretSuffix:      "-metadata-conn",
	ResultBackendConnSecretSuffix: "-result-backend-conn",
	BrokerUrlSecretSuffix:         "-broker-url",
	RedisPasswordSecretSuffix:     "-redis-password",
	LogReadConnSecretSuffix:       "-log-read-conn",
	PgbouncerConfigSecretSuffix:   "-pgbouncer-config",
	PgbouncerStatsSecretSuffix:    "-pgbouncer-stats",
	AdminAuthSecretSuffix:         "-admin-auth",
	FernetKeySecretSuffix:         "-fernet-key",
	ApiSecretKeySecretSuffix:      "-api-secret-key",
	WebserverSecretKeySuffix:      "-webserver-secret-key",
	JwtSecretSuffix:               "-jwt-secret",
	ConnectionKey:                 "connection",
	KedaConnectionKey:             "kedaConnection",
	FernetKeyKey:                  "fernet-key",
	FernetKeyStdB64Key:            "fernet-key-std-b64",
	ApiSecretKeyKey:               "api-secret-key",
	WebserverSecretKeyKey:         "webserver-secret-key",
	JwtSecretKey:                  "jwt-secret",
	AdminPasswordKey:              "password",
	HelmTimeoutSeconds:            900,
}
