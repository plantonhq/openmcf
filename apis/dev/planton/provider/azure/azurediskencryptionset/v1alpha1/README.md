# AzureDiskEncryptionSet

## Overview

`AzureDiskEncryptionSet` provisions an Azure Disk Encryption Set: the
resource that binds a customer-managed key (CMK) in Key Vault to managed
disks, snapshots, and images so their data is encrypted at rest with a key
you control rather than a platform-managed key. Managed disks, VMs, and VM
scale sets reference a set by ARM ID (`disk_encryption_set_id`) to opt their
disks into CMK encryption.

## Why a First-Class Resource?

A disk encryption set has its own lifecycle and is referenced by many
compute resources:

- **One key, many disks** -- a single set encrypts every disk that
  references it; rotate the key once and every disk follows
- **Referenced, not inline** -- `AzureManagedDisk`, `AzureVirtualMachine`,
  and `AzureVirtualMachineScaleSet` all point at the set's ARM ID
- **A distinct identity to grant** -- the set carries a managed identity
  that must be granted crypto access on the key; that grant is its own
  composable edge

## Key Requirements

- **Purge protection** -- the referenced Key Vault must have purge
  protection enabled; Azure requires it for any vault backing disk
  encryption
- **Key access** -- the set's identity must be granted "Key Vault Crypto
  Service Encryption User" (RBAC vault) or get/wrapKey/unwrapKey (access
  policy vault) on the key before disks can use the set
- **Key versioning matches rotation** -- a versionless key with
  `auto_key_rotation_enabled = true` (recommended), or a versioned key with
  it false

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | Azure region (must match the disks); fixed at creation |
| `resource_group` | StringValueOrRef | Yes | Resource group (defaults to AzureResourceGroup) |
| `name` | string | Yes | Set name, unique in the resource group; fixed at creation |
| `key_vault_key_id` | StringValueOrRef | Yes | The CMK (defaults to AzureKeyVaultKey `versionless_id`) |
| `auto_key_rotation_enabled` | bool | No | Follow the key's latest version (default false) |
| `encryption_type` | enum | No | Customer-key / platform+customer / confidential-VM; fixed at creation |
| `federated_client_id` | string | No | Multi-tenant app client id for cross-tenant CMK (UUID) |
| `identity` | message | Yes | The set's managed identity (system/user-assigned) |
| `tags` | map | No | User tags, merged over Planton-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `disk_encryption_set_id` | Full ARM ID -- the join key disks/VMs/scale sets reference |
| `disk_encryption_set_name` | The set's name as deployed |
| `identity_principal_id` | Principal to grant Key Vault crypto access (system-assigned) |
| `identity_tenant_id` | Tenant of the system-assigned identity |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDiskEncryptionSet
metadata:
  name: prod-des
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: prod-des
  keyVaultKeyId:
    valueFrom:
      name: disk-cmk
  autoKeyRotationEnabled: true
  identity:
    type: SYSTEM_ASSIGNED
```

Then grant the set's identity crypto access and reference it from a disk:

```yaml
spec:
  diskEncryptionSetId:
    valueFrom:
      name: prod-des
```

## Lifecycle Notes

- `region`, `name`, and `encryption_type` are **fixed at creation**
- With a system-assigned identity, grant `identity_principal_id` crypto
  access AFTER the set is created (its principal does not exist before);
  with a user-assigned identity, grant it BEFORE and reference it, avoiding
  the chicken-and-egg
- Changing the identity from a user-assigned flavor to system-assigned
  replaces the set

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
