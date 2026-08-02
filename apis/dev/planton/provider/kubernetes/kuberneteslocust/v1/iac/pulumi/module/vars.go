package module

// Chart identity — must stay byte-identical with the Terraform module's
// locals (helm_oci_repo / helm_chart_name): cross-engine chart-name drift
// deploys two different products from one manifest.
//
// SERVING TRUTH: the chart's home is the OCI registry — the classic
// index at charts.deliveryhero.io stalls at 0.31.6 (2024) while
// ghcr.io serves the live line. Pulumi's helm.v3.Release resolves OCI
// charts through the joined "oci://<registry-path>/<chart>" string with
// NO RepositoryOpts; the Terraform provider instead takes
// repository = the registry path plus the bare chart name.
var vars = struct {
	HelmOciRepo   string
	HelmChartName string
	// Chart 0.35.0 ships Locust 2.32.2; the chart pin governs.
	ChartVersion string
	// The official image and the release this kind is built against.
	DefaultImageRepository string
	DefaultImageTag        string
	// Name suffixes the chart derives from the fullname (the module
	// PINS fullnameOverride to metadata.name, so every child name is
	// deterministic): the master/worker Deployments and per-component
	// ServiceAccounts; the master Service is the bare fullname.
	MasterSuffix string
	WorkerSuffix string
	// Module-owned satellite name suffixes. LocustfileSuffix is the
	// longest derived name (11 chars) — it sets the name budget.
	LocustfileSuffix  string
	LibSuffix         string
	WebAuthCodeSuffix string
	AuthSecretSuffix  string
	// Keys inside the `<name>-auth` Secret. The username rides the
	// Secret too so the login backend reads ONE mounted source of
	// truth (and rotating the credential is a single-Secret edit).
	AuthUsernameKey       string
	AuthPasswordKey       string
	AuthFlaskSecretKeyKey string
	// In-pod mount paths: the login-backend code (a ConfigMap via the
	// chart's extraConfigMaps seam) and the credential files (the
	// chart's mount_external_secret projected volume).
	WebAuthCodeMountPath   string
	WebAuthSecretMountPath string
	// The chart's locustfile mount path (chart default, pinned
	// explicitly — the master's -f argument is derived from it).
	LocustfileMountPath string
	// The web-UI/REST port and the worker-connect port.
	WebPort        int
	MasterBindPort int
	// The pod-template annotation carrying the module's content hash
	// — the chart checksums only its OWN ConfigMaps, so module-owned
	// script changes roll the pods through this annotation instead.
	ChecksumAnnotation string
	// The web-UI login floor: the chart renders the modern
	// `--web-login` flag only for image tags >= 2.21.0; BELOW it the
	// chart falls onto `--web-auth=user:password` — credentials as a
	// LITERAL POD ARGUMENT, which this module refuses to render.
	AuthMinMajor int
	AuthMinMinor int
	// Name budget: the longest chart/module-derived child name is
	// `<name>-locustfile` (11-char suffix) — 63 - 11 = 52.
	NameBudget int
	// Helm wait budget: master + worker rollouts and the image pull;
	// pip_packages installs at pod start ride the readiness window.
	HelmTimeoutSeconds int
}{
	HelmOciRepo:            "oci://ghcr.io/deliveryhero/helm-charts",
	HelmChartName:          "locust",
	ChartVersion:           "0.35.0",
	DefaultImageRepository: "locustio/locust",
	DefaultImageTag:        "2.32.2",
	MasterSuffix:           "-master",
	WorkerSuffix:           "-worker",
	LocustfileSuffix:       "-locustfile",
	LibSuffix:              "-lib",
	WebAuthCodeSuffix:      "-web-auth",
	AuthSecretSuffix:       "-auth",
	AuthUsernameKey:        "username",
	AuthPasswordKey:        "password",
	AuthFlaskSecretKeyKey:  "flask-secret-key",
	WebAuthCodeMountPath:   "/opt/planton/web-auth-code",
	WebAuthSecretMountPath: "/opt/planton/web-auth",
	LocustfileMountPath:    "/mnt/locust",
	WebPort:                8089,
	MasterBindPort:         5557,
	ChecksumAnnotation:     "planton.dev/config-checksum",
	AuthMinMajor:           2,
	AuthMinMinor:           21,
	NameBudget:             52,
	HelmTimeoutSeconds:     600,
}
