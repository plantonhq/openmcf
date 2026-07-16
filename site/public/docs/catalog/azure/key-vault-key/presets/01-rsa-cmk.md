---
title: "RSA Customer-Managed Key"
description: "This preset creates the baseline CMK: a software-protected RSA-2048 key whose capability list is exactly the envelope-encryption pair (`WRAP_KEY` + `UNWRAP_KEY`). Every Azure CMK integration --..."
type: "preset"
rank: "01"
presetSlug: "01-rsa-cmk"
componentSlug: "key-vault-key"
componentTitle: "Key Vault Key"
provider: "azure"
icon: "package"
order: 1
---

# RSA Customer-Managed Key

This preset creates the baseline CMK: a software-protected RSA-2048 key
whose capability list is exactly the envelope-encryption pair
(`WRAP_KEY` + `UNWRAP_KEY`). Every Azure CMK integration -- storage
accounts, disk encryption sets, container registries, database services --
accepts this shape.

## When to Use

- Bringing your own key to any Azure service's encryption-at-rest
- The starting point before compliance requirements push you to HSM
  protection (see the auto-rotating HSM preset)

## Key Configuration Choices

- **RSA-2048** -- the universal baseline; move to 3072/4096 only when a
  compliance regime demands it (bigger keys wrap slower)
- **`keyOpts` limited to WRAP_KEY/UNWRAP_KEY** -- the key can do nothing
  else; Azure rejects encrypt/sign attempts, which is the least-privilege
  posture for a pure CMK
- **Consumers reference `versionless_id`** -- rotation (manual or policy)
  propagates without touching the consumer

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<key-vault-arm-id>` | The vault's ARM resource ID (or a `valueFrom` reference to an AzureKeyVault's `key_vault_id` output) | The vault's `status.outputs.key_vault_id` |
| `<purpose>` | What this key encrypts (tag value) | Your tagging convention |

## Downstream Wiring

```yaml
spec:
  encryption:
    identityClientId: <unwrapping-identity-client-id>
    keyVaultKeyId:
      valueFrom:
        name: storage-cmk
```

The consuming service's identity needs `KEY_WRAP_KEY`/`KEY_UNWRAP_KEY`
(access-policy mode) or the "Key Vault Crypto Service Encryption User"
role (RBAC mode) on the vault.
