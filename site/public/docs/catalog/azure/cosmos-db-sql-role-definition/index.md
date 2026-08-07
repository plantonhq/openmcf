---
title: "Cosmos DB SQL Role Definition"
description: "Cosmos DB SQL Role Definition deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbsqlroledefinition"
---

# Azure Cosmos DB SQL Role Definition

Authors a custom DATA-PLANE role for a Cosmos DB SQL (NoSQL) API account. Cosmos DB carries its own RBAC system, separate from ARM RBAC: an ARM role on the account (even Owner) governs managing it and grants no ability to read or write the documents inside. This kind names a reusable bundle of ALLOW data actions; an AzureCosmosdbSqlRoleAssignment binds it to a principal at a scope. Together with the account's local-authentication switch, this is Cosmos DB's keyless story — keys off, workload identities in, every grant explicit and auditable.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cosmos DB SQL Role Definition** -- a named, GUID-identified custom role inside the referenced account, carrying its assignable scopes and permission blocks

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A Cosmos DB account** speaking the SQL (NoSQL) API. Reference an AzureCosmosdbAccount Cloud Resource via ValueFromRef, or provide the account's ARM ID directly.

### Check the Built-ins First

Every SQL-API account already carries two built-in data roles — Data Reader (`00000000-0000-0000-0000-000000000001`) and Data Contributor (`00000000-0000-0000-0000-000000000002`). Assigning one needs NO definition resource: an AzureCosmosdbSqlRoleAssignment references the well-known ID directly. Author a custom definition when neither fits — "read-only on one container", "write items but never delete", "metadata-only for a monitoring probe".

## Deploy

### Console

Open the deployment store, find **Azure Cosmos DB SQL Role Definition**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **read-only-role** preset in the [Presets](#presets) tab for the reader shape.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbSqlRoleDefinition
metadata:
  name: orders-read-only
  org: acme-corp
  env: prod
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: orders-db
      fieldPath: status.outputs.cosmosdb_account_id
  roleName: orders-read-only
  assignableScopes:
    - valueFrom:
        kind: AzureCosmosdbAccount
        name: orders-db
        fieldPath: status.outputs.cosmosdb_account_id
  permissions:
    - dataActions:
        - Microsoft.DocumentDB/databaseAccounts/readMetadata
        - Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read
        - Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/executeQuery
        - Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/readChangeFeed
```

```shell
planton apply -f cosmosdb-sql-role-definition.yaml
```

This creates a read-only role assignable anywhere in the account. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the account, its roles, and their grants compose in one InfraPipeline: the pipeline deploys the account first, resolves `cosmosdb_account_id` into this definition, then resolves this definition's `role_definition_id` into its assignments.

## Key Configuration

**Permissions** -- ALLOW rules over Cosmos data operations, evaluated as a union across blocks. No carve-outs and no control-plane actions exist in this RBAC system; the four write operations (create/replace/upsert/delete) grant individually. Practically every role includes `readMetadata` — SDKs read metadata before any data call. Editing permissions later updates EVERY existing assignment of the role at once.

**Assignable scopes** -- Where assignments of this role may be CREATED: the account, a database (`{account-id}/dbs/{db}`), or a container (`{account-id}/dbs/{db}/colls/{container}`). At least one is required; database/container entries are literal paths composed on the account ID (references cannot append path suffixes).

**Role name and GUID** -- The display name is what audits read and renames in place; assignments track the role by GUID, generated at deploy unless pinned via `roleDefinitionId`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| AzureCosmosdbAccount | `cosmosdbAccountId` | `status.outputs.cosmosdb_account_id` |
| AzureCosmosdbAccount | `assignableScopes[]` | `status.outputs.cosmosdb_account_id` |

### What This Component Produces

| Output | Description | Consumed By |
|--------|-------------|-------------|
| `role_definition_id` | The fully-scoped ARM ID of the definition | AzureCosmosdbSqlRoleAssignment |
| `role_definition_guid` | The definition's GUID resource name | Audit tooling |
| `role_name` | The display name as deployed | Data-plane RBAC listings |
| `cosmosdb_account_name` | The parent account's name | Automation that needs the account/definition pair |
