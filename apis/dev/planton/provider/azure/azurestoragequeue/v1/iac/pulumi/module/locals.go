package module

import (
	azurestoragequeuev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestoragequeue/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageQueue *azurestoragequeuev1.AzureStorageQueue
	StorageAccountId  string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestoragequeuev1.AzureStorageQueueStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageQueue = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on queueServices/queues,
	// so the platform's identity tags live on the parent account.

	return locals
}
