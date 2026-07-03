# Data-Plane Role (Blob Auditor)

This preset teaches the distinction that trips up most custom-role authors:
control-plane `actions` manage Azure resources, data-plane `dataActions`
access the data inside them. A role with every Storage action but no
dataActions can delete a whole storage account yet cannot read one byte of
blob data.

The auditor shape combines minimal control-plane reads (so the principal can
discover accounts and containers in the portal or CLI) with a single
data-plane read (the blob contents themselves). The same pattern applies to
queues, tables, Key Vault (RBAC mode), Service Bus, and Event Hubs -- any
service with a data plane.

## When to Use

- Audit and compliance principals that must read data without managing anything
- Analytics identities reading blobs across many accounts (assign once at
  subscription scope instead of per-account built-in grants)
- As the template for any custom data-plane role (swap the operation patterns)

## Key Configuration Choices

- **Pair control-plane discovery with data-plane access** -- data-plane-only
  roles work for programmatic access by ID, but principals cannot browse to
  the data in the portal without the control-plane reads
- **Carve out with `notDataActions`** -- e.g. add
  `Microsoft.Storage/storageAccounts/blobServices/containers/blobs/delete`
  under a wildcard `dataActions` grant for write-but-never-delete roles

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<org-prefix>` | Your organization's role-name prefix (names are tenant-unique) | Your naming convention |
| `<subscription-arm-id>` | `/subscriptions/{subscription-id}` | `az account show --query id` (prepend `/subscriptions/`) |
