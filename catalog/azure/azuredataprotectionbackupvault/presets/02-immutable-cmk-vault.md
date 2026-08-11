# Immutable CMK Vault

This preset creates the compliance-grade vault: trial-run immutability, a 30-day soft-delete window, and backup data encrypted with your own Key Vault key. For organizations whose backup posture must survive both ransomware and a compromised administrator.

## When to Use

- Regulated environments where backup retention must not be quietly reducible
- Data-at-rest policies that require customer-controlled key material
- The step before `Locked` immutability -- prove the retention settings here first

## Key Configuration Choices

- **`immutability: Unlocked`** -- blocks backup deletion, still reversible; move to `Locked` only when certain (Locked is PERMANENT -- leaving it replaces the vault)
- **`retentionDurationInDays: 30`** -- deleted backups stay recoverable for a month; retention beyond 14 days may incur charges
- **`identity.type: SYSTEM_ASSIGNED`** -- required for CMK; Azure unwraps the key with the system identity (hardcoded service-side)
- **`encryption.keyId` via the versionless reference** -- key rotation propagates automatically, without touching the vault; note CMK is a ONE-WAY door (it can never be removed, only the key rotates)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the vault lives in | Your resource group resource's name |
| `<your-key-vault-key>` | The AzureKeyVaultKey encrypting backup data | Your Key Vault key resource's name |

Grant the vault's system-assigned identity wrap/unwrap access on the key BEFORE deploying with the encryption block -- the deploy fails without it.
