---
title: "EC Signing Key"
description: "This preset creates an elliptic-curve signing key: the private key signs inside the vault (`SIGN`), anyone with the public half verifies (`VERIFY`). The `public_key_pem` output exports the public..."
type: "preset"
rank: "03"
presetSlug: "03-ec-signing"
componentSlug: "key-vault-key"
componentTitle: "Key Vault Key"
provider: "azure"
icon: "package"
order: 3
---

# EC Signing Key

This preset creates an elliptic-curve signing key: the private key signs
inside the vault (`SIGN`), anyone with the public half verifies
(`VERIFY`). The `public_key_pem` output exports the public half so
verifiers never need vault access at all.

## When to Use

- JWT/token signing where the private key must never live in application
  memory or configuration
- Code signing, webhook signing, or any detached-signature flow
- Workloads preferring EC's smaller signatures and faster verification
  over RSA

## Key Configuration Choices

- **`curve: P_384`** -- a step above the interoperable P-256 baseline;
  drop to `P_256` for maximum ecosystem compatibility or use `P_256K`
  for blockchain-ecosystem compatibility
- **No `keySize`** -- an EC key's strength is its curve; the spec rejects
  a size on EC keys (and a curve on RSA keys) before Azure would
- **`keyOpts` limited to SIGN/VERIFY** -- EC keys cannot encrypt in Key
  Vault; the capability list makes that explicit

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | The vault's ARM resource ID (or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<purpose>` | What this key signs (tag value) | Your tagging convention |

## Downstream Wiring

Distribute `status.outputs.public_key_pem` to verifiers; signers call the
vault's sign operation with the "Key Vault Crypto User" role (RBAC) or
`KEY_SIGN` permission (access-policy mode).
