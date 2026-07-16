# Platform-Managed Encryption Scope

This preset creates an encryption scope with a Microsoft-managed key --
a distinct encryption boundary inside a shared account with ZERO key
management: Azure creates and rotates the key.

## When to Use

- Per-tenant data isolation where the requirement is a distinct key
  boundary, not customer key custody
- Separating data domains (uploads vs exports vs archives) under
  different keys on one account

## Key Configuration Choices

- **`source: MICROSOFT_STORAGE`** -- Azure owns rotation; upgrade a
  scope to a customer-managed key later by switching the source and
  adding a key reference (in place, no recreation)
- **`scopeName`** -- alphanumerics only (hyphens forbidden, unlike
  other storage names); this exact string is what containers pin

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `appdatascope` | 4-63 letters and digits, no hyphens | Your tenancy/data taxonomy |

## Downstream Wiring

Pin a container to this scope:

```yaml
# On an AzureStorageContainer
defaultEncryptionScope:
  valueFrom:
    kind: AzureStorageEncryptionScope
    name: my-tenant-scope
    fieldPath: status.outputs.encryption_scope_name
encryptionScopeOverrideEnabled: false
```
