# Serverless, Entra-Only Account

This preset creates a serverless SQL-API account with the fully
locked-down access posture: pay-per-request billing, no public
endpoint, and no key-based authentication -- every data-plane caller
rides Entra ID.

## When to Use

- Development and staging environments with idle periods (serverless
  bills nothing when nothing runs)
- Spiky or unpredictable workloads where provisioned RU/s would sit
  idle
- Regulated workloads that must not expose a public endpoint or share
  account keys

## Key Configuration Choices

- **`capabilities: [ENABLE_SERVERLESS]`** -- pay-per-request; databases
  and containers referencing this account must NOT declare throughput
  (Azure rejects it). Serverless accounts are single-region
- **`publicNetworkAccessEnabled: false`** -- reachable only through an
  AzurePrivateEndpoint referencing `cosmosdb_account_id`
- **`localAuthenticationEnabled: false`** -- the exported keys and
  connection strings stop authenticating; applications use Entra ID
  and Cosmos DB data-plane RBAC
- **`accessKeyMetadataWritesEnabled: false`** -- database/container
  management is restricted to ARM (Entra-authenticated) callers

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<globally-unique-account-name>` | 3-50 lowercase letters/digits/hyphens, unique across all of Azure | Your naming convention |
| `my-data-rg` | The AzureResourceGroup's Planton resource name | Your resource-group composition |

## Downstream Wiring

Private connectivity to the locked-down account:

```yaml
# On an AzurePrivateEndpoint
privateConnectionResourceId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: my-serverless-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
subresourceNames:
  - Sql
```
