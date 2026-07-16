# User-Assigned Set with Pre-Provisioned Access

This preset creates a disk encryption set backed by a user-assigned managed
identity. Because the identity exists independently of the set, you grant it
Key Vault crypto access BEFORE creating the set -- avoiding the
system-assigned chicken-and-egg where the principal only exists after
creation. This is the cleaner pattern for infra-as-code, where the identity,
the grant, and the set all deploy in one dependency-ordered pass.

## When to Use

- IaC pipelines that want the key grant applied before the set exists
- Multiple sets or resources sharing one managed identity and grant
- Cross-resource CMK setups where the identity is provisioned centrally

## Key Configuration Choices

- **`identity.type: USER_ASSIGNED`** with `identityIds` -- brings an
  identity you manage; grant it crypto access on the key first
- **`autoKeyRotationEnabled: true`** -- versionless key reference; hands-off
  rotation

## Requirements

- The referenced Key Vault must have **purge protection enabled**
- The referenced identity must be granted "Key Vault Crypto Service
  Encryption User" on the key before the set is used by disks

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the set in | The resource group's `status.outputs.resource_group_name` |
| `<disk-cmk-key>` | The AzureKeyVaultKey to encrypt disks with | Your key's Planton resource name |
| `<des-identity>` | The AzureUserAssignedIdentity the set uses | Your identity's Planton resource name |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
