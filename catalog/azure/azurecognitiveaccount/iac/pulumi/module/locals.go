package module

import (
	"strings"

	azurecognitiveaccountv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurecognitiveaccount/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureCognitiveAccount *azurecognitiveaccountv1alpha1.AzureCognitiveAccount

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

// identityTypeWire maps the spec's identity flavors to the provider's
// comma-joined wire values.
var identityTypeWire = map[azurecognitiveaccountv1alpha1.AzureCognitiveAccountIdentityType]string{
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// networkAclsBypassWire maps the spec's bypass enum to the provider's wire
// values. Unspecified is absent -- the property is omitted so ARM applies
// its default.
var networkAclsBypassWire = map[azurecognitiveaccountv1alpha1.AzureCognitiveAccountNetworkAclsBypass]string{
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountNetworkAclsBypass_AZURE_SERVICES: "AzureServices",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountNetworkAclsBypass_NONE:           "None",
}

// raiContentLevelWire maps the spec's content-severity enum to the
// provider's wire values.
var raiContentLevelWire = map[azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyContentLevel]string{
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyContentLevel_LOW:    "Low",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyContentLevel_MEDIUM: "Medium",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyContentLevel_HIGH:   "High",
}

// raiPolicyModeWire maps the spec's responsible-AI policy modes to the
// provider's wire values.
var raiPolicyModeWire = map[azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyMode]string{
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyMode_DEFAULT:             "Default",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyMode_BLOCKING:            "Blocking",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyMode_ASYNCHRONOUS_FILTER: "AsynchronousFilter",
	azurecognitiveaccountv1alpha1.AzureCognitiveAccountRaiPolicyMode_DEFERRED:            "Deferred",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurecognitiveaccountv1alpha1.AzureCognitiveAccountStackInput) *Locals {
	locals := &Locals{}

	locals.AzureCognitiveAccount = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureCognitiveAccount.String()),
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

	for k, v := range target.Spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
