package module

import (
	"strings"

	azurevirtualmachinescalesetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurevirtualmachinescaleset/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureVirtualMachineScaleSet *azurevirtualmachinescalesetv1.AzureVirtualMachineScaleSet
	ResourceGroupName           string
	AzureTags                   map[string]string
	IsUniform                   bool
	IsLinux                     bool
	ComputerNamePrefix          string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurevirtualmachinescalesetv1.AzureVirtualMachineScaleSetStackInput) *Locals {
	locals := &Locals{}

	locals.AzureVirtualMachineScaleSet = stackInput.Target
	target := stackInput.Target
	spec := target.Spec

	locals.ResourceGroupName = spec.ResourceGroup.GetValue()

	// The dispatch axes: ONE proto surface realizes onto azurerm's three
	// scale-set resources. Unset orchestration_mode applies FLEXIBLE
	// (Azure's recommendation for new workloads).
	locals.IsUniform = spec.OrchestrationMode == azurevirtualmachinescalesetv1.AzureVirtualMachineScaleSetOrchestrationMode_UNIFORM
	locals.IsLinux = spec.OsProfile.GetLinux() != nil

	// Instance computer names derive from this prefix (Azure appends a
	// unique suffix); unset defaults to the scale-set name.
	locals.ComputerNamePrefix = spec.OsProfile.ComputerNamePrefix

	// PARITY-EXCEPTION: resource_kind here is the lowered
	// CloudResourceKind enum string and resource_id is omitted when
	// metadata.id is empty, while the Terraform module hardcodes the
	// family-wide snake-case literal and falls back to metadata.name.
	// Output-neutral (tags never feed stack outputs); aligning the two
	// shapes is a family-wide convention change, not a per-kind fix.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureVirtualMachineScaleSet.String()),
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
	for key, value := range spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
