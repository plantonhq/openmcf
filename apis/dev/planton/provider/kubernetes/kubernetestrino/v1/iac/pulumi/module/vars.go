package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Chart 1.42.2 ships Trino 480; the chart pin governs.
	ChartVersion string
	// Name suffixes the chart derives from the fullname (the module
	// PINS fullnameOverride to metadata.name, so every child name is
	// deterministic): the coordinator/worker Deployments and Services
	// and the catalog ConfigMap.
	CoordinatorSuffix string
	WorkerSuffix      string
	// The coordinator HTTP port (service.port — REST API + Web UI).
	HttpPort int
	// Module-owned Secret names/keys. AuthSecretSuffix holds the
	// bootstrap admin credential: key `password` is the plaintext
	// clients use; key `password.db` is the htpasswd-format bcrypt
	// file the chart mounts through auth.passwordAuthSecret (the two
	// keys are a VERIFIED pairing — generated from one random).
	// InternalSecretSuffix holds the internal-communication shared
	// secret every node must present once authentication is on.
	AuthSecretSuffix     string
	AdminPasswordKey     string
	PasswordDbKey        string
	InternalSecretSuffix string
	SharedSecretKey      string
	// Environment variable names referenced from rendered properties
	// via Trino's `${ENV:VAR}` secrets substitution — the leak-free
	// path for secret material into ConfigMap-rendered configuration.
	SharedSecretEnvVar string
	// Per-catalog password env vars render as
	// CatalogPasswordEnvPrefix + upper(catalog_name) + "_PASSWORD".
	CatalogPasswordEnvPrefix string
	// Name-budget truths derived from the chart templates at the pin:
	// the schemas ConfigMap `<fullname>-schemas-volume-coordinator`
	// renders UNCONDITIONALLY (27-char suffix) and the resource-groups
	// ConfigMap `<fullname>-resource-groups-volume-coordinator`
	// (36-char suffix) renders only when resource groups use the
	// configmap type — both must fit 63.
	NameBudget               int
	NameBudgetResourceGroups int
	// Helm wait budget: coordinator+worker rollouts plus the first
	// image pull (~1GB trinodb/trino) on a fresh node.
	HelmTimeoutSeconds int
}{
	HelmChartName:            "trino",
	HelmChartRepo:            "https://trinodb.github.io/charts",
	ChartVersion:             "1.42.2",
	CoordinatorSuffix:        "-coordinator",
	WorkerSuffix:             "-worker",
	HttpPort:                 8080,
	AuthSecretSuffix:         "-auth",
	AdminPasswordKey:         "password",
	PasswordDbKey:            "password.db",
	InternalSecretSuffix:     "-internal",
	SharedSecretKey:          "shared-secret",
	SharedSecretEnvVar:       "TRINO_INTERNAL_SHARED_SECRET",
	CatalogPasswordEnvPrefix: "TRINO_CATALOG_",
	NameBudget:               36,
	NameBudgetResourceGroups: 27,
	HelmTimeoutSeconds:       900,
}
