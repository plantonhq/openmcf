package module

import (
	kubernetesbackendtlspolicyv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesbackendtlspolicy/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

// Locals holds the resolved inputs the module operates on: the full target
// resource plus the scalar identifiers used for the resource name, namespace,
// labels, and stack outputs.
type Locals struct {
	KubernetesBackendTlsPolicy *kubernetesbackendtlspolicyv1.KubernetesBackendTlsPolicy
	PolicyName                 string
	Namespace                  string
	Labels                     map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *kubernetesbackendtlspolicyv1.KubernetesBackendTlsPolicyStackInput) *Locals {
	target := stackInput.Target
	metadata := target.Metadata
	spec := target.Spec

	// namespace is a StringValueOrRef foreign key. The platform middleware
	// resolves valueFrom references to literal strings before the IaC module
	// runs, so GetValue() returns the resolved value.
	namespace := spec.GetNamespace().GetValue()

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesBackendTlsPolicy.String(),
	}
	if metadata.Id != "" {
		labels[kuberneteslabelkeys.ResourceId] = metadata.Id
	}
	if metadata.Org != "" {
		labels[kuberneteslabelkeys.Organization] = metadata.Org
	}
	if metadata.Env != "" {
		labels[kuberneteslabelkeys.Environment] = metadata.Env
	}

	return &Locals{
		KubernetesBackendTlsPolicy: target,
		PolicyName:                 metadata.Name,
		Namespace:                  namespace,
		Labels:                     labels,
	}
}
