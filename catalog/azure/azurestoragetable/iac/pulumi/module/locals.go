package module

import (
	azurestoragetablev1alpha1 "github.com/plantonhq/planton/catalog/azure/azurestoragetable/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageTable *azurestoragetablev1alpha1.AzureStorageTable
	StorageAccountId  string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestoragetablev1alpha1.AzureStorageTableStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageTable = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on tableServices/tables,
	// so the platform's identity tags live on the parent account.

	return locals
}
