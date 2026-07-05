# AzureStorageEncryptionScope

An encryption scope inside an AzureStorageAccount: a named encryption
boundary that lets different data in ONE account encrypt under DIFFERENT
keys. Containers and blobs opt in by name -- per-tenant keys inside a
shared account, or a customer-managed key for just the regulated subset,
without splitting into per-tenant accounts.

## When to Use

Use AzureStorageEncryptionScope when you need:

- **Per-tenant key isolation** -- one account, one scope per tenant,
  each container pinned to its tenant's scope
- **Mixed-sensitivity accounts** -- platform-managed keys for ordinary
  data, a Key Vault key for the regulated subset
- **A referenceable key boundary** -- containers
  (`default_encryption_scope`) and ADLS filesystems reference scopes by
  name

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation)
- `scope_name` -- 4-63 letters and digits (no hyphens), unique within
  the account; what containers and blobs reference
- `source` -- MICROSOFT_STORAGE (platform-managed key) or
  MICROSOFT_KEY_VAULT (your key; requires `key_vault_key_id` and an
  account identity with wrap/unwrap access on the vault)
- `infrastructure_encryption_required` -- double-encrypt just this
  scope's data (fixed at creation)

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Containers pin the scope through their `default_encryption_scope`
reference to this kind's `encryption_scope_name` output.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
