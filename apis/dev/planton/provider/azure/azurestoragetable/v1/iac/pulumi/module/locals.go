package module

import (
	azurestoragetablev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestoragetable/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageTable *azurestoragetablev1.AzureStorageTable
	StorageAccountId  string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestoragetablev1.AzureStorageTableStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageTable = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on tableServices/tables,
	// so the platform's identity tags live on the parent account.

	return locals
}
