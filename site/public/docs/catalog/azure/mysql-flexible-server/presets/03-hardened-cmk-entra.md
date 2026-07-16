---
title: "Hardened Server with CMK Encryption and an Entra Administrator"
description: "This preset creates a compliance-oriented MySQL Flexible Server: customer-managed-key (CMK) encryption unwrapped through a user-assigned identity, a Microsoft Entra administrator for token-based..."
type: "preset"
rank: "03"
presetSlug: "03-hardened-cmk-entra"
componentSlug: "mysql-flexible-server"
componentTitle: "MySQL Flexible Server"
provider: "azure"
icon: "package"
order: 3
---

# Hardened Server with CMK Encryption and an Entra Administrator

This preset creates a compliance-oriented MySQL Flexible Server:
customer-managed-key (CMK) encryption unwrapped through a user-assigned
identity, a Microsoft Entra administrator for token-based administration,
35-day backups, TLS-only connections, and the audit log enabled.

## When to Use

- Regulated environments where the encryption key must live in your own
  Key Vault with your own rotation policy
- Organizations standardizing database administration on Entra identities
  and groups instead of shared SQL passwords

## Key Configuration Choices

- **Password auth stays on** -- unlike PostgreSQL Flexible Server, MySQL
  cannot disable password authentication; the Entra administrator is
  additive. Keep the password in a secret manager and treat Entra tokens
  as the primary access path
- **CMK against `versionless_id`** -- key rotations propagate to the
  server automatically; pin a versioned key ID only when a compliance
  regime demands an immutable version. The key's vault needs purge
  protection enabled
- **One user-assigned identity serves both seams** -- it unwraps the CMK
  (needs wrap/unwrap access on the vault) and backs the Entra
  administrator grant (Azure uses it to read directory objects). Split
  them into two identities if your separation-of-duties model requires it
- **`audit_log_enabled: ON`** -- MySQL's server-side audit trail,
  consumable through Azure Monitor diagnostics

## Prerequisites

- An `AzureUserAssignedIdentity` with wrap/unwrap access on the key's
  vault (a "Key Vault Crypto Service Encryption User" role assignment)
- An `AzureKeyVaultKey` in a vault with purge protection enabled

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `myorg-hardened-mysql` | 3-63 lowercase chars, globally unique | Your naming convention |
| `mysqladmin` | The MySQL admin login | Your convention |
| `<admin-password>` | 8-128 chars from 3+ character classes | A secret manager; never commit literals |
| `<cmk-identity-resource-name>` | The user-assigned identity's Planton resource name | Your identity composition |
| `<cmk-key-resource-name>` | The Key Vault key's Planton resource name | Your Key Vault composition |
| `<entra-group-name>` | The Entra group's display name (the Entra login name) | Microsoft Entra ID -> Groups |
| `<entra-group-object-id>` | The Entra group's directory object ID | Microsoft Entra ID -> Groups -> Overview |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Entra connections authenticate with a directory token instead of the
password (the login name is the Entra principal's display name):

```text
mysql -h {status.outputs.fqdn} -u {entra-group-name} --enable-cleartext-plugin -p{entra-access-token}
```
