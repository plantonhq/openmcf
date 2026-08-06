package module

import (
	"strings"

	azurefirewallv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurefirewall/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureFirewall *azurefirewallv1alpha1.AzureFirewall

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurefirewallv1alpha1.AzureFirewallStackInput) *Locals {
	locals := &Locals{}

	locals.AzureFirewall = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureFirewall.String()),
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

// skuNameWireValue maps the deployment-model enum to the wire vocabulary.
// Unspecified deploys AZFW_VNet -- the standard hub-spoke shape, sent
// explicitly so both engines produce an identical payload.
func skuNameWireValue(skuName azurefirewallv1alpha1.AzureFirewallSkuName) string {
	if skuName == azurefirewallv1alpha1.AzureFirewallSkuName_AZFW_HUB {
		return "AZFW_Hub"
	}
	return "AZFW_VNet"
}

// skuTierWireValue maps the tier enum to the wire vocabulary. Unspecified
// deploys Standard -- the production default, sent explicitly so both
// engines produce an identical payload.
func skuTierWireValue(skuTier azurefirewallv1alpha1.AzureFirewallSkuTier) string {
	switch skuTier {
	case azurefirewallv1alpha1.AzureFirewallSkuTier_BASIC:
		return "Basic"
	case azurefirewallv1alpha1.AzureFirewallSkuTier_PREMIUM:
		return "Premium"
	default:
		return "Standard"
	}
}

// threatIntelModeWireValue maps the threat-intelligence mode to the wire
// vocabulary. Returns "" for unspecified so callers omit the field --
// the ARM field is server-defaulted (Alert) and the provider treats it
// as Computed, so omission lets Azure own the default.
func threatIntelModeWireValue(mode azurefirewallv1alpha1.AzureFirewallThreatIntelMode) string {
	switch mode {
	case azurefirewallv1alpha1.AzureFirewallThreatIntelMode_ALERT:
		return "Alert"
	case azurefirewallv1alpha1.AzureFirewallThreatIntelMode_DENY:
		return "Deny"
	case azurefirewallv1alpha1.AzureFirewallThreatIntelMode_OFF:
		return "Off"
	default:
		return ""
	}
}
