package module

import (
	ocicontainerenginenodepoolv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocicontainerenginenodepool/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciContainerEngineNodePool *ocicontainerenginenodepoolv1alpha1.OciContainerEngineNodePool
	DisplayName                string
	FreeformTags               map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *ocicontainerenginenodepoolv1alpha1.OciContainerEngineNodePoolStackInput) *Locals {
	locals := &Locals{}
	locals.OciContainerEngineNodePool = stackInput.Target

	if stackInput.Target.Spec.Name != "" {
		locals.DisplayName = stackInput.Target.Spec.Name
	} else {
		locals.DisplayName = stackInput.Target.Metadata.Name
	}

	locals.FreeformTags = map[string]string{
		"resource":      "true",
		"resource_kind": cloudresourcekind.CloudResourceKind_OciContainerEngineNodePool.String(),
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
