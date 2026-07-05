---
title: "Customer-Managed-Key Encryption Scope"
description: "This preset creates an encryption scope backed by YOUR Key Vault key -- customer key custody for just the data that opts into the scope, while the rest of the account stays on platform-managed keys."
type: "preset"
rank: "02"
presetSlug: "02-customer-managed-key-scope"
componentSlug: "storage-encryption-scope"
componentTitle: "Storage Encryption Scope"
provider: "azure"
icon: "package"
order: 2
---

# Customer-Managed-Key Encryption Scope

This preset creates an encryption scope backed by YOUR Key Vault key --
customer key custody for just the data that opts into the scope, while
the rest of the account stays on platform-managed keys.

## When to Use

- Regulated data subsets (PII, PHI, financial records) inside an
  otherwise ordinary account
- Tenants who contractually require holding their own encryption key
- Kill-switch requirements: revoking the vault key makes the scope's
  data unreadable

## Key Configuration Choices

- **`keyVaultKeyId` references the VERSIONLESS id** -- key rotation
  propagates to the scope automatically; pin a versioned URI only when
  a compliance regime demands a frozen version
- **The ACCOUNT needs identity plumbing** -- attach a managed identity
  to the parent account and grant it wrap/unwrap on the key's vault
  (RBAC: Key Vault Crypto Service Encryption User), exactly as for the
  account-level customer-managed key
- The spec rejects `MICROSOFT_KEY_VAULT` without a key reference at
  validation time -- the pairing cannot ship half-configured

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<scopeName>` | 4-63 letters and digits, no hyphens | Your tenancy/data taxonomy |
| `<key-vault-key-resource-name>` | The AzureKeyVaultKey's Planton resource name | Your key composition |
