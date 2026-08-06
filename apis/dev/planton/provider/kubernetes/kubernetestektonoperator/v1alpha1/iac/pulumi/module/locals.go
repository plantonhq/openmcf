package module

import (
	"strconv"

	kubernetestektonoperatorv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetestektonoperator/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds computed values derived from the stack input for use across
// the module. Every resolution here has an exact twin in the Terraform
// module's locals.tf — keep them in lockstep.
type Locals struct {
	KubernetesTektonOperator *kubernetestektonoperatorv1alpha1.KubernetesTektonOperator
	Spec                     *kubernetestektonoperatorv1alpha1.KubernetesTektonOperatorSpec

	// ResourceName keys the applied manifest bundle in the Pulumi state.
	ResourceName string

	// Labels tie the install back to the Planton resource. The manifest's
	// own documents keep their upstream labels untouched (faithful
	// apply); these are used on the module's own state identity only.
	Labels map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *kubernetestektonoperatorv1alpha1.KubernetesTektonOperatorStackInput) *Locals {
	target := stackInput.Target

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesTektonOperator.String(),
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

	return &Locals{
		KubernetesTektonOperator: target,
		Spec:                     target.Spec,
		ResourceName:             target.Metadata.Name,
		Labels:                   labels,
	}
}
