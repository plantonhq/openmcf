package module

import (
	"strings"

	azurestorageaccountv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurestorageaccount/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureStorageAccount *azurestorageaccountv1.AzureStorageAccount
	ResourceGroupName   string
	AzureTags           map[string]string
}

// accountKindStrings maps the spec's account-kind enum to ARM's values.
// Unspecified materializes StorageV2 (the spec's documented default) in
// main.go -- stack inputs never carry proto defaults.
var accountKindStrings = map[azurestorageaccountv1.AzureStorageAccountKind]string{
	azurestorageaccountv1.AzureStorageAccountKind_STORAGE_V2:         "StorageV2",
	azurestorageaccountv1.AzureStorageAccountKind_BLOB_STORAGE:       "BlobStorage",
	azurestorageaccountv1.AzureStorageAccountKind_BLOCK_BLOB_STORAGE: "BlockBlobStorage",
	azurestorageaccountv1.AzureStorageAccountKind_FILE_STORAGE:       "FileStorage",
	azurestorageaccountv1.AzureStorageAccountKind_STORAGE:            "Storage",
}

// accountTierStrings maps the spec's tier enum to ARM's values.
var accountTierStrings = map[azurestorageaccountv1.AzureStorageAccountTier]string{
	azurestorageaccountv1.AzureStorageAccountTier_STANDARD: "Standard",
	azurestorageaccountv1.AzureStorageAccountTier_PREMIUM:  "Premium",
}

// replicationTypeStrings maps the spec's replication enum to azurerm's SKU
// suffixes (the RA_ prefix collapses: RA_GRS -> RAGRS).
var replicationTypeStrings = map[azurestorageaccountv1.AzureStorageAccountReplicationType]string{
	azurestorageaccountv1.AzureStorageAccountReplicationType_LRS:     "LRS",
	azurestorageaccountv1.AzureStorageAccountReplicationType_ZRS:     "ZRS",
	azurestorageaccountv1.AzureStorageAccountReplicationType_GRS:     "GRS",
	azurestorageaccountv1.AzureStorageAccountReplicationType_GZRS:    "GZRS",
	azurestorageaccountv1.AzureStorageAccountReplicationType_RA_GRS:  "RAGRS",
	azurestorageaccountv1.AzureStorageAccountReplicationType_RA_GZRS: "RAGZRS",
}

// accessTierStrings maps the spec's access-tier enum to ARM's values.
// Unspecified is not sent at all -- Azure computes Hot on the kinds that
// support tiers, mirroring the Terraform module's null.
var accessTierStrings = map[azurestorageaccountv1.AzureStorageAccountAccessTier]string{
	azurestorageaccountv1.AzureStorageAccountAccessTier_HOT:                 "Hot",
	azurestorageaccountv1.AzureStorageAccountAccessTier_COOL:                "Cool",
	azurestorageaccountv1.AzureStorageAccountAccessTier_COLD:                "Cold",
	azurestorageaccountv1.AzureStorageAccountAccessTier_ACCESS_TIER_PREMIUM: "Premium",
}

// minTlsVersionStrings maps the spec's TLS-floor enum; ARM's values happen
// to match the proto value names verbatim.
var minTlsVersionStrings = map[azurestorageaccountv1.AzureStorageAccountMinTlsVersion]string{
	azurestorageaccountv1.AzureStorageAccountMinTlsVersion_TLS1_0: "TLS1_0",
	azurestorageaccountv1.AzureStorageAccountMinTlsVersion_TLS1_1: "TLS1_1",
	azurestorageaccountv1.AzureStorageAccountMinTlsVersion_TLS1_2: "TLS1_2",
}

// allowedCopyScopeStrings maps the copy-scope restriction enum. Unspecified
// is not sent -- copy stays unrestricted (Azure's default).
var allowedCopyScopeStrings = map[azurestorageaccountv1.AzureStorageAccountAllowedCopyScope]string{
	azurestorageaccountv1.AzureStorageAccountAllowedCopyScope_AAD:          "AAD",
	azurestorageaccountv1.AzureStorageAccountAllowedCopyScope_PRIVATE_LINK: "PrivateLink",
}

// dnsEndpointTypeStrings maps the DNS-architecture enum. Unspecified is not
// sent -- azurerm defaults the create-only choice to Standard itself.
var dnsEndpointTypeStrings = map[azurestorageaccountv1.AzureStorageAccountDnsEndpointType]string{
	azurestorageaccountv1.AzureStorageAccountDnsEndpointType_DNS_ENDPOINT_STANDARD: "Standard",
	azurestorageaccountv1.AzureStorageAccountDnsEndpointType_AZURE_DNS_ZONE:        "AzureDnsZone",
}

// encryptionKeyTypeStrings maps the queue/table key-scope enum. Unspecified
// is not sent -- Azure defaults to the Service scope.
var encryptionKeyTypeStrings = map[azurestorageaccountv1.AzureStorageAccountEncryptionKeyType]string{
	azurestorageaccountv1.AzureStorageAccountEncryptionKeyType_SERVICE: "Service",
	azurestorageaccountv1.AzureStorageAccountEncryptionKeyType_ACCOUNT: "Account",
}

