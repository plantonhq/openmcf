# AzureCosmosdbMongoDatabase

A MongoDB API database inside an AzureCosmosdbAccount: the namespace
collections live in and the boundary for SHARED throughput. A database
either provisions throughput (fixed RU/s or autoscale) that every
collection without its own shares, or provisions nothing and lets each
collection bring dedicated throughput -- the common production shape,
because shared throughput couples noisy neighbors.

## When to Use

Use AzureCosmosdbMongoDatabase when you need:

- **A collection namespace** -- every AzureCosmosdbMongoCollection
  lives in exactly one database and references it
- **A shared-throughput pool** -- one RU/s budget (fixed or autoscale)
  split across many small collections that would each waste a
  dedicated minimum
- **A pure namespace with no throughput** -- collections provision
  their own dedicated RU/s, isolated from each other; also the only
  legal shape on serverless accounts

## Key Configuration

- `cosmosdb_account_id` -- the parent account, referenced from an
  AzureCosmosdbAccount's output (fixed at creation); the account must
  be a MONGO_DB-kind account with the ENABLE_MONGO capability
- `database_name` -- 1-255 characters, unique within the account;
  renaming replaces the database AND everything in it
- `throughput` -- fixed shared RU/s (minimum 400, increments of 100);
  mutually exclusive with autoscale
- `autoscale_max_throughput` -- autoscale ceiling (minimum 1000,
  increments of 1000); Azure scales between 10% and 100% of it
- Leave BOTH throughput fields unset when collections bring their own
  RU/s, and always on serverless accounts (Azure rejects provisioned
  throughput on serverless at apply)

## Composition

```yaml
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: app-cosmos-mongo
    fieldPath: status.outputs.cosmosdb_account_id
```

Collections reference this database the same way, through its
`mongo_database_id` output. Connectivity lives on the ACCOUNT --
applications combine the account's MongoDB connection-string outputs
with this database's `mongo_database_name`.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
