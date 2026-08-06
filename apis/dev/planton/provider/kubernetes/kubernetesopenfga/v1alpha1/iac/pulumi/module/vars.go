package module

var vars = struct {
	// HelmChartName is the official openfga chart.
	HelmChartName string

	// HelmChartRepo is the served chart index.
	HelmChartRepo string

	// DefaultChartVersion is the fallback when spec.chart_version is
	// unset AND the platform's defaulting middleware did not run. Keep
	// aligned with the spec default. Chart and app versions move in
	// lockstep (chart 0.3.10 = OpenFGA v1.18.1).
	DefaultChartVersion string

	// MaxMetadataNameLength is the NAME BUDGET: the chart truncates its
	// fullname at 63 characters and separately derives `<fullname>-migrate`
	// for the migration Job — whose pod label value also caps at 63.
	// 63 - len("-migrate") = 55; a longer resource name would silently
	// truncate the fullname while the Job suffix pushes past the label
	// limit. Both engines fail loudly instead (Terraform twin: a
	// lifecycle precondition on the helm_release).
	MaxMetadataNameLength int

	// HttpPort / GrpcPort are the chart's fixed Service ports
	// (http.addr 0.0.0.0:8080, grpc.addr 0.0.0.0:8081 — chart defaults
	// this module never moves).
	HttpPort int
	GrpcPort int

	// AuthnKeysSecretSuffix names the module-owned Secret carrying
	// declared pre-shared API keys (`<metadata.name>-authn-keys`);
	// AuthnKeysSecretDataKey is the chart's keysSecret contract — it
	// reads the comma-separated key list from the data key `keys`.
	AuthnKeysSecretSuffix  string
	AuthnKeysSecretDataKey string
}{
	HelmChartName:          "openfga",
	HelmChartRepo:          "https://openfga.github.io/helm-charts",
	DefaultChartVersion:    "0.3.10",
	MaxMetadataNameLength:  55,
	HttpPort:               8080,
	GrpcPort:               8081,
	AuthnKeysSecretSuffix:  "-authn-keys",
	AuthnKeysSecretDataKey: "keys",
}
