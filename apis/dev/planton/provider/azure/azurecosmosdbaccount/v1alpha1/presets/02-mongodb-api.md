# MongoDB API Account

This preset creates a MongoDB-compatible account: existing MongoDB
drivers, tools, and code work unchanged against a fully managed,
globally distributable backend. Databases and collections are deployed
as AzureCosmosdbMongoDatabase / AzureCosmosdbMongoCollection resources
referencing this account.

## When to Use

- Migrating a MongoDB application to a managed service without driver
  changes
- New workloads whose teams know the MongoDB query language and
  tooling
- Consuming the account's `primary_mongodb_connection_string` output
  directly in application configuration

## Key Configuration Choices

- **`kind: MONGO_DB`** -- fixed at creation; the wire protocol shapes
  storage
- **`capabilities: [ENABLE_MONGO, ENABLE_MONGO_RETRYABLE_WRITES]`** --
  the API capability is declared explicitly, and retryable writes are
  the driver behavior modern MongoDB applications expect (also one of
  the two capabilities Azure can remove in place)
- **`mongoServerVersion: MONGO_7_0`** -- the newest wire protocol;
  drivers must be compatible with the declared version

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-catalog-mongo` | 3-50 lowercase letters/digits/hyphens, unique across all of Azure | Your naming convention |
| `my-data-rg` | The AzureResourceGroup's Planton resource name | Your resource-group composition |

## Downstream Wiring

Mongo databases reference the account's ARM id:

```yaml
# On an AzureCosmosdbMongoDatabase
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: my-mongo-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
```
