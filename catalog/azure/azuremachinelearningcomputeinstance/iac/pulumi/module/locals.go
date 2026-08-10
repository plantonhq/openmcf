package module

import (
	"strings"

	azuremachinelearningcomputeinstancev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuremachinelearningcomputeinstance/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMachineLearningComputeInstance *azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstance

	// WorkspaceId is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal ARM ID.
	WorkspaceId string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order. NOTE: instance tags are ForceNew on
	// the provider -- changing any tag replaces the instance (its OS disk
	// and local files go with it).
	AzureTags map[string]string
}

// identityTypeWire maps the spec's identity flavors to the provider's
// comma-joined wire values.
var identityTypeWire = map[azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstanceIdentityType]string{
	azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstanceIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstanceIdentityType_USER_ASSIGNED:            "UserAssigned",
	azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstanceIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremachinelearningcomputeinstancev1alpha1.AzureMachineLearningComputeInstanceStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMachineLearningComputeInstance = stackInput.Target
	target := stackInput.Target

	locals.WorkspaceId = target.Spec.WorkspaceId.GetValue()

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMachineLearningComputeInstance.String()),
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
