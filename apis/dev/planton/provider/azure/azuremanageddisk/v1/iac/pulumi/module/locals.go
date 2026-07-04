package module

import (
	"strings"

	azuremanageddiskv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremanageddisk/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureManagedDisk *azuremanageddiskv1.AzureManagedDisk

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// StorageAccountType and CreateOption are the ARM strings for the
	// spec's required enums.
	StorageAccountType string
	CreateOption       string

	// OsType, HyperVGeneration, SecurityType, and NetworkAccessPolicy are
	// the ARM strings for the spec's optional enums, or empty when
	// unspecified so both engines send nothing and Azure's defaults apply.
	OsType              string
	HyperVGeneration    string
	SecurityType        string
	NetworkAccessPolicy string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremanageddiskv1.AzureManagedDiskStackInput) *Locals {
	locals := &Locals{}

	locals.AzureManagedDisk = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	switch target.Spec.StorageAccountType {
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_STANDARD_LRS:
		locals.StorageAccountType = "Standard_LRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_STANDARD_SSD_LRS:
		locals.StorageAccountType = "StandardSSD_LRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_STANDARD_SSD_ZRS:
		locals.StorageAccountType = "StandardSSD_ZRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_PREMIUM_LRS:
		locals.StorageAccountType = "Premium_LRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_PREMIUM_ZRS:
		locals.StorageAccountType = "Premium_ZRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_PREMIUM_V2_LRS:
		locals.StorageAccountType = "PremiumV2_LRS"
	case azuremanageddiskv1.AzureManagedDiskStorageAccountType_ULTRA_SSD_LRS:
		locals.StorageAccountType = "UltraSSD_LRS"
	}

	switch target.Spec.CreateOption {
	case azuremanageddiskv1.AzureManagedDiskCreateOption_EMPTY:
		locals.CreateOption = "Empty"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_COPY:
		locals.CreateOption = "Copy"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_FROM_IMAGE:
		locals.CreateOption = "FromImage"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_IMPORT:
		locals.CreateOption = "Import"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_IMPORT_SECURE:
		locals.CreateOption = "ImportSecure"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_RESTORE:
		locals.CreateOption = "Restore"
	case azuremanageddiskv1.AzureManagedDiskCreateOption_UPLOAD:
		locals.CreateOption = "Upload"
	}

	switch target.Spec.OsType {
	case azuremanageddiskv1.AzureManagedDiskOsType_LINUX:
		locals.OsType = "Linux"
	case azuremanageddiskv1.AzureManagedDiskOsType_WINDOWS:
		locals.OsType = "Windows"
	}

	switch target.Spec.HyperVGeneration {
	case azuremanageddiskv1.AzureManagedDiskHyperVGeneration_V1:
		locals.HyperVGeneration = "V1"
	case azuremanageddiskv1.AzureManagedDiskHyperVGeneration_V2:
		locals.HyperVGeneration = "V2"
	}

	switch target.Spec.SecurityType {
	case azuremanageddiskv1.AzureManagedDiskSecurityType_CONFIDENTIAL_VM_VMGUEST_STATE_ONLY_ENCRYPTED_WITH_PLATFORM_KEY:
		locals.SecurityType = "ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey"
	case azuremanageddiskv1.AzureManagedDiskSecurityType_CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_PLATFORM_KEY:
		locals.SecurityType = "ConfidentialVM_DiskEncryptedWithPlatformKey"
	case azuremanageddiskv1.AzureManagedDiskSecurityType_CONFIDENTIAL_VM_DISK_ENCRYPTED_WITH_CUSTOMER_KEY:
		locals.SecurityType = "ConfidentialVM_DiskEncryptedWithCustomerKey"
	}

	switch target.Spec.NetworkAccessPolicy {
	case azuremanageddiskv1.AzureManagedDiskNetworkAccessPolicy_ALLOW_ALL:
		locals.NetworkAccessPolicy = "AllowAll"
	case azuremanageddiskv1.AzureManagedDiskNetworkAccessPolicy_ALLOW_PRIVATE:
		locals.NetworkAccessPolicy = "AllowPrivate"
	case azuremanageddiskv1.AzureManagedDiskNetworkAccessPolicy_DENY_ALL:
		locals.NetworkAccessPolicy = "DenyAll"
	}

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureManagedDisk.String()),
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

	for k, v := range target.Spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
