package module

import (
	"strings"

	azurenetworksecuritygroupv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurenetworksecuritygroup/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureNetworkSecurityGroup *azurenetworksecuritygroupv1.AzureNetworkSecurityGroup

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurenetworksecuritygroupv1.AzureNetworkSecurityGroupStackInput) *Locals {
	locals := &Locals{}

	locals.AzureNetworkSecurityGroup = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureNetworkSecurityGroup.String()),
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

// directionToArm maps the spec's direction enum to ARM's
// SecurityRuleDirection string.
func directionToArm(direction azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleDirection) string {
	switch direction {
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleDirection_INBOUND:
		return "Inbound"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleDirection_OUTBOUND:
		return "Outbound"
	}
	return ""
}

// accessToArm maps the spec's access enum to ARM's SecurityRuleAccess
// string.
func accessToArm(access azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleAccess) string {
	switch access {
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleAccess_ALLOW:
		return "Allow"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleAccess_DENY:
		return "Deny"
	}
	return ""
}

// protocolToArm maps the spec's protocol enum to ARM's SecurityRuleProtocol
// string (ANY is ARM's "*").
func protocolToArm(protocol azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol) string {
	switch protocol {
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_ANY:
		return "*"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_TCP:
		return "Tcp"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_UDP:
		return "Udp"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_ICMP:
		return "Icmp"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_AH:
		return "Ah"
	case azurenetworksecuritygroupv1.AzureNetworkSecurityGroupRuleProtocol_ESP:
		return "Esp"
	}
	return ""
}
