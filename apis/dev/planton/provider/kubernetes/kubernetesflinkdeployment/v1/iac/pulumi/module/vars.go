package module

var vars = struct {
	// ApiVersion / Kind of the rendered custom resource.
	ApiVersion string
	Kind       string

	// NameBudget is the ceiling on metadata.name: the operator derives
	// child names by suffixing (`<name>-rest`,
	// `<name>-taskmanager-N-M`) and Kubernetes object names cap at 63
	// characters — 45 keeps every derived name inside the budget. Both
	// engines fail loudly past it.
	NameBudget int

	// Operator naming contract: the JobManager REST Service is
	// `<name>-rest` (IngressUtils' REST_SVC_NAME_SUFFIX at the pinned
	// operator), serving the Flink REST API and web UI on RestPort.
	RestServiceSuffix string
	RestPort          int

	// DefaultServiceAccount is the account the KubernetesFlinkOperator's
	// chart creates with reconcile RBAC — the spec's own default.
	DefaultServiceAccount string

	// DefaultCpuCores / DefaultMemory are the fallback tier sizing (the
	// spec's default_container_resources values: 1 CPU / 2Gi) — the
	// operator's validator REQUIRES resource memory on both tiers, so
	// the blocks always render.
	DefaultCpuCores float64
	DefaultMemory   string

	// MainContainerName is the operator's pod-template merge contract
	// (its own examples: "Do not change the main container name").
	MainContainerName string

	// PodTemplateName is the metadata.name of the rendered pod template.
	PodTemplateName string
}{
	ApiVersion:            "flink.apache.org/v1beta1",
	Kind:                  "FlinkDeployment",
	NameBudget:            45,
	RestServiceSuffix:     "-rest",
	RestPort:              8081,
	DefaultServiceAccount: "flink",
	DefaultCpuCores:       1.0,
	DefaultMemory:         "2Gi",
	MainContainerName:     "flink-main-container",
	PodTemplateName:       "pod-template",
}
