package module

import (
	"strings"

	azureuserassignedidentityv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureuserassignedidentity/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureUserAssignedIdentity *azureuserassignedidentityv1.AzureUserAssignedIdentity

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal.
	ResourceGroupName string

	// IsolationScope is the ARM string for the spec's enum ("Regional"), or
	// empty when unspecified so the provider applies ARM's default (no
	// isolation).
	IsolationScope string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureuserassignedidentityv1.AzureUserAssignedIdentityStackInput) *Locals {
	locals := &Locals{}

	locals.AzureUserAssignedIdentity = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	if target.Spec.IsolationScope == azureuserassignedidentityv1.AzureUserAssignedIdentityIsolationScope_REGIONAL {
		locals.IsolationScope = "Regional"
	}

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureUserAssignedIdentity.String()),
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
