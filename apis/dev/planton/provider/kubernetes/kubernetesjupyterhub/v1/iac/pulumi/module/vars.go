package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_chart_name / helm_chart_repo): cross-engine chart-name drift
// deploys two different products from one manifest.
var vars = struct {
	HelmChartName string
	HelmChartRepo string
	// Fallback when spec.chart_version is unset — mirror of the proto
	// field's default option and the Terraform module's coalesce. Chart
	// 4.4.0 ships JupyterHub 5.5.0; the chart pin governs.
	DefaultChartVersion string
	// CHART-FIXED resource names: at the chart's default
	// fullnameOverride "" every resource renders a bare name — the
	// deployment is a per-NAMESPACE singleton and these names are the
	// exported composition handles.
	ProxyPublicServiceName string
	HubServiceName         string
	// The proxy-public Service's HTTP port (the chart's
	// service.ports.http) and the hub API port (hub Service).
	ProxyPublicPort int
	HubApiPort      int
	// Module-owned Secret names/keys. HubSecretSuffix is the
	// hub.existingSecret the hub mounts for the external database
	// password (chart contract key `hub.db.password`);
	// AuthSecretSuffix carries the shared sign-in password on the
	// dummy arm.
	HubSecretSuffix   string
	HubDbPasswordKey  string
	AuthSecretSuffix  string
	SharedPasswordKey string
	// Environment variable names the module's extraConfig snippets
	// read — the leak-free path for secret material into JupyterHub's
	// python configuration (everything in Helm values lands readable
	// inside the chart-owned hub Secret; env vars from secretKeyRef do
	// not).
	SharedPasswordEnvVar    string
	OauthClientSecretEnvVar string
	// Helm wait budget: with the pre-puller hook enabled (the chart
	// default) the install waits for the notebook image to pull onto
	// EVERY node before the hook-image-awaiter Job releases it — on a
	// fresh cluster this is the multi-GB pull budget, not the hub's.
	HelmTimeoutSeconds int
}{
	HelmChartName:           "jupyterhub",
	HelmChartRepo:           "https://hub.jupyter.org/helm-chart/",
	DefaultChartVersion:     "4.4.0",
	ProxyPublicServiceName:  "proxy-public",
	HubServiceName:          "hub",
	ProxyPublicPort:         80,
	HubApiPort:              8081,
	HubSecretSuffix:         "-hub-secret",
	HubDbPasswordKey:        "hub.db.password",
	AuthSecretSuffix:        "-auth",
	SharedPasswordKey:       "password",
	SharedPasswordEnvVar:    "PLANTON_SHARED_PASSWORD",
	OauthClientSecretEnvVar: "PLANTON_OAUTH_CLIENT_SECRET",
	HelmTimeoutSeconds:      1200,
}
