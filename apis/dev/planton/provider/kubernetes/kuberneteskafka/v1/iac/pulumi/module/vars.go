package module

var vars = struct {
	// ApiVersion of the Strimzi Kafka and KafkaNodePool custom resources
	// (the `v1` API — the only served API from Strimzi 1.0 onward). Must
	// stay byte-identical with the Terraform module's rendered
	// apiVersion.
	ApiVersion string

	// ClusterLabelKey is the Strimzi label that binds KafkaNodePool
	// (and KafkaTopic/KafkaUser) resources to their Kafka cluster.
	ClusterLabelKey string

	// MetricsConfigKey is the ConfigMap key the metricsConfig block
	// points at when spec.metrics.enabled renders the module-owned JMX
	// exporter rules ConfigMap.
	MetricsConfigKey string
}{
	ApiVersion:       "kafka.strimzi.io/v1",
	ClusterLabelKey:  "strimzi.io/cluster",
	MetricsConfigKey: "kafka-metrics-config.yml",
}
