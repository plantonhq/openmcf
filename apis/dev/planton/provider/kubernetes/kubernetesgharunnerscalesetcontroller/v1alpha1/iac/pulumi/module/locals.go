package module

import (
	"strconv"

	kubernetesgharunnerscalesetcontrollerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesgharunnerscalesetcontroller/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds computed values derived from the stack input for use across
// the module. Every resolution here has an exact twin in the Terraform
// module's locals.tf — keep them in lockstep.
type Locals struct {
	Spec *kubernetesgharunnerscalesetcontrollerv1alpha1.KubernetesGhaRunnerScaleSetControllerSpec

	// Resource-identity labels stamped on the module-created satellites
	// (the namespace — never injected into the chart's own resources;
	// Helm owns those).
	Labels map[string]string

	// Namespace the controller installs into (resolved literal from the
	// spec's value-or-ref).
	Namespace string

	// Helm release name — metadata.name. Several controllers can coexist
	// only when each is fenced to its own namespace
	// (flags.watch_single_namespace); one cluster-wide controller is the
	// sane default.
	ReleaseName string

	// Chart version resolved to the pinned default when unset, so both
	// engines install the same chart whether or not the platform's
	// defaulting middleware ran.
	ChartVersion string

	// ServiceAccountName is the controller's ServiceAccount —
	// `fullnameOverride` pins the chart fullname to the release name, and
	// the chart names the created ServiceAccount exactly the fullname, so
	// this equals metadata.name. Exported: runner scale sets reference it
	// when this controller watches a single namespace.
	ServiceAccountName string
}

// initializeLocals extracts and transforms spec fields into module-local
// values.
func initializeLocals(_ *pulumi.Context, stackInput *kubernetesgharunnerscalesetcontrollerv1alpha1.KubernetesGhaRunnerScaleSetControllerStackInput) *Locals {
	target := stackInput.Target
	spec := target.Spec

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesGhaRunnerScaleSetController.String(),
	}
	if target.Metadata.Id != "" {
		labels[kuberneteslabelkeys.ResourceId] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		labels[kuberneteslabelkeys.Organization] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		labels[kuberneteslabelkeys.Environment] = target.Metadata.Env
	}

	chartVersion := spec.GetChartVersion()
	if chartVersion == "" {
		chartVersion = vars.DefaultChartVersion
	}

	return &Locals{
		Spec:               spec,
		Labels:             labels,
		Namespace:          spec.Namespace.GetValue(),
		ReleaseName:        target.Metadata.Name,
		ChartVersion:       chartVersion,
		ServiceAccountName: target.Metadata.Name,
	}
}
