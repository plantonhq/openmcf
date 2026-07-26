package module

var vars = struct {
	// AmqpPort / AmqpTlsPort are RabbitMQ's fixed client listener ports
	// (5672 plain, 5671 TLS). The operator's generated Services expose
	// them and the spec does not remap them.
	AmqpPort    int
	AmqpTlsPort int

	// ManagementPort / ManagementTlsPort are the management API / UI
	// ports (15672 plain, 15671 TLS). The management plugin is one of
	// the operator's always-on essentials.
	ManagementPort    int
	ManagementTlsPort int

	// DefaultUserSecretSuffix builds the operator-generated credentials
	// Secret name `<metadata.name><suffix>` (operator naming contract:
	// ChildResourceName("default-user")).
	DefaultUserSecretSuffix string

	// HeadlessServiceSuffix builds the inter-node Service name
	// `<metadata.name><suffix>` (operator naming contract:
	// ChildResourceName("nodes")).
	HeadlessServiceSuffix string

	// NodeAffinityKubernetesHostnameKey is the topology key the
	// spread_across_nodes anti-affinity pins on.
	NodeAffinityKubernetesHostnameKey string

	// AppLabelKey identifies the cluster's own pods for the
	// spread_across_nodes anti-affinity: the operator stamps every pod
	// with `app.kubernetes.io/name: <cluster-name>`.
	AppLabelKey string
}{
	AmqpPort:                          5672,
	AmqpTlsPort:                       5671,
	ManagementPort:                    15672,
	ManagementTlsPort:                 15671,
	DefaultUserSecretSuffix:           "-default-user",
	HeadlessServiceSuffix:             "-nodes",
	NodeAffinityKubernetesHostnameKey: "kubernetes.io/hostname",
	AppLabelKey:                       "app.kubernetes.io/name",
}
