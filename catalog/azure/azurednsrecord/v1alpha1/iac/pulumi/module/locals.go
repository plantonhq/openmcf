package module

import (
	"strings"

	azurednsrecordv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurednsrecord/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureDnsRecord    *azurednsrecordv1alpha1.AzureDnsRecord
	ResourceGroupName string
	ZoneName          string
	AzureTags         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurednsrecordv1alpha1.AzureDnsRecordStackInput) *Locals {
	locals := &Locals{}

	locals.AzureDnsRecord = stackInput.Target

	target := stackInput.Target

	// resource_group and zone_name are StringValueOrRef fields. The
	// platform middleware resolves valueFrom references before IaC modules
	// run, so .GetValue() always returns the resolved literal string.
	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()
	locals.ZoneName = target.Spec.ZoneName.GetValue()

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide. Record-set tags land in ARM's record-set metadata.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureDnsRecord.String()),
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
