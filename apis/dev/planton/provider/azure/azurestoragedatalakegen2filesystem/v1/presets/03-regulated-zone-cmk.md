# Regulated Zone with Customer-Managed Key

This preset creates a compliance zone whose data encrypts under YOUR
Key Vault key (via an encryption scope), with a deny-by-default root
ACL -- key custody and least-privilege for just the regulated subset of
the lake.

## When to Use

- PII/PHI/financial zones inside an otherwise ordinary lake
- Kill-switch requirements: revoking the vault key makes the zone's
  data unreadable

## Key Configuration Choices

- **`defaultEncryptionScope` is fixed at creation** -- choose the
  scope before data lands; every object that doesn't name its own
  scope encrypts under it
- **The scope (not the filesystem) carries the Key Vault wiring** --
  create an AzureStorageEncryptionScope with `source:
  MICROSOFT_KEY_VAULT` first; the account needs an identity with
  wrap/unwrap on the vault
- **Deny-by-default ACL** -- add qualified USER/GROUP entries (or RBAC
  grants at the `filesystem_id` scope) for exactly the principals the
  regulation allows

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The HNS AzureStorageAccount's Planton resource name | Your lake composition |
| `<encryption-scope-resource-name>` | The AzureStorageEncryptionScope's Planton resource name | Your key-isolation composition |
