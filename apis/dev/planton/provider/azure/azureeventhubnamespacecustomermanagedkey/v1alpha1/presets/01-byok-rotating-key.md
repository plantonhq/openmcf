# BYOK with Rotating Key

This preset applies customer-managed-key encryption onto a single-tenant
Event Hubs namespace, with a versionless key reference so vault-side
rotation propagates automatically.

## When to Use

- Compliance regimes requiring tenant-controlled encryption keys for
  event data at rest
- After the namespace exists with its identity attached and granted on
  the vault -- CMK is a second step by Azure's design (the sequencing a
  create-time block could never express)

## Key Configuration Choices

- **Versionless key reference** -- rotation-follows-latest; pin a
  versioned `key_id` only when compliance demands immutable versions
- **The identity contract** -- `userAssignedIdentityId` must already be
  in the namespace's identity block with wrap/unwrap vault access
  (a "Key Vault Crypto Service Encryption User" AzureRoleAssignment);
  omit it to use the namespace's system-assigned identity instead
- **ADD-ONLY** -- once applied, CMK can never be removed; deleting this
  resource changes nothing (returning to Microsoft-managed keys means
  replacing the namespace). The keys' vault needs purge protection

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-premium-hubs` | A PREMIUM or dedicated-cluster namespace | Your streaming composition |
| `my-streaming-key` | The AzureKeyVaultKey encrypting event data | Your key-management composition |
| `my-cmk-identity` | The unwrapping AzureUserAssignedIdentity | Your identity composition |
