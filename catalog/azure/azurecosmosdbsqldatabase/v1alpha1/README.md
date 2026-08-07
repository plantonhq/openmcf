# AzureCosmosdbSqlDatabase

A SQL (NoSQL) API database inside an AzureCosmosdbAccount: the
namespace containers live in and the boundary for SHARED throughput.
A database either provisions throughput (fixed RU/s or autoscale) that
every container without its own shares, or provisions nothing and lets
each container bring dedicated throughput -- the common production
shape, because shared throughput couples noisy neighbors.

## When to Use

Use AzureCosmosdbSqlDatabase when you need:

- **A container namespace** -- every AzureCosmosdbSqlContainer lives in
  exactly one database and references it
- **A shared-throughput pool** -- one RU/s budget (fixed or autoscale)
  split across many small containers that would each waste a dedicated
  minimum
- **A pure namespace with no throughput** -- containers provision their
  own dedicated RU/s, isolated from each other; also the only legal
  shape on serverless accounts

## Key Configuration

- `cosmosdb_account_id` -- the parent account, referenced from an
  AzureCosmosdbAccount's output (fixed at creation); the account must
  be a GLOBAL_DOCUMENT_DB (SQL API) account
- `database_name` -- 1-255 characters, unique within the account;
  renaming replaces the database AND everything in it
- `throughput` -- fixed shared RU/s (minimum 400, increments of 100);
  mutually exclusive with autoscale
- `autoscale_max_throughput` -- autoscale ceiling (minimum 1000,
  increments of 1000); Azure scales between 10% and 100% of it
- Leave BOTH throughput fields unset when containers bring their own
  RU/s, and always on serverless accounts (Azure rejects provisioned
  throughput on serverless at apply)

## Composition

```yaml
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: app-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
```

Containers reference this database the same way, through its
`sql_database_id` output. Connectivity and keys live on the ACCOUNT --
applications combine the account's endpoint or connection-string
outputs with this database's `sql_database_name`.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
