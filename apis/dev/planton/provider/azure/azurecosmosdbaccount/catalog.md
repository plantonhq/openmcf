# Azure Cosmos DB Account

Deploys an Azure Cosmos DB account — the globally distributed, multi-model database account that owns regions, consistency, network posture, encryption, and backup for everything stored inside it. Databases and containers are first-class kinds (AzureCosmosdbSqlDatabase / AzureCosmosdbSqlContainer for the SQL API, AzureCosmosdbMongoDatabase / AzureCosmosdbMongoCollection for MongoDB) that reference this account's `cosmosdb_account_id` output — nothing is embedded in this spec. The account integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cosmos DB Account** -- a globally distributed database account in the specified Azure region and resource group, with a globally-unique DNS endpoint (`https://{accountName}.documents.azure.com`), the chosen API kind (SQL/NoSQL or MongoDB), consistency policy, geo-locations, capabilities, backup policy, and network access controls
- **Managed Identity** -- created when `identity` is set; the account's own Entra identity (system-assigned, user-assigned, or both) — what unwraps a customer-managed encryption key
- **Virtual Network Rules** -- created when `virtualNetworkRules` entries are provided; restricts account access to specified subnets carrying the Microsoft.AzureCosmosDB service endpoint
- **Azure Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the account for tracking and governance

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **An Azure Resource Group** where the Cosmos DB account will be created. Provide the name directly or reference an AzureResourceGroup Cloud Resource via ValueFromRef.
- **A globally unique account name** -- `accountName` becomes the endpoint hostname (`https://{accountName}.documents.azure.com`). Must be 3-50 characters, lowercase letters, numbers, and hyphens only.
- **Subnets with service endpoints** (optional) -- required when using virtual network filtering. Each subnet must have the `Microsoft.AzureCosmosDB` service endpoint enabled.
- **A Key Vault key** (optional) -- when using customer-managed encryption, reference an AzureKeyVaultKey's `versionless_id` output. The vault must have purge protection enabled.

## Deploy

### Console

Open the deployment store, find **Azure Cosmos DB Account**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **SQL API** preset in the [Presets](#presets) tab for the most common configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbAccount
metadata:
  name: orders-db
  org: acme-corp
  env: prod
spec:
  region: eastus
  resourceGroup:
    value: "acme-prod-rg"
  accountName: acme-orders-prod
  consistencyPolicy:
    consistencyLevel: SESSION
  geoLocations:
    - location: eastus
      failoverPriority: 0
```

```shell
planton apply -f cosmosdb-account.yaml
```

This creates a Cosmos DB account with the SQL API (default), Session consistency, a single geo-location, and Azure's defaults everywhere else. Databases attach afterwards as AzureCosmosdbSqlDatabase resources referencing this account. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Cosmos DB account to a resource group deployed in the same InfraPipeline:

```yaml
spec:
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: production-rg
      fieldPath: status.outputs.resource_group_name
```

The InfraPipeline resolves the dependency graph, deploys the resource group first, then provisions the Cosmos DB account with the resolved value.

## Key Configuration

These are the most important decisions when configuring a Cosmos DB account. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**API kind** -- The `kind` field selects the wire protocol: unset / `GLOBAL_DOCUMENT_DB` for SQL-like queries over JSON documents (and the base for Cassandra, Gremlin, and Table via capabilities), or `MONGO_DB` for MongoDB driver compatibility. This is a ForceNew field — changing it destroys and recreates the account.

**Consistency** -- Five well-defined levels from Strong (linearizable) to Eventual. Session — read-your-writes within a session — is Azure's recommendation for most applications and the default when unset.

**Global distribution** -- `geoLocations` declares the regions data lives in. Priority 0 is the write region; higher numbers promote on failover. Adding and removing regions is an in-place update. Set `automaticFailoverEnabled` for unattended promotion, or `multipleWriteLocationsEnabled` for active-active writes (with per-container conflict resolution).

**Capabilities** -- Customize what the account can do: serverless billing (`ENABLE_SERVERLESS`), extra APIs on a SQL account, MongoDB feature switches, vector/full-text search. Most capability changes recreate the account — settle them before production.

**Network posture** -- `publicNetworkAccessEnabled: false` restricts all access to private endpoints (AzurePrivateEndpoint targeting `cosmosdb_account_id`, subresource `"Sql"` or `"MongoDB"`). When public, `isVirtualNetworkFilterEnabled` + `virtualNetworkRules` and `ipRangeFilter` scope who can reach the account.

**Customer-managed encryption** -- `keyVaultKeyId` points encryption at a Key Vault key you own (the VERSIONLESS key id so rotation propagates). Fixed at creation. Pair with a user-assigned `defaultIdentity` so the key can be unwrapped before the account's own identity exists.

**Backup** -- Unset means Periodic (every 4 hours, retained 8 hours, geo-redundant). Continuous enables point-in-time restore and is required for `createMode: RESTORE`. Periodic → Continuous upgrades in place; Continuous → Periodic recreates the account.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| AzureResourceGroup | `resourceGroup` | `status.outputs.resource_group_name` |
| AzureSubnet | `virtualNetworkRules[].subnetId` | `status.outputs.subnet_id` |
| AzureKeyVaultKey | `keyVaultKeyId` | `status.outputs.versionless_id` |
| AzureUserAssignedIdentity | `identity.identityIds[]`, `defaultIdentity.userAssignedIdentityId` | `status.outputs.identity_id` |

### What This Component Produces

| Output | Description | Consumed By |
|--------|-------------|-------------|
| `cosmosdb_account_id` | The ARM ID of the account | AzureCosmosdbSqlDatabase, AzureCosmosdbMongoDatabase, AzurePrivateEndpoint |
| `cosmosdb_account_name` | The globally unique DNS label | Connection string composition |
| `endpoint` | The document endpoint SDKs connect to | Application configuration |
| `read_endpoints` / `write_endpoints` | Per-region endpoints | Latency-sensitive readers / multi-write writers |
| `primary_key` / `secondary_key` / readonly variants | Account keys (secret-bearing) | Legacy key-auth clients — prefer Entra RBAC |
| `primary_sql_connection_string` (+ 3 variants) | Ready-made SQL-API connection strings (secret-bearing) | SQL SDK clients |
| `primary_mongodb_connection_string` (+ 3 variants) | Ready-made MongoDB connection strings (secret-bearing) | MongoDB drivers on MongoDB accounts |
| `identity_principal_id` | System-assigned identity principal | Role assignments the account needs |

When `localAuthenticationEnabled` is false, the keys and connection strings stop authenticating — data-plane access rides Entra ID instead.
