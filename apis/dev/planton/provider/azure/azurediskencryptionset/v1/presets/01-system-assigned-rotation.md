# System-Assigned Set with Auto Key Rotation

This preset creates a disk encryption set with a system-assigned identity
and automatic key rotation -- the recommended posture. The set follows the
key's latest version as it rotates (referencing the key's `versionless_id`),
so a key rotation propagates to every disk without redeploying.

Because the identity is system-assigned, its principal does not exist until
the set is created; grant it "Key Vault Crypto Service Encryption User" on
the key AFTER deployment (read `status.outputs.identity_principal_id`).

## When to Use

- Standard CMK disk encryption where you want hands-off key rotation
- A single set encrypting many disks/VMs in one region

## Key Configuration Choices

- **`autoKeyRotationEnabled: true`** -- requires a versionless key reference
  (the default `AzureKeyVaultKey.versionless_id`); the set tracks new
  versions automatically
- **`identity.type: SYSTEM_ASSIGNED`** -- Azure manages the identity's
  lifecycle with the set; grant it crypto access after creation

## Requirements

- The referenced Key Vault must have **purge protection enabled** (Azure
  requires it for disk encryption)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the set in | The resource group's `status.outputs.resource_group_name` |
| `<disk-cmk-key>` | The AzureKeyVaultKey to encrypt disks with | Your key's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Encrypt a managed disk with this set:

```yaml
spec:
  diskEncryptionSetId:
    valueFrom:
      name: my-disk-encryption-set
```
