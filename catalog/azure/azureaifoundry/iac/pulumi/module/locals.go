package module

import (
	"strings"

	azureaifoundryv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureaifoundry/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureAiFoundry *azureaifoundryv1alpha1.AzureAiFoundry

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
var identityTypeWire = map[azureaifoundryv1alpha1.AzureAiFoundryIdentityType]string{
	azureaifoundryv1alpha1.AzureAiFoundryIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azureaifoundryv1alpha1.AzureAiFoundryIdentityType_USER_ASSIGNED:            "UserAssigned",
	azureaifoundryv1alpha1.AzureAiFoundryIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// isolationModeWire maps the spec's isolation modes to the provider's
// wire values. Unspecified is absent -- the property is omitted and the
// value is read back.
var isolationModeWire = map[azureaifoundryv1alpha1.AzureAiFoundryIsolationMode]string{
	azureaifoundryv1alpha1.AzureAiFoundryIsolationMode_DISABLED:                     "Disabled",
	azureaifoundryv1alpha1.AzureAiFoundryIsolationMode_ALLOW_INTERNET_OUTBOUND:      "AllowInternetOutbound",
	azureaifoundryv1alpha1.AzureAiFoundryIsolationMode_ALLOW_ONLY_APPROVED_OUTBOUND: "AllowOnlyApprovedOutbound",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureaifoundryv1alpha1.AzureAiFoundryStackInput) *Locals {
	locals := &Locals{}

	locals.AzureAiFoundry = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureAiFoundry.String()),
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
