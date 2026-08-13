---
title: "Immutable CMK Vault"
description: "This preset creates the compliance-grade vault: immutability enforced (reversibly, at `Unlocked`), backup data encrypted with your own Key Vault key through a user-assigned identity, and the public..."
type: "preset"
rank: "02"
presetSlug: "02-immutable-cmk-vault"
componentSlug: "recovery-services-vault"
componentTitle: "Recovery Services Vault"
provider: "azure"
icon: "package"
order: 2
---

# Immutable CMK Vault

This preset creates the compliance-grade vault: immutability enforced (reversibly, at `Unlocked`), backup data encrypted with your own Key Vault key through a user-assigned identity, and the public endpoint closed. For organizations whose backup posture must survive both ransomware and policy audits.

## When to Use

- Regulated environments where backup deletion and retention reduction must be blocked
- Data-at-rest policies that require customer-managed key material
- Private-endpoint-only network architectures

## Key Configuration Choices

- **`immutability: Unlocked`** -- full enforcement, still reversible. Graduate to `Locked` only after living with the retention settings; `Locked` is permanent (leaving it replaces the vault)
- **User-assigned identity for CMK** -- the Key Vault wrap/unwrap grant composes BEFORE the vault exists, avoiding the system-assigned bootstrap hop; `useSystemAssignedIdentity: false` points key unwrapping at it
- **Versionless key reference** -- key rotation propagates automatically, no vault update needed
- **`infrastructureEncryptionEnabled: true`** -- double encryption; decide now, it can NEVER change once encryption is on
- **`publicNetworkAccessEnabled: false`** -- pair with private endpoints for the vault's backup traffic

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the vault lives in | Your resource group resource's name |
| `<your-identity>` | The AzureUserAssignedIdentity that unwraps the key | Grant it wrap/unwrap (Key Vault Crypto Service Encryption User) on the key first |
| `<your-key-vault-key>` | The AzureKeyVaultKey encrypting backup data | Its versionless_id output wires automatically |

Consider adding `resourceGuardId` for multi-user authorization -- the guard belongs in a DIFFERENT administrator's scope to mean anything.
