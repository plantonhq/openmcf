# AzureStorageObjectReplication

An object replication policy between TWO Azure Storage Accounts:
asynchronous, rule-driven copying of block blobs from containers on a
source account to containers on a destination account -- cross-region
DR, read-local data distribution, and archival fan-out without
application-side copy jobs.

## When to Use

Use AzureStorageObjectReplication when you need:

- **Cross-region DR for blobs** -- replicate to an account in a paired
  region and fail reads over when a region degrades
- **Read-local distribution** -- fan content out to an account close to
  each consumer
- **Archival/offboarding fan-out** -- continuously copy a tenant's
  containers toward an archive account

## Key Configuration

- `source_storage_account_id` / `destination_storage_account_id` --
  the account pair, referenced from AzureStorageAccount outputs (both
  fixed at creation)
- `rules` -- one source container → one destination container each (up
  to 1000), with `copy_blobs_created_after` (OnlyNewObjects /
  Everything / an RFC 3339 instant) and include-prefix filters
- **Account prerequisites** (apply-time, on the ACCOUNTS): blob
  versioning + change feed on the source, blob versioning on the
  destination (`blob_properties` on the account spec) -- which also
  means neither account can be hierarchical-namespace

## Composition

```yaml
rules:
  - sourceContainerName:
      valueFrom:
        kind: AzureStorageContainer
        name: invoices
        fieldPath: status.outputs.container_name
    destinationContainerName:
      valueFrom:
        kind: AzureStorageContainer
        name: invoices-replica
        fieldPath: status.outputs.container_name
```

Azure materializes the one logical policy on BOTH accounts under one
GUID; this kind IS the pair, and the outputs carry both ARM ids plus
the shared `policy_id`.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
