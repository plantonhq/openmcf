package module

import (
	"strconv"

	kuberneteskarpenterec2nodeclassv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskarpenterec2nodeclass/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds the resolved inputs the module operates on: the full target
// resource plus the scalar identifiers used for the resource name, labels,
// and stack outputs. EC2NodeClass is cluster-scoped, so there is no
// namespace.
type Locals struct {
	KubernetesKarpenterEc2NodeClass *kuberneteskarpenterec2nodeclassv1alpha1.KubernetesKarpenterEc2NodeClass
	NodeClassName                   string
	Labels                          map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *kuberneteskarpenterec2nodeclassv1alpha1.KubernetesKarpenterEc2NodeClassStackInput) *Locals {
	target := stackInput.Target
	metadata := target.Metadata

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesKarpenterEc2NodeClass.String(),
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
		KubernetesKarpenterEc2NodeClass: target,
		NodeClassName:                   metadata.Name,
		Labels:                          labels,
	}
}
