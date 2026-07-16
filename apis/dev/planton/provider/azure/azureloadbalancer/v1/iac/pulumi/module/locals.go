package module

import (
	"strings"

	azureloadbalancerv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureloadbalancer/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureLoadBalancer *azureloadbalancerv1.AzureLoadBalancer
	ResourceGroupName string
	AzureTags         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureloadbalancerv1.AzureLoadBalancerStackInput) *Locals {
	locals := &Locals{}

	locals.AzureLoadBalancer = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module hardcodes the family-wide snake-case
	// literal and falls back to metadata.name. Output-neutral (tags never
	// feed stack outputs); aligning the two shapes is a family-wide
	// convention change, not a per-kind fix.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureLoadBalancer.String()),
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

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
