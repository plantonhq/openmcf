# Azure Data Factory Linked Service

Deploys one linked service inside an Azure Data Factory -- a saved connection telling pipelines, datasets, and data flows where an external system lives and how to authenticate to it, in any of 23 connection types plus a raw-JSON custom form. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions exactly one linked service of the type the spec's variant block declares:

- **Storage** -- Azure Blob Storage, Azure Files, Table Storage, Data Lake Storage Gen2
- **Databases** -- Azure SQL Database, SQL Server, SQL Managed Instance, Synapse, PostgreSQL, MySQL, Cosmos DB (SQL API), Cosmos DB for MongoDB, Snowflake, Kusto, ODBC
- **Services** -- Key Vault, Azure AI Search, Azure Function, Databricks
- **Protocol endpoints** -- web (HTTP), OData, SFTP
- **Custom** -- any other Data Factory connector, as its ARM type name plus type-properties JSON

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Data Factory** -- the linked service lives in a factory; reference an AzureDataFactory's ID output (or provide a literal).
- **For Key-Vault-sourced secrets**: deploy the KEY VAULT variant of this same kind first -- other connections' secret blocks reference it by name.

### Azure Subscription

- **Creating a connection does not test it** -- Azure saves the definition without dialing the target; a wrong password or unreachable host surfaces when a pipeline first USES the connection. Use Studio's Test connection button after deploy.
- **Managed identity needs data-plane permission**: grant the factory's identity the right role on the target (e.g. Storage Blob Data Contributor, Key Vault get/list on secrets) -- the connection saves fine without it and fails at run time.
- **Private networks need a self-hosted integration runtime** -- set `integrationRuntimeName` to reach systems the shared Azure runtime cannot.

## Deploy

### Console

Open the deployment store, find **Azure Data Factory Linked Service**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Key Vault Connection** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f data-factory-linked-service.yaml
```

## After Deploy

The connection appears in the factory's Studio under Manage -> Linked services. Use **Test connection** there to verify reachability and credentials before pointing datasets at it. Secrets never read back through ARM: what you see in Studio is masked, and imports carry the address, not the credential.
