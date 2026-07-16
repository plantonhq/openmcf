package module

import (
	azurestorageobjectreplicationv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestorageobjectreplication/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageObjectReplication *azurestorageobjectreplicationv1.AzureStorageObjectReplication
	SourceStorageAccountId        string
	DestinationStorageAccountId   string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestorageobjectreplicationv1.AzureStorageObjectReplicationStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageObjectReplication = stackInput.Target
	locals.SourceStorageAccountId = stackInput.Target.Spec.SourceStorageAccountId.GetValue()
	locals.DestinationStorageAccountId = stackInput.Target.Spec.DestinationStorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on
	// objectReplicationPolicies, so the platform's identity tags live on
	// the two accounts.

	return locals
}
