package module

import (
	"strings"

	azureaifoundryprojectv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureaifoundryproject/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureAiFoundryProject *azureaifoundryprojectv1alpha1.AzureAiFoundryProject

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

// identityTypeWire maps the spec's identity flavors to the provider's
// comma-joined wire values.
var identityTypeWire = map[azureaifoundryprojectv1alpha1.AzureAiFoundryProjectIdentityType]string{
	azureaifoundryprojectv1alpha1.AzureAiFoundryProjectIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azureaifoundryprojectv1alpha1.AzureAiFoundryProjectIdentityType_USER_ASSIGNED:            "UserAssigned",
	azureaifoundryprojectv1alpha1.AzureAiFoundryProjectIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureaifoundryprojectv1alpha1.AzureAiFoundryProjectStackInput) *Locals {
	locals := &Locals{}

	locals.AzureAiFoundryProject = stackInput.Target
	target := stackInput.Target

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureAiFoundryProject.String()),
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
