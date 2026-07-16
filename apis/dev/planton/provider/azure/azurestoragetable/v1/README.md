# AzureStorageTable

A Storage table inside an AzureStorageAccount: the serverless NoSQL
key/value store of Azure storage. Applications store schemaless entities
addressed by partition key + row key -- device state, user profiles,
audit trails, IoT telemetry -- at petabyte scale with
single-digit-millisecond point reads and no capacity provisioning.

## When to Use

Use AzureStorageTable when you need:

- **Cheap, huge key/value storage** -- entities addressed by
  partition + row key, no throughput provisioning
- **Audit trails and telemetry** -- append-heavy, point-read workloads
- **A role-assignment scope** -- grant Storage Table Data
  Reader/Contributor on `table_id` for table-level access

Cosmos DB's Table API is the premium sibling (global distribution,
throughput SLAs, secondary indexes) at a very different price point.

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation); the account must
  keep shared keys enabled -- the provider drives table create/ACLs
  through the data plane with shared-key auth
- `table_name` -- 3-63 alphanumerics starting with a letter (no
  hyphens; never the literal word "table"), unique within the account
- `acls` -- stored access policies anchoring revocable SAS tokens
  (table policies require the full validity window)

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Client URLs compose from the ACCOUNT's endpoint plus this table's name:
`{primary_table_endpoint}{table_name}`.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
