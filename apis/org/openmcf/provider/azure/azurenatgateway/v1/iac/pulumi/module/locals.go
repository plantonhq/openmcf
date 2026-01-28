package module

import (
	"fmt"
	"strings"

	azurenatgatewayv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/azure/azurenatgateway/v1"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureNatGateway *azurenatgatewayv1.AzureNatGateway
	NatGatewayName  string
	SubnetId        string
	ResourceGroup   string
	Location        string
	AzureTags       map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurenatgatewayv1.AzureNatGatewayStackInput) *Locals {
	locals := &Locals{}

	locals.AzureNatGateway = stackInput.Target
	target := stackInput.Target

	// Generate NAT Gateway name from metadata
	locals.NatGatewayName = fmt.Sprintf("natgw-%s", target.Metadata.Name)

	// Get subnet ID (either direct value or from reference)
	locals.SubnetId = target.Spec.SubnetId.GetValue()

	// Parse subnet ID to extract resource group and location
	// Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vnet}/subnets/{subnet}
	parts := strings.Split(locals.SubnetId, "/")
	if len(parts) >= 5 {
		for i, part := range parts {
			if part == "resourceGroups" && i+1 < len(parts) {
				locals.ResourceGroup = parts[i+1]
			}
		}
	}

	// For location, we'll need to query the resource group or use a default
	// For simplicity, we'll extract from subnet or use a sensible default
	locals.Location = "eastus" // This should ideally be extracted or passed

	// Create Azure tags for resource tagging
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureNatGateway.String()),
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

	// Merge user-provided tags
	for k, v := range target.Spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}

// getIdleTimeoutMinutes returns the idle timeout or default of 4
func getIdleTimeoutMinutes(spec *azurenatgatewayv1.AzureNatGatewaySpec) int {
	if spec.IdleTimeoutMinutes != nil {
		return int(*spec.IdleTimeoutMinutes)
	}
	return 4 // Azure default
}
