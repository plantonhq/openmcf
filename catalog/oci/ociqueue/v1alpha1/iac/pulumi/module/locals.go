package module

import (
	ociqueuev1alpha1 "github.com/plantonhq/planton/catalog/oci/ociqueue/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciQueue     *ociqueuev1alpha1.OciQueue
	QueueName    string
	FreeformTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *ociqueuev1alpha1.OciQueueStackInput) *Locals {
	locals := &Locals{}
	locals.OciQueue = stackInput.Target
	locals.QueueName = stackInput.Target.Metadata.Name

	locals.FreeformTags = map[string]string{
		"resource":      "true",
		"resource_kind": cloudresourcekind.CloudResourceKind_OciQueue.String(),
		"resource_id":   stackInput.Target.Metadata.Id,
	}
	if stackInput.Target.Metadata.Org != "" {
		locals.FreeformTags["organization"] = stackInput.Target.Metadata.Org
	}
	if stackInput.Target.Metadata.Env != "" {
		locals.FreeformTags["environment"] = stackInput.Target.Metadata.Env
	}
	for k, v := range stackInput.Target.Metadata.Labels {
		locals.FreeformTags[k] = v
	}

	return locals
}
