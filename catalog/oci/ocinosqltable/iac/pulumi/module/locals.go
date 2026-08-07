package module

import (
	ocinosqltablev1alpha1 "github.com/plantonhq/planton/catalog/oci/ocinosqltable/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	OciNosqlTable *ocinosqltablev1alpha1.OciNosqlTable
	TableName     string
	FreeformTags  map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *ocinosqltablev1alpha1.OciNosqlTableStackInput) *Locals {
	locals := &Locals{}
	locals.OciNosqlTable = stackInput.Target

	locals.TableName = stackInput.Target.Spec.Name

	locals.FreeformTags = map[string]string{
		"resource":      "true",
		"resource_kind": cloudresourcekind.CloudResourceKind_OciNosqlTable.String(),
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
