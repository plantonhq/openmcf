package module

import (
	azuredataprotectionbackupinstancev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredataprotectionbackupinstance/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureDataProtectionBackupInstance *azuredataprotectionbackupinstancev1alpha1.AzureDataProtectionBackupInstance

	// VaultId and BackupPolicyId are StringValueOrRef fields; the
	// platform middleware resolves valueFrom references before IaC
	// modules run, so GetValue() always returns the resolved literal
	// ARM ID. Per-variant datasource references resolve the same way
	// inside each variant's create function.
	VaultId        string
	BackupPolicyId string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuredataprotectionbackupinstancev1alpha1.AzureDataProtectionBackupInstanceStackInput) *Locals {
	locals := &Locals{}

	locals.AzureDataProtectionBackupInstance = stackInput.Target

	locals.VaultId = stackInput.Target.Spec.VaultId.GetValue()
	locals.BackupPolicyId = stackInput.Target.Spec.BackupPolicyId.GetValue()

	// Note: backup instances carry NO tags argument (the provider has
	// none on any of the six variant resources) -- there is no tag map
	// here, unlike sibling modules.
	return locals
}
