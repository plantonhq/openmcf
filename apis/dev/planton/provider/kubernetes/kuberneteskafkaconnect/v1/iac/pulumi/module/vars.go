package module

var vars = struct {
	// ApiVersion of the Strimzi KafkaConnect custom resource (the `v1`
	// API — the only served API from Strimzi 1.0 onward). Must stay
	// byte-identical with the Terraform module's rendered apiVersion.
	ApiVersion string

	// UseConnectorResourcesAnnotationKey is the Strimzi annotation that
	// switches connector management to DECLARATIVE mode: with it set to
	// "true" the operator reconciles KafkaConnector resources against
	// this cluster and reverts any change made through the Connect REST
	// API. The module stamps it unconditionally — KubernetesKafkaConnector
	// resources only work against clusters carrying it.
	UseConnectorResourcesAnnotationKey string

	// MetricsConfigKey is the ConfigMap key the metricsConfig block
	// points at when spec.metrics.enabled renders the module-owned JMX
	// exporter rules ConfigMap (the upstream Strimzi example's key).
	MetricsConfigKey string
}{
	ApiVersion:                         "kafka.strimzi.io/v1",
	UseConnectorResourcesAnnotationKey: "strimzi.io/use-connector-resources",
	MetricsConfigKey:                   "metrics-config.yml",
}
