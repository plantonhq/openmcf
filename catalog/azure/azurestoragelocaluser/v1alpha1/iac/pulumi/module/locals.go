package module

import (
	azurestoragelocaluserv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurestoragelocaluser/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageLocalUser *azurestoragelocaluserv1alpha1.AzureStorageLocalUser
	StorageAccountId      string
}

// serviceStrings maps the spec's permission-scope service enum to the
// API's lowercase wire values.
var serviceStrings = map[azurestoragelocaluserv1alpha1.AzureStorageLocalUserPermissionService]string{
	azurestoragelocaluserv1alpha1.AzureStorageLocalUserPermissionService_BLOB: "blob",
	azurestoragelocaluserv1alpha1.AzureStorageLocalUserPermissionService_FILE: "file",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestoragelocaluserv1alpha1.AzureStorageLocalUserStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageLocalUser = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on localUsers, so the
	// platform's identity tags live on the parent account.

	return locals
}
