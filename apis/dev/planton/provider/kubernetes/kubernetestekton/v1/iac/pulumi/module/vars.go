package module

var vars = struct {
	// ApiVersion / Kind of the rendered custom resource.
	ApiVersion string
	Kind       string

	// TektonConfigName is the FIXED name of the rendered TektonConfig.
	// The operator's admission webhook allows exactly one TektonConfig
	// per cluster and requires this name — metadata.name of the Planton
	// resource never reaches the CR (it keys the state identity only).
	TektonConfigName string

	// DefaultTargetNamespace is the upstream default installation
	// namespace for the Tekton components.
	DefaultTargetNamespace string

	// DashboardServiceName / DashboardPort are the dashboard's fixed
	// Service handles in the target namespace (installed on profile
	// `all`).
	DashboardServiceName string
	DashboardPort        int
}{
	ApiVersion:             "operator.tekton.dev/v1alpha1",
	Kind:                   "TektonConfig",
	TektonConfigName:       "config",
	DefaultTargetNamespace: "tekton-pipelines",
	DashboardServiceName:   "tekton-dashboard",
	DashboardPort:          9097,
}
