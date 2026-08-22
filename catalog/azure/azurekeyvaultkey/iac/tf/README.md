# AzureKeyVaultKey Terraform Module

## Overview

This Terraform module provisions a cryptographic key inside an Azure Key
Vault using the `azurerm` provider. It creates a single
`azurerm_key_vault_key` -- a DATA-PLANE object: the provider talks to the
vault's `{name}.vault.azure.net` endpoint, not ARM, so creation fails with
a 403 when the deploying credential lacks key permissions on the vault,
even if it owns the subscription.

Key material is immutable by design: type, size, and curve are fixed at
creation; changing any of them replaces the key and every consumer
re-encrypts through the new key on its next unwrap.

## Resources Created

- `azurerm_key_vault_key.main` -- the key, with its optional rotation
  policy (Azure models the policy as a sub-resource updated in place)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Key specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | 1-127 letters/digits/hyphens, unique among the vault's keys; ForceNew |
| `key_vault_id` | yes | The vault's ARM ID (resolved from an AzureKeyVault reference); ForceNew |
| `key_type` | yes | `RSA`/`RSA_HSM`/`EC`/`EC_HSM`, mapped to Azure's hyphenated wire values; `_HSM` needs a Premium vault; ForceNew |
| `key_size` | RSA only | 2048/3072/4096 (spec validation enforces the RSA/EC pairing); ForceNew |
| `curve` | no | EC curve; unset lets Azure default to P-256; ForceNew |
| `key_opts` | yes | The capability boundary, mapped to Azure's camelCase operation strings |
| `not_before_date` / `expiration_date` | no | RFC 3339 activation/expiry |
| `rotation_policy` | no | expire_after + notify_before_expiry (paired) and/or the automatic rotation trigger |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `key_id` | Versioned data-plane ID (pins consumers) |
| `versionless_id` | Versionless data-plane ID -- the CMK reference that follows rotation |
| `key_name` / `version` | Name and current version |
| `resource_id` / `resource_versionless_id` | ARM (control-plane) IDs |
| `public_key_pem` / `public_key_openssh` | The public half |

## Operational Notes

- A deleted key's name stays reserved in the vault for the soft-delete
  retention window; the provider's default features purge soft-deleted
  keys on destroy so the name frees up immediately.
- Once `expiration_date` is set it cannot be fully unset on the underlying
  key, even across delete/recreate (Azure restores purged names' state).

## Required Permissions

The deployer permissions this module needs are cataloged in
[`../permissions.yaml`](../permissions.yaml).
