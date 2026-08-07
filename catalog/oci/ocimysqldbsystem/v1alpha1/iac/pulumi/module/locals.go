package module

import (
	ocimysqldbsystemv1alpha1 "github.com/plantonhq/planton/catalog/oci/ocimysqldbsystem/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciMysqlDbSystem *ocimysqldbsystemv1alpha1.OciMysqlDbSystem
	DisplayName      string
	FreeformTags     map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *ocimysqldbsystemv1alpha1.OciMysqlDbSystemStackInput) *Locals {
	locals := &Locals{}
	locals.OciMysqlDbSystem = stackInput.Target

	if stackInput.Target.Spec.DisplayName != "" {
		locals.DisplayName = stackInput.Target.Spec.DisplayName
	} else {
		locals.DisplayName = stackInput.Target.Metadata.Name
	}

	locals.FreeformTags = map[string]string{
		"resource":      "true",
		"resource_kind": cloudresourcekind.CloudResourceKind_OciMysqlDbSystem.String(),
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
