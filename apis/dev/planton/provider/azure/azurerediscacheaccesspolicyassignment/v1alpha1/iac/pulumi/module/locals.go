package module

import (
	azurerediscacheaccesspolicyassignmentv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurerediscacheaccesspolicyassignment/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureRedisCacheAccessPolicyAssignment *azurerediscacheaccesspolicyassignmentv1alpha1.AzureRedisCacheAccessPolicyAssignment
	RedisCacheId                          string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurerediscacheaccesspolicyassignmentv1alpha1.AzureRedisCacheAccessPolicyAssignmentStackInput) *Locals {
	locals := &Locals{}

	locals.AzureRedisCacheAccessPolicyAssignment = stackInput.Target
	locals.RedisCacheId = stackInput.Target.Spec.RedisCacheId.GetValue()

	// No Azure tags: ARM does not support tags on access policy
	// assignments (cache children), so the platform's identity tags live
	// on the cache.

	return locals
}
