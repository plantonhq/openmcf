package module

import (
	"strings"

	azureservicebusnamespacev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureservicebusnamespace/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureServiceBusNamespace *azureservicebusnamespacev1alpha1.AzureServiceBusNamespace
	ResourceGroupName        string
	AzureTags                map[string]string
}

// skuStrings maps the spec's SKU enum to ARM's wire values. The unspecified
// row deploys STANDARD -- the full-featured multi-tenant tier -- because an
// unmapped enum would send the empty string, which the provider rejects.
var skuStrings = map[azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceSku]string{
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceSku_azure_service_bus_namespace_sku_unspecified: "Standard",
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceSku_BASIC:                                       "Basic",
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceSku_STANDARD:                                    "Standard",
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceSku_PREMIUM:                                     "Premium",
}

// identityTypeStrings maps the identity-model enum to ARM's values --
// Service Bus namespaces support all three managed-identity models.
var identityTypeStrings = map[azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceIdentityType]string{
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceIdentityType_USER_ASSIGNED:            "UserAssigned",
	azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// networkDefaultActionStrings maps the firewall's default-action enum to
// ARM's values. The unspecified row keeps Azure's open default (Allow) --
// the block may be declared just for trusted-services or public-access
// dials.
var networkDefaultActionStrings = map[azureservicebusnamespacev1alpha1.AzureServiceBusNetworkDefaultAction]string{
	azureservicebusnamespacev1alpha1.AzureServiceBusNetworkDefaultAction_azure_service_bus_network_default_action_unspecified: "Allow",
	azureservicebusnamespacev1alpha1.AzureServiceBusNetworkDefaultAction_ALLOW:                                                "Allow",
	azureservicebusnamespacev1alpha1.AzureServiceBusNetworkDefaultAction_DENY:                                                 "Deny",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureservicebusnamespacev1alpha1.AzureServiceBusNamespaceStackInput) *Locals {
	locals := &Locals{}

	locals.AzureServiceBusNamespace = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureServiceBusNamespace.String()),
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
	// literal "azure_service_bus_namespace" for resource_kind and fall back
	// to metadata.name for resource_id, while this module emits the lowered
	// enum string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
