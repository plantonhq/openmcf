package module

import (
	"strings"

	azureprivatednsresolverforwardingrulesetv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureprivatednsresolverforwardingruleset/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzurePrivateDnsResolverForwardingRuleset *azureprivatednsresolverforwardingrulesetv1alpha1.AzurePrivateDnsResolverForwardingRuleset

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// OutboundEndpointIds are the resolved literal ARM ids of the outbound
	// endpoints the ruleset binds (references pre-resolved by the platform
	// middleware).
	OutboundEndpointIds []string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order. Individual forwarding rules carry NO
	// tags (ARM gives them a metadata map instead) -- these tags land on
	// the ruleset only.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureprivatednsresolverforwardingrulesetv1alpha1.AzurePrivateDnsResolverForwardingRulesetStackInput) *Locals {
	locals := &Locals{}

	locals.AzurePrivateDnsResolverForwardingRuleset = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	for _, endpointId := range target.Spec.OutboundEndpointIds {
		locals.OutboundEndpointIds = append(locals.OutboundEndpointIds, endpointId.GetValue())
	}

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzurePrivateDnsResolverForwardingRuleset.String()),
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
