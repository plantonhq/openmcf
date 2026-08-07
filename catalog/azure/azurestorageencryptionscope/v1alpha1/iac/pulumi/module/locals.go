package module

import (
	azurestorageencryptionscopev1alpha1 "github.com/plantonhq/planton/catalog/azure/azurestorageencryptionscope/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageEncryptionScope *azurestorageencryptionscopev1alpha1.AzureStorageEncryptionScope
	StorageAccountId            string
}

// sourceStrings maps the spec's source enum to ARM's dotted wire values.
var sourceStrings = map[azurestorageencryptionscopev1alpha1.AzureStorageEncryptionScopeSource]string{
	azurestorageencryptionscopev1alpha1.AzureStorageEncryptionScopeSource_MICROSOFT_STORAGE:   "Microsoft.Storage",
	azurestorageencryptionscopev1alpha1.AzureStorageEncryptionScopeSource_MICROSOFT_KEY_VAULT: "Microsoft.KeyVault",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestorageencryptionscopev1alpha1.AzureStorageEncryptionScopeStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageEncryptionScope = stackInput.Target
	locals.StorageAccountId = stackInput.Target.Spec.StorageAccountId.GetValue()

	// No Azure tags: ARM does not support tags on encryptionScopes, so
	// the platform's identity tags live on the parent account.

	return locals
}
