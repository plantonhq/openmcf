package module

import (
	azuremanagedredisaccesspolicyassignmentv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremanagedredisaccesspolicyassignment/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureManagedRedisAccessPolicyAssignment *azuremanagedredisaccesspolicyassignmentv1alpha1.AzureManagedRedisAccessPolicyAssignment
	ManagedRedisId                          string
	ObjectId                                string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremanagedredisaccesspolicyassignmentv1alpha1.AzureManagedRedisAccessPolicyAssignmentStackInput) *Locals {
	locals := &Locals{}

	locals.AzureManagedRedisAccessPolicyAssignment = stackInput.Target
	locals.ManagedRedisId = stackInput.Target.Spec.ManagedRedisId.GetValue()
	locals.ObjectId = stackInput.Target.Spec.ObjectId.GetValue()

	// No Azure tags: ARM does not support tags on access policy
	// assignments (database children), so the platform's identity tags
	// live on the cluster.

	return locals
}
