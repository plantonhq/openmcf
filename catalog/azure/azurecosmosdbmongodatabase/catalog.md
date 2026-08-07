# Azure Cosmos DB Mongo Database

Deploys a MongoDB API database inside an Azure Cosmos DB account — the namespace collections live in and the boundary for SHARED throughput. Databases are many-per-account with independent lifecycles, which is why they are a first-class kind referencing the account rather than a list folded into the account's spec. Collections are their own kind (AzureCosmosdbMongoCollection) referencing this database's `mongo_database_id` output.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cosmos DB Mongo Database** -- a named database inside the referenced Cosmos DB account (which must be a MONGO_DB-kind account with the ENABLE_MONGO capability)
- **Shared Throughput** (optional) -- fixed RU/s or an autoscale ceiling that every collection in the database shares, when either is declared; omit both to let each collection bring its own dedicated throughput (or on serverless accounts, where provisioned throughput is rejected)

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A Cosmos DB account** speaking the MongoDB API (kind MONGO_DB with the ENABLE_MONGO capability). Reference an AzureCosmosdbAccount Cloud Resource via ValueFromRef, or provide the account's ARM ID directly.

## Deploy

### Console

Open the deployment store, find **Azure Cosmos DB Mongo Database**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **dedicated-collection-throughput** preset in the [Presets](#presets) tab for the common production shape.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbMongoDatabase
metadata:
  name: catalog-database
  org: acme-corp
  env: prod
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: catalog-db
      fieldPath: status.outputs.cosmosdb_account_id
  databaseName: catalog
```

```shell
planton apply -f cosmosdb-mongo-database.yaml
```

This creates a database with no shared throughput — each collection brings its own. Add `throughput` (fixed RU/s) or `autoscaleMaxThroughput` (a ceiling) for the shared model. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the account, database, and collections compose in one InfraPipeline: the pipeline deploys the account first, resolves `cosmosdb_account_id` into this database, then resolves this database's `mongo_database_id` into its collections.

## Key Configuration

**Throughput model** -- The single real decision. Leave both fields unset (the production norm) and each collection provisions its own dedicated throughput; set `throughput` for a fixed shared budget (minimum 400 RU/s, increments of 100); or set `autoscaleMaxThroughput` (minimum 1000 RU/s, increments of 1000) and Azure floats the shared budget between 10% and 100% of it. The two are mutually exclusive. On serverless accounts (the ENABLE_SERVERLESS capability) neither may be set — Azure rejects provisioned throughput at apply.

**Database name** -- Unique within the account, 1-255 characters (Azure's only constraint), and what Mongo drivers address after connecting to the account. Changing the name replaces the database AND everything in it.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| AzureCosmosdbAccount | `cosmosdbAccountId` | `status.outputs.cosmosdb_account_id` |

### What This Component Produces

| Output | Description | Consumed By |
|--------|-------------|-------------|
| `mongo_database_id` | The ARM ID of the database | AzureCosmosdbMongoCollection |
| `mongo_database_name` | The name Mongo drivers reference | Application configuration |
| `cosmosdb_account_name` | The parent account's name | Connection string composition |
