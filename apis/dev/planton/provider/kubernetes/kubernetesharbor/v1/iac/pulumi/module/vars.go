package module

// vars carries the module's pinned chart identity and naming contracts.
// The chart version is the SERVED index truth (https://helm.goharbor.io)
// — chart 1.19.1 pairs with Harbor 2.15.1 (the chart's appVersion pins
// every component image tag).
var vars = struct {
	HelmChartName       string
	HelmChartRepo       string
	DefaultChartVersion string

	// Suffixes of the module-owned Secrets. `-admin-auth` is the
	// exported credential handle (key HARBOR_ADMIN_PASSWORD);
	// `-internal-auth` carries every generated inter-component
	// credential under the chart's per-site contract keys;
	// `-redis-auth` and `-storage-auth` materialize declared external
	// credentials so nothing credential-bearing rides rendered values.
	AdminAuthSecretSuffix    string
	InternalAuthSecretSuffix string
	RedisAuthSecretSuffix    string
	StorageAuthSecretSuffix  string

	// The chart's contract key for the admin password Secret.
	AdminPasswordSecretKey string

	// The internal registry basic-auth username (chart default kept:
	// the username is not secret material — the module-generated
	// PASSWORD is what replaces the chart's public default).
	RegistryCredentialUsername string

	// Harbor's admin login is a server constant.
	AdminUsername string

	// MaxNameLength is the fullname budget, derived from the chart's
	// name helpers: the chart truncates its fullname at 63 and then
	// APPENDS component suffixes. The longest derived object name is
	// `<fullname>-jobservice-internal-tls` (24 characters of suffix,
	// rendered whenever internalTLS runs in auto mode), so names cap
	// at 63 − 24 = 39. Enforced fail-loud in BOTH engines.
	MaxNameLength int
}{
	HelmChartName:       "harbor",
	HelmChartRepo:       "https://helm.goharbor.io",
	DefaultChartVersion: "1.19.1",

	AdminAuthSecretSuffix:    "-admin-auth",
	InternalAuthSecretSuffix: "-internal-auth",
	RedisAuthSecretSuffix:    "-redis-auth",
	StorageAuthSecretSuffix:  "-storage-auth",

	AdminPasswordSecretKey: "HARBOR_ADMIN_PASSWORD",

	RegistryCredentialUsername: "harbor_registry_user",

	AdminUsername: "admin",

	MaxNameLength: 39,
}
