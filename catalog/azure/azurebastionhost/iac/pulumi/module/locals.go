package module

import (
	"strings"

	azurebastionhostv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurebastionhost/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureBastionHost *azurebastionhostv1alpha1.AzureBastionHost

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// SkuWireValue is the ARM wire value for the spec's SKU enum.
	// Unspecified deploys "Basic", the provider's own default, kept
	// explicit so both engines send identical wire shapes.
	SkuWireValue string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

// skuWire maps the spec's SKU enum NAME to the ARM wire value.
var skuWire = map[azurebastionhostv1alpha1.AzureBastionHostSku]string{
	azurebastionhostv1alpha1.AzureBastionHostSku_DEVELOPER: "Developer",
	azurebastionhostv1alpha1.AzureBastionHostSku_BASIC:     "Basic",
	azurebastionhostv1alpha1.AzureBastionHostSku_STANDARD:  "Standard",
	azurebastionhostv1alpha1.AzureBastionHostSku_PREMIUM:   "Premium",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurebastionhostv1alpha1.AzureBastionHostStackInput) *Locals {
	locals := &Locals{}

	locals.AzureBastionHost = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	locals.SkuWireValue = "Basic"
	if wire, ok := skuWire[target.Spec.Sku]; ok {
		locals.SkuWireValue = wire
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureBastionHost.String()),
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
