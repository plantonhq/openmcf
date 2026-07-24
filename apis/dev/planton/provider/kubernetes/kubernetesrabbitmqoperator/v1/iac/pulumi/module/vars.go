package module

var vars = struct {
	// OperatorRelease is the pinned rabbitmq/cluster-operator release tag.
	//
	// MUST stay in sync with `rabbitmq_cluster_operator_release` in
	// pkg/kubernetes/kubernetestypes/Makefile and the Terraform module's
	// operator_release local, so the installed CRD schema matches the
	// crd2pulumi-generated typed SDK that KubernetesRabbitMq is built
	// against. Always an exact release TAG, never a branch: a branch ref
	// moves as patches land, so the same deployed resource would install
	// different operators at different times — tag pinning keeps installs
	// reproducible.
	OperatorRelease string

	// Namespace is the fixed installation namespace. The release
	// manifest bakes `rabbitmq-system` into its own cross-references
	// (webhook client configuration, Certificate DNS names, CA-injection
	// annotations, cluster-role binding subjects) — it is not
	// configurable, and exactly one install per cluster is the upstream
	// contract (the admission webhooks are cluster-scoped singletons).
	Namespace string

	// DeploymentName / MetricsServiceName / CrdName are the release
	// manifest's fixed object names the outputs point at.
	DeploymentName     string
	MetricsServiceName string
	CrdName            string

	// MetricsPort is the operator's Prometheus metrics port (the metrics
	// Service's port in the release manifest).
	MetricsPort int
}{
	OperatorRelease:    "v2.22.3",
	Namespace:          "rabbitmq-system",
	DeploymentName:     "rabbitmq-cluster-operator",
	MetricsServiceName: "rabbitmq-cluster-operator-metrics-service",
	CrdName:            "rabbitmqclusters.rabbitmq.com",
	MetricsPort:        8080,
}

// ManifestURL is the released single-file manifest for the pinned tag —
// the operator's OFFICIAL distribution (it has no Helm chart). The release
// pipeline pins the operator image tag inside this asset (the in-repo
// kustomization's `latest` is build-time only; verified at the pin).
func ManifestURL() string {
	return "https://github.com/rabbitmq/cluster-operator/releases/download/" +
		vars.OperatorRelease + "/cluster-operator.yml"
}
