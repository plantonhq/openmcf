---
title: "Double-Encryption Scope"
description: "This preset creates a scope with infrastructure (double) encryption: two independent encryption layers with independent keys and algorithms for the data that opts in -- WITHOUT enabling account-wide..."
type: "preset"
rank: "03"
presetSlug: "03-double-encryption-scope"
componentSlug: "storage-encryption-scope"
componentTitle: "Storage Encryption Scope"
provider: "azure"
icon: "package"
order: 3
---

# Double-Encryption Scope

This preset creates a scope with infrastructure (double) encryption:
two independent encryption layers with independent keys and algorithms
for the data that opts in -- WITHOUT enabling account-wide
infrastructure encryption (which is creation-time-only on the account).

## When to Use

- Compliance regimes that mandate defense against a single-algorithm or
  single-key compromise for specific data classes
- Adding double encryption to a subset of data on an EXISTING account
  that was created without account-level infrastructure encryption

## Key Configuration Choices

- **`infrastructureEncryptionRequired: true`** -- fixed at creation;
  the second layer always uses platform-managed keys regardless of the
  scope's source
- **Independent of the account switch** -- the account's own
  `infrastructureEncryptionEnabled` can stay off; this scope carries
  its own second layer
- Combine with `source: MICROSOFT_KEY_VAULT` when the regime also
  requires customer key custody (see the customer-managed-key preset)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<scopeName>` | 4-63 letters and digits, no hyphens | Your data classification |
