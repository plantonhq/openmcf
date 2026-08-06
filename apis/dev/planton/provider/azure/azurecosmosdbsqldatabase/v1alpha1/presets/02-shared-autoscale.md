# Shared Autoscale Database

This preset creates a database with an autoscale throughput budget that
every container inside it SHARES (unless a container brings its own).
Azure scales the budget between 10% and 100% of the ceiling as traffic
moves -- the economical shape for many small containers whose combined
traffic never justifies per-container dedicated throughput.

## When to Use

- Many small containers (reference data, settings, feature flags,
  per-tenant metadata) that would each waste a dedicated 400 RU/s
  minimum
- Spiky combined traffic that a fixed shared budget would either
  over-provision or throttle
- Development and staging environments consolidating cost

## Key Configuration Choices

- **`autoscaleMaxThroughput: 4000`** -- the ceiling; Azure scales the
  shared budget between 400 and 4000 RU/s. Minimum ceiling 1000, in
  increments of 1000. The 10% floor always bills, so size the ceiling
  to real peaks
- **Shared means coupled** -- one hot container can consume the whole
  budget and starve its siblings; keep genuinely hot containers OUT of
  the pool by giving them their own `throughput` or
  `autoscaleMaxThroughput` on the container
- **XOR with fixed throughput** -- `throughput` and
  `autoscaleMaxThroughput` are mutually exclusive; the spec rejects
  both together before anything reaches Azure

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name (a GLOBAL_DOCUMENT_DB account) | Your Cosmos composition |
| `<database-name>` | 1-255 characters, unique within the account | Your naming convention |

## Downstream Wiring

Containers that should share this budget simply omit their own
throughput fields:

```yaml
# On an AzureCosmosdbSqlContainer -- no throughput fields: it shares
# the database's autoscale budget.
sqlDatabaseId:
  valueFrom:
    kind: AzureCosmosdbSqlDatabase
    name: my-shared-database
    fieldPath: status.outputs.sql_database_id
```
