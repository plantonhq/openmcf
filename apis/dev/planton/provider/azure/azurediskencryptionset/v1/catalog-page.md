# Azure Disk Encryption Set

Creates an Azure Disk Encryption Set — the resource that binds a customer-managed Key Vault key to managed disks, snapshots, and images for encryption at rest with a key you control. Managed disks, VMs, and VM scale sets reference the set by ARM ID to opt their disks into CMK encryption.

## What Gets Created

When you deploy an AzureDiskEncryptionSet resource, Planton provisions:

- **Disk Encryption Set** — an `azurerm_disk_encryption_set` with a managed identity, bound to a Key Vault key

Granting the set's identity crypto access on the key, and encrypting specific disks with the set, are separate steps (a role assignment and each disk's `diskEncryptionSetId`).

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A Key Vault key** (an `AzureKeyVaultKey`) on a vault with **purge protection enabled** — Azure requires it for disk encryption
- **Compute + Key Vault rights**: `Microsoft.Compute/diskEncryptionSets/write` plus the ability to grant the set crypto access on the key

## Quick Start

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureDiskEncryptionSet
metadata:
  name: prod-des
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureDiskEncryptionSet.prod-des
spec:
  region: eastus
  resourceGroup:
    value: platform-rg
  name: prod-des
  keyVaultKeyId:
    valueFrom:
      name: disk-cmk
  autoKeyRotationEnabled: true
  identity:
    type: SYSTEM_ASSIGNED
```

Deploy:

```shell
planton apply -f disk-encryption-set.yaml
```

After deployment, grant `status.outputs.identity_principal_id` the "Key Vault Crypto Service Encryption User" role on the key, then reference `status.outputs.disk_encryption_set_id` from your disks.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | Azure region; must match the disks. Fixed at creation. |
| `resourceGroup` | `StringValueOrRef` | Resource group. Defaults to an `AzureResourceGroup` reference. |
| `name` | `string` | Set name, unique in the resource group. Fixed at creation. |
| `keyVaultKeyId` | `StringValueOrRef` | The CMK. Defaults to an `AzureKeyVaultKey` `versionless_id` reference. |
| `identity` | `object` | The set's managed identity (`type`: `SYSTEM_ASSIGNED` / `USER_ASSIGNED` / `SYSTEM_AND_USER_ASSIGNED`; `identityIds` for the user-assigned flavors). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `autoKeyRotationEnabled` | `bool` | `false` | Follow the key's latest version. Requires a versionless key when true, a versioned key when false. |
| `encryptionType` | `enum` | `ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY` | Customer-key / platform+customer / confidential-VM. Fixed at creation. |
| `federatedClientId` | `string` | `""` | Multi-tenant app client id for cross-tenant CMK (UUID). |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags (user wins). |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `disk_encryption_set_id` | `string` | Full ARM ID — referenced by disks, VMs, and scale sets |
| `disk_encryption_set_name` | `string` | The set's name as deployed |
| `identity_principal_id` | `string` | Principal to grant Key Vault crypto access (system-assigned) |
| `identity_tenant_id` | `string` | Tenant of the system-assigned identity |

## Related Components

- [AzureKeyVaultKey](/docs/catalog/azure/key-vault-key) — the customer-managed key the set wraps disks with
- [AzureManagedDisk](/docs/catalog/azure/managed-disk) — references the set via `diskEncryptionSetId`
- [AzureVirtualMachine](/docs/catalog/azure/virtual-machine) — references the set for OS/data disk CMK encryption
- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the identity you grant vault access before the set exists
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants the set's identity crypto access on the key
