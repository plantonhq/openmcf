# AzureKeyVaultKey

## Overview

`AzureKeyVaultKey` provisions a cryptographic key inside an Azure Key Vault
-- the customer-managed-key (CMK) building block of an Azure platform. The
private part never leaves the vault (or its HSM on the Premium tier);
consumers call the vault to encrypt/decrypt, wrap/unwrap, or sign/verify.

This is how "bring your own key" works across Azure: Storage accounts, disk
encryption sets, container registries, and database services all encrypt
their data with a data-encryption key that THIS key wraps -- which is why
revoking or rotating it revokes access to everything downstream.

## Why a First-Class Resource?

- **Many-per-vault, independent lifecycle** -- keys are created, rotated,
  and retired on their own schedules inside a long-lived vault
- **FK-referenced by CMK consumers** -- `AzureContainerRegistry`
  encryption and `AzureAksCluster` KMS reference the key's outputs
  directly; anything referenceable must be a first-class node
- **Rotation is a key-level policy** -- automatic rotation, per-version
  expiry, and near-expiry notification ride the key, not the vault

## Key Features

- **All four algorithm families** -- RSA / RSA_HSM (2048/3072/4096) and
  EC / EC_HSM (P-256/P-256K/P-384/P-521); the HSM variants require the
  vault's PREMIUM SKU
- **Capability boundary** -- `key_opts` lists exactly the operations the
  key may perform (wrap/unwrap for CMK, encrypt/decrypt, sign/verify);
  Azure rejects everything else
- **Automatic rotation policy** -- rotate on a schedule
  (`time_after_creation`) or relative to expiry (`time_before_expiry`),
  stamp per-version expiry (`expire_after`), and raise Event Grid
  notifications before expiry
- **Versionless composition** -- CMK consumers reference the
  `versionless_id` output so rotation propagates automatically; pin
  `key_id` only when a compliance regime demands a frozen version

## Spec Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | 1-127 letters/digits/hyphens, unique among the vault's keys; ForceNew |
| `key_vault_id` | StringValueOrRef | Yes | -- | The vault, by ARM ID (defaults to an AzureKeyVault reference); ForceNew |
| `key_type` | enum | Yes | -- | `RSA` / `RSA_HSM` / `EC` / `EC_HSM`; ForceNew |
| `key_size` | int32 | RSA only | -- | 2048 / 3072 / 4096; required for RSA, forbidden for EC; ForceNew |
| `curve` | enum | No | P-256 | `P_256` / `P_256K` / `P_384` / `P_521`; EC only; ForceNew |
| `key_opts` | enum list | Yes | -- | decrypt / encrypt / sign / unwrapKey / verify / wrapKey |
| `not_before_date` | string | No | -- | RFC 3339 activation instant |
| `expiration_date` | string | No | -- | RFC 3339 expiry; prefer rotation_policy for keys encrypting live data |
| `rotation_policy` | message | No | -- | expire_after + notify_before_expiry (paired) and/or automatic rotation trigger |
| `tags` | map | No | -- | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `key_id` | Versioned data-plane ID (pins consumers; rotation does not follow) |
| `versionless_id` | Versionless data-plane ID -- the reference CMK consumers should use |
| `key_name` | The key's name within the vault |
| `version` | The current version identifier |
| `resource_id` / `resource_versionless_id` | ARM (control-plane) IDs |
| `public_key_pem` / `public_key_openssh` | The public half, for consumers that verify or encrypt outside the vault |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureKeyVaultKey
metadata:
  name: registry-cmk
spec:
  name: registry-cmk
  keyVaultId:
    valueFrom:
      name: platform-vault
  keyType: RSA
  keySize: 2048
  keyOpts:
    - WRAP_KEY
    - UNWRAP_KEY
  rotationPolicy:
    expireAfter: P2Y
    notifyBeforeExpiry: P30D
    automatic:
      timeBeforeExpiry: P60D
```

Wire it to a CMK consumer:

```yaml
spec:
  encryption:
    identityClientId: ...
    keyVaultKeyId:
      valueFrom:
        name: registry-cmk
```

## Lifecycle Notes

- **Key material is immutable by design**: type, size, and curve are fixed
  at creation; changing any of them replaces the key and every consumer
  re-encrypts through the new key on its next unwrap
- **The deploying credential needs data-plane key permissions on the
  vault** ("Key Vault Administrator" or "Key Vault Crypto Officer" RBAC
  role, or key permissions in a legacy access policy) -- subscription
  Owner alone is not enough
- A deleted key's name stays reserved for the vault's soft-delete
  retention window unless purged (the IaC engines purge on destroy by
  default)
- Once `expiration_date` is set, it cannot be fully unset on the
  underlying key even across delete/recreate (Azure restores purged
  names' state)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
