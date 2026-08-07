package module

import (
	"strings"

	azureappinsightswebtestv1 "github.com/plantonhq/planton/catalog/azure/azureapplicationinsightsstandardwebtest/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureApplicationInsightsStandardWebTest *azureappinsightswebtestv1.AzureApplicationInsightsStandardWebTest

	ResourceGroupName     string
	ApplicationInsightsId string

	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureappinsightswebtestv1.AzureApplicationInsightsStandardWebTestStackInput) *Locals {
	locals := &Locals{}

	locals.AzureApplicationInsightsStandardWebTest = stackInput.Target
	target := stackInput.Target
	spec := target.Spec

	locals.ResourceGroupName = spec.ResourceGroup.GetValue()
	locals.ApplicationInsightsId = spec.ApplicationInsightsId.GetValue()

	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module emits the family-wide snake-case literal
	// and falls back to metadata.name. Output-neutral (tags never feed stack
	// outputs); aligning the two shapes is a family-wide convention change.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureApplicationInsightsStandardWebTest.String()),
	}
	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}
	for k, v := range spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
