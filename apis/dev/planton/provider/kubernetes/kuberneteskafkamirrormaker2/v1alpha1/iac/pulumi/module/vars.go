package module

var vars = struct {
	// ApiVersion of the Strimzi KafkaMirrorMaker2 custom resource (the
	// `v1` API — the only served API from Strimzi 1.0 onward). Must stay
	// byte-identical with the Terraform module's rendered apiVersion.
	ApiVersion string

	// MetricsConfigKey is the ConfigMap key the metricsConfig block
	// points at when spec.metrics.enabled renders the module-owned JMX
	// exporter rules ConfigMap (the upstream connect-metrics example's
	// key).
	MetricsConfigKey string
}{
	ApiVersion:       "kafka.strimzi.io/v1",
	MetricsConfigKey: "metrics-config.yml",
}
