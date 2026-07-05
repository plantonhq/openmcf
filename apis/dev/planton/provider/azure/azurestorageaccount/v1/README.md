# AzureStorageAccount

An Azure Storage Account: the multi-service storage primitive fronting
Blob (objects), Files (SMB/NFS shares), Queues, Tables, and Data Lake
Storage Gen2 behind one globally-unique DNS name. Kind, performance
tier, and replication pick the SKU; service-level blocks tune the data
services; blob containers are first-class `AzureStorageContainer`
resources referencing this account.

## When to Use

Use AzureStorageAccount when you need:

- **Object storage** -- application assets, uploads, backups, data-lake
  raw zones (add `is_hns_enabled` for ADLS Gen2 analytics)
- **The storage backend an app service requires** -- Function Apps and
  Web Apps bind to the account's name and access key
- **File shares, queues, or tables** -- one account fronts all of them
- **A static website or CDN origin** -- `static_website` plus the
  `primary_web_host` / `primary_blob_host` outputs

## Key Configuration

### The SKU trio

- `account_kind` -- STORAGE_V2 (default, right for virtually everything),
  BLOCK_BLOB_STORAGE / FILE_STORAGE (premium specializations),
  BLOB_STORAGE / STORAGE (legacy)
- `account_tier` -- STANDARD (default) or PREMIUM (SSD, pairs with the
  specialized kinds); fixed at creation
- `replication_type` -- LRS (default) → ZRS (zone loss) → GZRS (regional
  loss) → RA_* (adds a read-only paired-region endpoint, lighting up the
  secondary endpoint outputs)

### Security posture

- `https_traffic_only_enabled` + `min_tls_version` -- transport floor
  (both default to Azure's secure side)
- `shared_access_key_enabled: false` -- Entra-only data plane (verify
  every consumer first)
- `allow_nested_items_to_be_public: false` -- makes anonymous container
  access unrepresentable account-wide
- `network_rules` -- DENY + IP rules / subnet references /
  trusted-service bypass; `public_network_access_enabled: false` removes
  the public endpoint entirely
- `customer_managed_key` + `identity` -- bring-your-own-key encryption
  against an `AzureKeyVaultKey`; `infrastructure_encryption_enabled`
  double-encrypts at rest

### Data protection

- `blob_properties` -- versioning, blob + container soft delete, change
  feed, point-in-time restore, CORS, last-access tracking
- `immutability_policy` -- account-wide WORM (requires versioning; LOCKED
  is irreversible)
- `lifecycle_rules` -- tier down (cool/cold/archive) and delete on
  age/access schedules, filtered by prefix, blob type, and index tags

## Composition

Reference the account's outputs from other resources:

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Key outputs: `storage_account_id` (containers, role-assignment scopes),
`storage_account_name` + `primary_access_key` (app-service bindings),
per-service endpoints and hosts (applications, CDN origins).

## Documentation

- [Design research](docs/README.md) -- field mapping, decomposition
  decisions, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
