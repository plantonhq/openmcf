package module

import (
	"strings"

	azureexpressroutecircuitv1alpha1 "github.com/plantonhq/planton/catalog/azure/azureexpressroutecircuit/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureExpressRouteCircuit *azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuit

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitStackInput) *Locals {
	locals := &Locals{}

	locals.AzureExpressRouteCircuit = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureExpressRouteCircuit.String()),
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

// skuTierWireValue maps the spec's tier enum onto ARM's SKU tier
// vocabulary ("Basic", "Local", "Standard", "Premium").
func skuTierWireValue(tier azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuTier) string {
	switch tier {
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuTier_BASIC:
		return "Basic"
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuTier_LOCAL:
		return "Local"
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuTier_STANDARD:
		return "Standard"
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuTier_PREMIUM:
		return "Premium"
	}
	// Unreachable: the spec's sku_tier_required contract rejects
	// unspecified before the module runs.
	return ""
}

// skuFamilyWireValue maps the spec's family enum onto ARM's SKU family
// vocabulary ("MeteredData", "UnlimitedData").
func skuFamilyWireValue(family azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuFamily) string {
	switch family {
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuFamily_METERED_DATA:
		return "MeteredData"
	case azureexpressroutecircuitv1alpha1.AzureExpressRouteCircuitSkuFamily_UNLIMITED_DATA:
		return "UnlimitedData"
	}
	// Unreachable: the spec's sku_family_required contract rejects
	// unspecified before the module runs.
	return ""
}
