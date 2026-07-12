package module

import (
	"strings"

	azureprivateendpointv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureprivateendpoint/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzurePrivateEndpoint *azureprivateendpointv1.AzurePrivateEndpoint
	ResourceGroupName    string
	AzureTags            map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureprivateendpointv1.AzurePrivateEndpointStackInput) *Locals {
	locals := &Locals{}

	locals.AzurePrivateEndpoint = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Metadata-derived tags first; the user's spec tags merge over them in
	// main.go so an org's governance conventions win on key collision.
	//
	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module emits the family-wide snake-case literal
	// and falls back to metadata.name. Output-neutral (tags never feed stack
	// outputs); aligning the two shapes is a family-wide convention change.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzurePrivateEndpoint.String()),
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
