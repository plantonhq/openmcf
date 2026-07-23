package module

import (
	"strconv"

	kubernetesperconamysqloperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesperconamysqloperator/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds computed values derived from the stack input for use across
// the module. Every resolution here has an exact twin in the Terraform
// module's locals.tf — keep them in lockstep.
type Locals struct {
	Spec *kubernetesperconamysqloperatorv1.KubernetesPerconaMysqlOperatorSpec

	// Resource-identity labels stamped on the module-created satellites
	// (the namespace — never injected into the chart's own resources;
	// Helm owns those).
	Labels map[string]string

	// Namespace the operator installs into (resolved literal from the
	// spec's value-or-ref). With the default watch scope this is also
	// the only namespace the operator reconciles databases in.
	Namespace string

	// Helm release name — metadata.name. The chart derives every
	// resource name from the release (Deployment, ServiceAccount, RBAC
	// through its fullname helper), so distinct release names keep
	// multiple namespace-scoped installations from colliding.
	ReleaseName string

	// Chart version resolved to the pinned default when unset, so both
	// engines install the same chart whether or not the platform's
	// defaulting middleware ran.
	ChartVersion string

	// WatchWidened is true when the spec widens the watch scope
	// (cluster-wide or a namespace fence) — the arms in which the chart
	// grants cluster-scoped RBAC and the module owns the operator's
	// cluster-scoped validation webhook (see webhook.go).
	WatchWidened bool
}

// initializeLocals extracts and transforms spec fields into module-local
// values.
func initializeLocals(_ *pulumi.Context, stackInput *kubernetesperconamysqloperatorv1.KubernetesPerconaMysqlOperatorStackInput) *Locals {
	target := stackInput.Target
	spec := target.Spec

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesPerconaMysqlOperator.String(),
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
		Spec:         spec,
		Labels:       labels,
		Namespace:    spec.Namespace.GetValue(),
		ReleaseName:  target.Metadata.Name,
		ChartVersion: chartVersion,
		WatchWidened: spec.GetWatch().GetClusterWide() || len(spec.GetWatch().GetNamespaces()) > 0,
	}
}
