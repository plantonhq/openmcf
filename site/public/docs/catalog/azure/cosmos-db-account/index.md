---
title: "Cosmos DB Account"
description: "Cosmos DB Account deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbaccount"
---

# Azure Cosmos DB Account

Creates an Azure Cosmos DB account -- the globally distributed, multi-model database account that owns regions, consistency, network posture, encryption, and backup. Databases and containers are their own composable resources (AzureCosmosdbSqlDatabase, AzureCosmosdbSqlContainer, AzureCosmosdbMongoDatabase, AzureCosmosdbMongoCollection) referencing this account.

## What Gets Created

When you deploy an AzureCosmosdbAccount resource, Planton provisions:

- **Cosmos DB Account** -- an `azurerm_cosmosdb_account` with your chosen API (SQL or MongoDB), regions, consistency policy, capabilities, network rules, managed identity, customer-managed-key encryption, backup mode, and tags

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the account in (referenced through `resourceGroup`)
- For customer-managed keys: an **AzureKeyVaultKey** in a purge-protected vault, and an **AzureUserAssignedIdentity** holding get/wrapKey/unwrapKey on it (referenced through `keyVaultKeyId`, `identity`, and `defaultIdentity`)

## Quick Start

Create a file `cosmos-account.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbAccount
metadata:
  name: app-cosmos
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbAccount.app-cosmos
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: my-data-rg
      fieldPath: status.outputs.resource_group_name
  accountName: my-org-app-cosmos
  consistencyPolicy:
    consistencyLevel: SESSION
  geoLocations:
    - location: eastus
      failoverPriority: 0
```

Deploy:

```shell
planton apply -f cosmos-account.yaml
```

The account speaks the SQL (NoSQL) API by default. For MongoDB compatibility set `kind: MONGO_DB`, declare the `ENABLE_MONGO` capability, and pick a `mongoServerVersion`.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `cosmosdb_account_id` | The ARM id database kinds and private endpoints reference |
| `cosmosdb_account_name` | The globally unique DNS label |
| `endpoint` | The document endpoint SDKs connect to |
| `read_endpoints` / `write_endpoints` | Per-region endpoints in failover-priority order |
| `primary_key` / `secondary_key` (+ readonly pair) | The account keys (secret-bearing; rotate via the secondary) |
| `primary_sql_connection_string` (+ 3 variants) | Ready-made SQL-API connection strings (secret-bearing) |
| `primary_mongodb_connection_string` (+ 3 variants) | Ready-made MongoDB connection strings (secret-bearing) |
| `identity_principal_id` | The system-assigned identity's principal for role assignments |

When `localAuthenticationEnabled: false`, the keys and connection strings stop authenticating and data-plane access rides Entra ID.

## Related Resources

- **AzureCosmosdbSqlDatabase** / **AzureCosmosdbSqlContainer** -- the SQL-API data containers referencing this account
- **AzureCosmosdbMongoDatabase** / **AzureCosmosdbMongoCollection** -- the MongoDB-API data containers
- **AzureResourceGroup** -- the account's home
- **AzureKeyVaultKey** + **AzureUserAssignedIdentity** -- the customer-managed-key encryption pair
- **AzurePrivateEndpoint** -- private connectivity to a locked-down account