// identityTypeStrings maps the managed-identity flavor enum to ARM's
// comma-separated vocabulary.
var identityTypeStrings = map[azurestorageaccountv1.AzureStorageAccountIdentityType]string{
	azurestorageaccountv1.AzureStorageAccountIdentityType_SYSTEM_ASSIGNED:          "SystemAssigned",
	azurestorageaccountv1.AzureStorageAccountIdentityType_USER_ASSIGNED:            "UserAssigned",
	azurestorageaccountv1.AzureStorageAccountIdentityType_SYSTEM_AND_USER_ASSIGNED: "SystemAssigned, UserAssigned",
}

// networkDefaultActionStrings maps the firewall default-action enum.
var networkDefaultActionStrings = map[azurestorageaccountv1.AzureStorageAccountNetworkDefaultAction]string{
	azurestorageaccountv1.AzureStorageAccountNetworkDefaultAction_ALLOW: "Allow",
	azurestorageaccountv1.AzureStorageAccountNetworkDefaultAction_DENY:  "Deny",
}

// networkBypassStrings maps the firewall bypass-class enum.
var networkBypassStrings = map[azurestorageaccountv1.AzureStorageAccountNetworkBypass]string{
	azurestorageaccountv1.AzureStorageAccountNetworkBypass_AZURE_SERVICES: "AzureServices",
	azurestorageaccountv1.AzureStorageAccountNetworkBypass_LOGGING:        "Logging",
	azurestorageaccountv1.AzureStorageAccountNetworkBypass_METRICS:        "Metrics",
	azurestorageaccountv1.AzureStorageAccountNetworkBypass_NONE:           "None",
}

// routingChoiceStrings maps the routing-preference enum. Unspecified
// materializes MicrosoftRouting (Azure's default) when the block is present.
var routingChoiceStrings = map[azurestorageaccountv1.AzureStorageAccountRoutingChoice]string{
	azurestorageaccountv1.AzureStorageAccountRoutingChoice_MICROSOFT_ROUTING: "MicrosoftRouting",
	azurestorageaccountv1.AzureStorageAccountRoutingChoice_INTERNET_ROUTING:  "InternetRouting",
}

// sasExpirationActionStrings maps the SAS-policy action enum. Unspecified
// materializes Log (Azure's default) when the block is present.
var sasExpirationActionStrings = map[azurestorageaccountv1.AzureStorageAccountSasExpirationAction]string{
	azurestorageaccountv1.AzureStorageAccountSasExpirationAction_LOG:   "Log",
	azurestorageaccountv1.AzureStorageAccountSasExpirationAction_BLOCK: "Block",
}

// immutabilityStateStrings maps the account-level WORM state enum.
var immutabilityStateStrings = map[azurestorageaccountv1.AzureStorageAccountImmutabilityState]string{
	azurestorageaccountv1.AzureStorageAccountImmutabilityState_DISABLED: "Disabled",
	azurestorageaccountv1.AzureStorageAccountImmutabilityState_UNLOCKED: "Unlocked",
	azurestorageaccountv1.AzureStorageAccountImmutabilityState_LOCKED:   "Locked",
}

// directoryTypeStrings maps the Azure Files directory-service enum; ARM's
// values match the proto value names verbatim.
var directoryTypeStrings = map[azurestorageaccountv1.AzureStorageAccountDirectoryServiceType]string{
	azurestorageaccountv1.AzureStorageAccountDirectoryServiceType_AADDS:   "AADDS",
	azurestorageaccountv1.AzureStorageAccountDirectoryServiceType_AADKERB: "AADKERB",
	azurestorageaccountv1.AzureStorageAccountDirectoryServiceType_AD:      "AD",
}

// defaultSharePermissionStrings maps the Azure Files default-share-permission
// enum to ARM's role-name vocabulary.
var defaultSharePermissionStrings = map[azurestorageaccountv1.AzureStorageAccountDefaultSharePermission]string{
	azurestorageaccountv1.AzureStorageAccountDefaultSharePermission_SHARE_PERMISSION_NONE:                 "None",
	azurestorageaccountv1.AzureStorageAccountDefaultSharePermission_SHARE_PERMISSION_READER:               "StorageFileDataSmbShareReader",
	azurestorageaccountv1.AzureStorageAccountDefaultSharePermission_SHARE_PERMISSION_CONTRIBUTOR:          "StorageFileDataSmbShareContributor",
	azurestorageaccountv1.AzureStorageAccountDefaultSharePermission_SHARE_PERMISSION_ELEVATED_CONTRIBUTOR: "StorageFileDataSmbShareElevatedContributor",
}

// lifecycleBlobTypeStrings maps the lifecycle blob-type enum to ARM's
// camelCase wire values.
var lifecycleBlobTypeStrings = map[azurestorageaccountv1.AzureStorageAccountLifecycleBlobType]string{
	azurestorageaccountv1.AzureStorageAccountLifecycleBlobType_BLOCK_BLOB:  "blockBlob",
	azurestorageaccountv1.AzureStorageAccountLifecycleBlobType_APPEND_BLOB: "appendBlob",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurestorageaccountv1.AzureStorageAccountStackInput) *Locals {
	locals := &Locals{}

	locals.AzureStorageAccount = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// PARITY-EXCEPTION: resource_kind here is the lowered
	// CloudResourceKind enum string and resource_id is omitted when
	// metadata.id is empty, while the Terraform module emits the
	// family-wide snake-case literal and falls back to metadata.name.
	// Output-neutral (tags never feed stack outputs); aligning the two
	// shapes is a family-wide convention change, not a per-kind fix.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureStorageAccount.String()),
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

	// The user's spec tags merge over the metadata-derived tags -- user
	// tags deliberately win so an org's governance conventions can
	// override the derived values where they collide.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
