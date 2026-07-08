# AzureStorageDataLakeGen2Filesystem

A Data Lake Storage Gen2 filesystem inside an AzureStorageAccount: the
namespace unit of an analytics data lake. Spark, Databricks, Synapse,
and the abfss:// driver address data as
`abfss://{filesystem}@{account}.dfs.core.windows.net/path`, and the
filesystem is the grant boundary for data-plane RBAC and root-path
POSIX ACLs -- the one-filesystem-per-zone (raw, curated, gold) pattern.

## When to Use

Use AzureStorageDataLakeGen2Filesystem when you need:

- **Data-lake zones** -- one filesystem per zone, each carrying its own
  access posture (ACLs at the root, RBAC at the filesystem scope)
- **Analytics workspace storage** -- the container Databricks/Synapse
  mount and Spark jobs write to
- **POSIX access control** -- per-principal rwx entries on the root
  path, with DEFAULT entries that newly created children inherit

The account must have hierarchical namespace enabled (`is_hns_enabled`)
for the POSIX surface; on a flat account, create an
AzureStorageContainer instead.

## Key Configuration

- `storage_account_id` -- the parent HNS account, referenced from an
  AzureStorageAccount's output (fixed at creation)
- `filesystem_name` -- 3-63 lowercase letters/digits/hyphens; the
  container segment of every abfss:// URL (renaming replaces the
  filesystem and everything in it)
- `owner` / `group` -- the root path's owning principal/group (Entra
  object IDs or `$superuser`)
- `aces` -- the root path's POSIX ACL; DEFAULT-scope entries are the
  inheritance template for new children
- `default_encryption_scope` -- sub-account key isolation for just this
  filesystem, referencing an AzureStorageEncryptionScope

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: lake-storage
    fieldPath: status.outputs.storage_account_id
```

The `filesystem_id` output is the ARM container-proxy ID -- the scope
data-plane role assignments (Storage Blob Data Reader/Contributor/
Owner via AzureRoleAssignment) target for filesystem-level access.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
