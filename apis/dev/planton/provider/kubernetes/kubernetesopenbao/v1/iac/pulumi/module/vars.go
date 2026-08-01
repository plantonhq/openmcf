package module

// vars carries the module's pinned chart identity. The chart version is
// the SERVED index truth (https://openbao.github.io/openbao-helm) — chart
// 0.28.6 pairs with OpenBao v2.6.1 (the chart's appVersion pins the
// server image tag).
var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string

	// OpenBao's API and cluster ports (chart constants — the listener
	// template binds [::]:8200 / [::]:8201).
	ApiPort     int
	ClusterPort int

	// Data/audit mount paths (chart constants — the config's storage
	// stanzas and the PVC mounts must agree).
	DataMountPath  string
	AuditMountPath string

	// Where the module mounts the TLS Secret when tls.enabled — the
	// listener's tls_cert_file/tls_key_file point here.
	TlsMountPath string

	// Suffix of the module-owned Secret carrying declared seal
	// credentials (AWS secret key / Azure client secret / transit
	// token), delivered to the server as environment variables so
	// nothing credential-bearing lands in the config ConfigMap.
	SealCredentialsSecretSuffix string

	// Name-length budgets, derived from the chart's name helpers: the
	// chart truncates its fullname at 63 but appends suffixes AFTER
	// truncation. Longest Service suffixes: `-internal` (9) always,
	// `-agent-injector-svc` (19) when the injector is enabled — Services
	// cap at 63 characters.
	MaxNameLength             int
	MaxNameLengthWithInjector int
}{
	HelmChartName:       "openbao",
	HelmChartRepo:       "https://openbao.github.io/openbao-helm",
	DefaultChartVersion: "0.28.6",

	ApiPort:     8200,
	ClusterPort: 8201,

	DataMountPath:  "/openbao/data",
	AuditMountPath: "/openbao/audit",
	TlsMountPath:   "/openbao/tls",

	SealCredentialsSecretSuffix: "-seal-credentials",

	MaxNameLength:             54,
	MaxNameLengthWithInjector: 44,
}
