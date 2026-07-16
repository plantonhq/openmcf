package module

import (
	azurestorageencryptionscopev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestorageencryptionscope/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageEncryptionScope *azurestorageencryptionscopev1.AzureStorageEncryptionScope
	StorageAccountId            string
}

// sourceStrings maps the spec's source enum to ARM's dotted wire values.
var sourceStrings = map[azurestorageencryptionscopev1.AzureStorageEncryptionScopeSource]string{
	azurestorageencryptionscopev1.AzureStorageEncryptionScopeSource_MICROSOFT_STORAGE:   "Microsoft.Storage",
	azurestorageencryptionscopev1.AzureStorageEncryptionScopeSource_MICROSOFT_KEY_VAULT: "Microsoft.KeyVault",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestorageencryptionscopev1.AzureStorageEncryptionScopeStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageEncryptionScope = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on encryptionScopes, so
	// the platform's identity tags live on the parent account.

	return locals
}
