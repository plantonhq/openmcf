package module

import (
	"strings"

	azurefrontdoorendpointv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefrontdoorendpoint/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFrontDoorEndpoint *azurefrontdoorendpointv1alpha1.AzureFrontDoorEndpoint
	ProfileId              string
	AzureTags              map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefrontdoorendpointv1alpha1.AzureFrontDoorEndpointStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFrontDoorEndpoint = stackInput.Target
	target := stackInput.Target

	locals.ProfileId = target.Spec.ProfileId.GetValue()

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureFrontDoorEndpoint.String()),
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

	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
