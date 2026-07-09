---
title: "Cosmos DB SQL Role Assignment"
description: "Cosmos DB SQL Role Assignment deployment documentation"
icon: "package"
order: 100
componentName: "azurecosmosdbsqlroleassignment"
---

# Azure Cosmos DB SQL Role Assignment

Creates a Cosmos DB SQL (NoSQL) API role assignment inside an AzureCosmosdbAccount -- the grant binding a Cosmos data-plane role to a Microsoft Entra principal at an account, database, or container scope. This is how workload identities get data access; with the account's key authentication disabled, these grants are the only way clients connect.

## What Gets Created

When you deploy an AzureCosmosdbSqlRoleAssignment resource, Planton provisions:

- **Cosmos DB SQL Role Assignment** -- an `azurerm_cosmosdb_sql_role_assignment` on the referenced account, binding the role definition to the principal at the declared scope

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbAccount** (referenced through `cosmosdbAccountId`); must be a GLOBAL_DOCUMENT_DB (SQL API) account
- **A principal** to grant to -- typically an AzureUserAssignedIdentity's `principal_id` output
- **A role** to bind -- a built-in's well-known ID, or an AzureCosmosdbSqlRoleDefinition's output

## Quick Start

Create a file `grant.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbSqlRoleAssignment
metadata:
  name: app-data-access
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbSqlRoleAssignment.app-data-access
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: my-app-cosmos
      fieldPath: status.outputs.cosmosdb_account_id
  # The built-in Data Contributor, by its well-known ID composed on
  # the account's ARM ID.
  roleDefinitionId:
    value: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.DocumentDB/databaseAccounts/my-cosmos/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002
  principalId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: my-app-identity
      fieldPath: status.outputs.principal_id
  scope:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: my-app-cosmos
      fieldPath: status.outputs.cosmosdb_account_id
```

Deploy:

```shell
planton apply -f grant.yaml
```

The principal must be the Entra OBJECT ID -- a client (application) ID is accepted by ARM but grants nothing, because no directory object carries it. Scope the grant as narrowly as the workload allows: `{account-id}/dbs/{db}` for one database, `.../colls/{container}` for one container.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `role_assignment_id` | The fully-scoped ARM id of the grant record |
| `role_assignment_guid` | The assignment's GUID resource name (pinned or generated) |
| `cosmosdb_account_name` | The account/grant pair, without a second reference |

## Related Resources

- [Azure Cosmos DB Account](/docs/catalog/azure/cosmos-db-account) -- the parent account
- [Azure Cosmos DB SQL Role Definition](/docs/catalog/azure/cosmos-db-sql-role-definition) -- custom roles this grant can bind
- [Azure User Assigned Identity](/docs/catalog/azure/user-assigned-identity) -- the workload identity most grants target
