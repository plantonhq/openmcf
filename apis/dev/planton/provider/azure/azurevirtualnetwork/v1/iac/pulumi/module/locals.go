package module

import (
	"strings"

	azurevirtualnetworkv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurevirtualnetwork/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureVirtualNetwork *azurevirtualnetworkv1.AzureVirtualNetwork

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal.
	ResourceGroupName string

	// EncryptionEnforcement is the ARM string for the spec's enum
	// ("AllowUnencrypted"/"DropUnencrypted"), or empty when unspecified so
	// no encryption block is sent (ARM's default: encryption off).
	EncryptionEnforcement string

	// PrivateEndpointVnetPolicies is the ARM string for the spec's enum
	// ("Basic"), or empty when unspecified so the provider applies ARM's
	// default ("Disabled").
	PrivateEndpointVnetPolicies string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurevirtualnetworkv1.AzureVirtualNetworkStackInput) *Locals {
	locals := &Locals{}

	locals.AzureVirtualNetwork = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	switch target.Spec.Encryption {
	case azurevirtualnetworkv1.AzureVirtualNetworkEncryptionEnforcement_ALLOW_UNENCRYPTED:
		locals.EncryptionEnforcement = "AllowUnencrypted"
	case azurevirtualnetworkv1.AzureVirtualNetworkEncryptionEnforcement_DROP_UNENCRYPTED:
		locals.EncryptionEnforcement = "DropUnencrypted"
	}

	if target.Spec.PrivateEndpointVnetPolicies == azurevirtualnetworkv1.AzureVirtualNetworkPrivateEndpointVnetPolicies_BASIC {
		locals.PrivateEndpointVnetPolicies = "Basic"
	}

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureVirtualNetwork.String()),
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
