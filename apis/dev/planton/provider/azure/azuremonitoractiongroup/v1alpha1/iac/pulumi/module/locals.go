package module

import (
	"strings"

	azuremonitoractiongroupv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremonitoractiongroup/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMonitorActionGroup *azuremonitoractiongroupv1alpha1.AzureMonitorActionGroup
	ResourceGroupName       string
	AzureTags               map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremonitoractiongroupv1alpha1.AzureMonitorActionGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMonitorActionGroup = stackInput.Target
	target := stackInput.Target

	// The resource_group field is a StringValueOrRef. The platform middleware
	// resolves valueFrom references before IaC modules run, so .GetValue()
	// always returns the resolved literal string.
	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Identity tags derived from metadata; user tags merge OVER these (the
	// governance surface belongs to the user).
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMonitorActionGroup.String()),
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

	// PARITY-EXCEPTION: the Terraform module's base tags use the snake_case
	// literal "azure_monitor_action_group" for resource_kind and fall back to
	// metadata.name for resource_id, while this module emits the lowered enum
	// string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
