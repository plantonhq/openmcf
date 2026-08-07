package module

import (
	ociloggroupv1alpha1 "github.com/plantonhq/planton/catalog/oci/ociloggroup/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciLogGroup  *ociloggroupv1alpha1.OciLogGroup
	GroupName    string
	FreeformTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *ociloggroupv1alpha1.OciLogGroupStackInput) *Locals {
	locals := &Locals{}
	locals.OciLogGroup = stackInput.Target
	locals.GroupName = stackInput.Target.Metadata.Name

	locals.FreeformTags = map[string]string{
		"resource":      "true",
		"resource_kind": cloudresourcekind.CloudResourceKind_OciLogGroup.String(),
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
