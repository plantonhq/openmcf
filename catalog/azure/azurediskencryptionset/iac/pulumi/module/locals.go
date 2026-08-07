package module

import (
	"strings"

	azurediskencryptionsetv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurediskencryptionset/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureDiskEncryptionSet *azurediskencryptionsetv1alpha1.AzureDiskEncryptionSet

	ResourceGroupName string

	// KeyVaultKeyId is the resolved literal Key Vault key URL (versionless
	// for rotation-on, versioned for a pinned version).
	KeyVaultKeyId string

	// EncryptionType is the ARM string for the spec enum, or empty when
	// unspecified so both engines let Azure apply its default
	// (EncryptionAtRestWithCustomerKey).
	EncryptionType string

	// IdentityType is the ARM identity string; IdentityIds are the resolved
	// user-assigned identity ARM IDs.
	IdentityType string
	IdentityIds  []string

	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetStackInput) *Locals {
	locals := &Locals{}

	locals.AzureDiskEncryptionSet = stackInput.Target
	target := stackInput.Target
	spec := target.Spec

	locals.ResourceGroupName = spec.ResourceGroup.GetValue()
	locals.KeyVaultKeyId = spec.KeyVaultKeyId.GetValue()

	switch spec.EncryptionType {
	case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetEncryptionType_ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY:
		locals.EncryptionType = "EncryptionAtRestWithCustomerKey"
	case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetEncryptionType_ENCRYPTION_AT_REST_WITH_PLATFORM_AND_CUSTOMER_KEYS:
		locals.EncryptionType = "EncryptionAtRestWithPlatformAndCustomerKeys"
	case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetEncryptionType_CONFIDENTIAL_VM_ENCRYPTED_WITH_CUSTOMER_KEY:
		locals.EncryptionType = "ConfidentialVmEncryptedWithCustomerKey"
	}

	if spec.Identity != nil {
		switch spec.Identity.Type {
		case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetIdentityType_SYSTEM_ASSIGNED:
			locals.IdentityType = "SystemAssigned"
		case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetIdentityType_USER_ASSIGNED:
			locals.IdentityType = "UserAssigned"
		case azurediskencryptionsetv1alpha1.AzureDiskEncryptionSetIdentityType_SYSTEM_AND_USER_ASSIGNED:
			locals.IdentityType = "SystemAssigned, UserAssigned"
		}
		for _, id := range spec.Identity.IdentityIds {
			locals.IdentityIds = append(locals.IdentityIds, id.GetValue())
		}
	}

	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module emits the family-wide snake-case literal
	// and falls back to metadata.name. Output-neutral (tags never feed stack
	// outputs); aligning the two shapes is a family-wide convention change.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureDiskEncryptionSet.String()),
	}
	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}
	for k, v := range spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
