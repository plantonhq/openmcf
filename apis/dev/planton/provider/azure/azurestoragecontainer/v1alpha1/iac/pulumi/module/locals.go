package module

import (
	azurestoragecontainerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestoragecontainer/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageContainer *azurestoragecontainerv1alpha1.AzureStorageContainer
	StorageAccountId      string
}

// containerAccessTypeStrings maps the spec's access-type enum to azurerm's
// lowercase wire values. Unspecified materializes "private" in main.go --
// the container is born locked down unless the spec says otherwise.
var containerAccessTypeStrings = map[azurestoragecontainerv1alpha1.AzureStorageContainerAccessType]string{
	azurestoragecontainerv1alpha1.AzureStorageContainerAccessType_PRIVATE:   "private",
	azurestoragecontainerv1alpha1.AzureStorageContainerAccessType_BLOB:      "blob",
	azurestoragecontainerv1alpha1.AzureStorageContainerAccessType_CONTAINER: "container",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestoragecontainerv1alpha1.AzureStorageContainerStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageContainer = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on blobServices/containers,
	// so the platform's identity tags live on the parent account.

	return locals
}
