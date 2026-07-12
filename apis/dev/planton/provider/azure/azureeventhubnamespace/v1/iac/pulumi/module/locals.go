package module

import (
	"strings"

	azureeventhubnamespacev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureeventhubnamespace/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureEventHubNamespace *azureeventhubnamespacev1.AzureEventHubNamespace
	ResourceGroupName      string
	AzureTags              map[string]string
}

// skuStrings maps the spec's SKU enum to ARM's wire values. The unspecified
// row deploys STANDARD -- the full-featured multi-tenant tier -- because an
// unmapped enum would send the empty string, which the provider rejects.
var skuStrings = map[azureeventhubnamespacev1.AzureEventHubNamespaceSku]string{
	azureeventhubnamespacev1.AzureEventHubNamespaceSku_azure_event_hub_namespace_sku_unspecified: "Standard",
	azureeventhubnamespacev1.AzureEventHubNamespaceSku_BASIC:                                     "Basic",
	azureeventhubnamespacev1.AzureEventHubNamespaceSku_STANDARD:                                  "Standard",
	azureeventhubnamespacev1.AzureEventHubNamespaceSku_PREMIUM:                                   "Premium",
}

// identityTypeStrings maps the identity-model enum to ARM's values --
// Event Hubs namespaces support all three managed-identity models.
var identityTypeStrings = map[azureeventhubnamespacev1.AzureEventHubNamespaceIdentityType]string{
	azureeventhubnamespacev1.AzureEventHubNamespaceIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azureeventhubnamespacev1.AzureEventHubNamespaceIdentityType_USER_ASSIGNED:            "UserAssigned",
	azureeventhubnamespacev1.AzureEventHubNamespaceIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// networkDefaultActionStrings maps the firewall's default-action enum to
// ARM's values. No unspecified row: the spec enum requires an explicit
// choice (Azure itself requires default_action when the rule set is
// declared).
var networkDefaultActionStrings = map[azureeventhubnamespacev1.AzureEventHubNetworkRuleSetDefaultAction]string{
	azureeventhubnamespacev1.AzureEventHubNetworkRuleSetDefaultAction_ALLOW: "Allow",
	azureeventhubnamespacev1.AzureEventHubNetworkRuleSetDefaultAction_DENY:  "Deny",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureeventhubnamespacev1.AzureEventHubNamespaceStackInput) *Locals {
	locals := &Locals{}

	locals.AzureEventHubNamespace = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureEventHubNamespace.String()),
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
	// literal "azure_event_hub_namespace" for resource_kind and fall back
	// to metadata.name for resource_id, while this module emits the lowered
	// enum string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
