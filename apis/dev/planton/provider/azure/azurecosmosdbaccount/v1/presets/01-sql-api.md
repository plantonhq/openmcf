# Production SQL API Account

This preset creates a two-region SQL (NoSQL) API account with automatic
failover and continuous backup -- the production baseline for document
workloads. Databases and containers are deployed as their own
AzureCosmosdbSqlDatabase / AzureCosmosdbSqlContainer resources
referencing this account.

## When to Use

- The primary document store for an application that must survive a
  regional outage
- Teams that want point-in-time restore instead of scheduled snapshots
- The starting point most SQL-API workloads should remix

## Key Configuration Choices

- **`consistencyLevel: SESSION`** -- read-your-writes within a session;
  the sweet spot between STRONG's latency and EVENTUAL's anarchy
- **Two `geoLocations` + `automaticFailoverEnabled`** -- eastus writes,
  westus2 reads and is promoted automatically on failover; regions can
  be added later in place
- **`backup: CONTINUOUS / CONTINUOUS_30_DAYS`** -- a 30-day
  point-in-time restore window; note the one-way door (Continuous
  cannot go back to Periodic without recreating the account)
- **No throughput here** -- capacity is provisioned on the databases
  and containers; add a `capacity.totalThroughputLimit` to cap the
  account-wide spend

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<globally-unique-account-name>` | 3-50 lowercase letters/digits/hyphens, unique across all of Azure | Your naming convention |
| `my-data-rg` | The AzureResourceGroup's Planton resource name | Your resource-group composition |

## Downstream Wiring

Databases reference the account's ARM id:

```yaml
# On an AzureCosmosdbSqlDatabase
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: my-app-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
```
