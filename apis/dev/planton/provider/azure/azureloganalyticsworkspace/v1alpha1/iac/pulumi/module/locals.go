package module

import (
	"strings"

	azureloganalyticsworkspacev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureloganalyticsworkspace/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureLogAnalyticsWorkspace *azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspace
	ResourceGroupName          string
	AzureTags                  map[string]string
}

// skuStrings maps the spec's SKU enum to ARM's wire values. The unspecified
// row deploys Azure's recommended PerGB2018 pay-as-you-go tier -- an
// unmapped enum would send the empty string, which the provider rejects.
// Standard/Premium/LACluster/Unlimited are deliberately not mapped: Azure
// blocks creating workspaces on them (see the spec enum comment).
var skuStrings = map[azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku]string{
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku_azure_log_analytics_workspace_sku_unspecified: "PerGB2018",
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku_PER_GB_2018:                                   "PerGB2018",
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku_CAPACITY_RESERVATION:                          "CapacityReservation",
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku_PER_NODE:                                      "PerNode",
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceSku_STANDALONE:                                    "Standalone",
}

// identityTypeStrings maps the identity-model enum to ARM's values.
// Workspaces accept exactly SystemAssigned or UserAssigned -- the combined
// model does not exist on this resource.
var identityTypeStrings = map[azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceIdentityType]string{
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceIdentityType_SYSTEM_ASSIGNED: "SystemAssigned",
	azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceIdentityType_USER_ASSIGNED:   "UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureloganalyticsworkspacev1alpha1.AzureLogAnalyticsWorkspaceStackInput) *Locals {
	locals := &Locals{}

	locals.AzureLogAnalyticsWorkspace = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureLogAnalyticsWorkspace.String()),
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
	// literal "azure_log_analytics_workspace" for resource_kind and fall back
	// to metadata.name for resource_id, while this module emits the lowered
	// enum string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
