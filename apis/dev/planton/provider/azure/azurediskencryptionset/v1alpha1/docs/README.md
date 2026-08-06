# AzureDiskEncryptionSet -- Design Research

## The Resource

An Azure Disk Encryption Set (`Microsoft.Compute/diskEncryptionSets`) binds
a customer-managed Key Vault key to managed disks, snapshots, and images for
server-side encryption at rest with a customer-managed key. The component
maps onto `azurerm_disk_encryption_set` (azurerm v4.x,
`internal/services/compute/disk_encryption_set_resource.go`), parity-verified
against pulumi-azure v6 (`compute.DiskEncryptionSet`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `key_vault_key_id` | `key_vault_key_id` | FK → AzureKeyVaultKey (versionless_id default) |
| `auto_key_rotation_enabled` | `auto_key_rotation_enabled` | Azure default false |
| `encryption_type` | `encryption_type` enum | ForceNew; default EncryptionAtRestWithCustomerKey |
| `federated_client_id` | `federated_client_id` | UUID; cross-tenant CMK |
| `identity` (Required) | `identity` message | System/User/both; identity_ids for user-assigned |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `id` (computed) | `disk_encryption_set_id` output | Join key for disk/VM/VMSS references |
| `identity.0.principal_id`/`tenant_id` (computed) | `identity_principal_id`/`_tenant_id` outputs | Grant target |

## Key Design Decisions

- **Versioned-vs-versionless key contract lives in the reference, not a
  CEL.** Azure enforces at apply that `auto_key_rotation_enabled = true`
  needs a VERSIONLESS key URL and `false` needs a VERSIONED one. A
  field-level CEL cannot express this because it cannot dereference a
  `StringValueOrRef`'s resolved value. The spec's default field path is
  `AzureKeyVaultKey.versionless_id` (the rotation-on posture); pointing at
  `key_id` (versioned) pairs with rotation off. Both modules pass the
  referenced id through and let the provider validate the pairing. This is
  the same protovalidate constraint recorded across the catalog for
  StringValueOrRef sub-field dereferencing.
- **Identity is required and modeled as a block.** A set cannot read its key
  without a managed identity. The system-assigned flavor creates the
  chicken-and-egg (its principal does not exist until the set is created, so
  the key grant comes after); the user-assigned flavor lets the grant be
  applied before the set exists. The `identity_ids`-match-type CEL mirrors
  the ACR/registry precedent.
- **Purge protection is an Azure runtime requirement, not a provider
  validator.** azurerm does not validate it, but Azure rejects a set whose
  vault lacks purge protection. This is documented on the spec and README
  rather than CEL'd (it is a property of the referenced vault, not this
  resource).

## Composition Seams

`AzureManagedDisk` (`disk_encryption_set_id`, `secure_vm_disk_encryption_set_id`),
`AzureVirtualMachine` (same two), and `AzureVirtualMachineScaleSet`
(OS-disk + data-disk sets) all reference this kind's
`disk_encryption_set_id` output; `AzureAksCluster` references it for node
OS-disk and persistent-volume CMK encryption. The set consumes
`AzureKeyVaultKey` (the key) and optionally `AzureUserAssignedIdentity` (the
unwrapping identity).

## Live E2E Exclusion (recorded)

Live E2E is profile-deferred because a disk encryption set requires a
purge-protection-enabled Key Vault, and a purge-protected vault cannot be
purged on teardown -- Azure holds the soft-deleted vault (and its name) for
the retention period (7+ days). Every other Azure vault-bearing E2E lane
purges its vault on destroy; this one structurally cannot, so an ephemeral
create-verify-destroy cycle cannot reach zero orphans in the shared test
subscription. The scenario and its fixtures (identity, purge-protected
vault, key, crypto grant) ship ready-to-run for a throwaway subscription;
the module is proven on the offline gate.
