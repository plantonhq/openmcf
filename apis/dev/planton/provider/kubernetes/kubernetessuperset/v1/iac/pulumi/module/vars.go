package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Chart 0.22.4 ships Superset 6.1.0; the chart pin governs.
	ChartVersion string
	// The web Service port (the chart's service.port).
	HttpPort int
	// Name suffixes the chart derives from the fullname (the module
	// PINS fullnameOverride to metadata.name): the env/config Secrets
	// the chart consumes by name, and the longest child
	// (`<name>-celerybeat`) that sets the name budget.
	EnvSecretSuffix string
	// Module-owned Secret names/keys. AdminAuthSuffix holds the
	// bootstrap admin password (the exported handle); SecretKeySuffix
	// holds the Flask session-signing key (STABLE — rotation logs out
	// every session and orphans encrypted datasource credentials).
	AdminAuthSuffix  string
	AdminPasswordKey string
	SecretKeySuffix  string
	SecretKeyKey     string
	// Environment variable names: the ws server reads JwtSecretEnvVar
	// from its environment natively (config.ts env priority); the
	// module's configOverrides snippets read the same variables via
	// env() in superset_config.py.
	JwtSecretEnvVar string
	// The chart mounts superset_config.py + the init/bootstrap scripts
	// at this path (the chart's configMountPath default).
	ConfigMountPath string
	// Longest chart-derived child suffix (`-celerybeat`, 11 chars) —
	// every child name must fit 63.
	NameBudget int
	// Helm wait budget: web+worker rollouts plus the init Job's schema
	// migration against the composed database (the post-install hook
	// runs inside the release operation) and the first image pull.
	HelmTimeoutSeconds int
}{
	HelmChartName:      "superset",
	HelmChartRepo:      "https://apache.github.io/superset",
	ChartVersion:       "0.22.4",
	HttpPort:           8088,
	EnvSecretSuffix:    "-env",
	AdminAuthSuffix:    "-admin-auth",
	AdminPasswordKey:   "password",
	SecretKeySuffix:    "-secret-key",
	SecretKeyKey:       "secret_key",
	JwtSecretEnvVar:    "JWT_SECRET",
	ConfigMountPath:    "/app/pythonpath",
	NameBudget:         52,
	HelmTimeoutSeconds: 900,
}
