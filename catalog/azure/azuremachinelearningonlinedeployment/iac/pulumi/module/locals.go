package module

import (
	"fmt"
	"strings"

	azuremachinelearningonlinedeploymentv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuremachinelearningonlinedeployment/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMachineLearningOnlineDeployment *azuremachinelearningonlinedeploymentv1alpha1.AzureMachineLearningOnlineDeployment

	// EndpointId is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal ARM ID.
	EndpointId string

	// ResourceGroupName / WorkspaceName / EndpointName are parsed from
	// EndpointId: azure-native addresses ARM children by their ancestor
	// NAMES where the raw ARM layer takes the parent's full ID -- the two
	// engines consume the same spec reference either way.
	ResourceGroupName string
	WorkspaceName     string
	EndpointName      string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

// parseEndpointId splits an online-endpoint ARM ID into its
// resource-group, workspace and endpoint names. The ID shape is fixed
// by ARM: /subscriptions/{sub}/resourceGroups/{rg}/providers/
// Microsoft.MachineLearningServices/workspaces/{ws}/onlineEndpoints/{name}
func parseEndpointId(endpointId string) (resourceGroupName, workspaceName, endpointName string, err error) {
	segments := strings.Split(strings.Trim(endpointId, "/"), "/")
	for i := 0; i+1 < len(segments); i += 2 {
		switch strings.ToLower(segments[i]) {
		case "resourcegroups":
			resourceGroupName = segments[i+1]
		case "workspaces":
			workspaceName = segments[i+1]
		case "onlineendpoints":
			endpointName = segments[i+1]
		}
	}
	if resourceGroupName == "" || workspaceName == "" || endpointName == "" {
		return "", "", "", fmt.Errorf("endpoint id %q does not carry resourceGroups, workspaces and onlineEndpoints segments", endpointId)
	}
	return resourceGroupName, workspaceName, endpointName, nil
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremachinelearningonlinedeploymentv1alpha1.AzureMachineLearningOnlineDeploymentStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMachineLearningOnlineDeployment = stackInput.Target
	target := stackInput.Target

	locals.EndpointId = target.Spec.EndpointId.GetValue()

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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMachineLearningOnlineDeployment.String()),
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
