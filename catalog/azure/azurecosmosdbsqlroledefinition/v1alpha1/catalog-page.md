# Azure Cosmos DB SQL Role Definition

Creates a Cosmos DB SQL (NoSQL) API role definition inside an AzureCosmosdbAccount -- a named bundle of data-plane permissions that Cosmos DB SQL role assignments bind to Microsoft Entra principals. This is Cosmos DB's own RBAC system, separate from ARM RBAC: ARM roles manage the account, these roles govern the documents inside it.

## What Gets Created

When you deploy an AzureCosmosdbSqlRoleDefinition resource, Planton provisions:

- **Cosmos DB SQL Role Definition** -- an `azurerm_cosmosdb_sql_role_definition` on the referenced account, carrying the role's data actions and assignable scopes

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureCosmosdbAccount** to create the definition in (referenced through `cosmosdbAccountId`); the account must be a GLOBAL_DOCUMENT_DB (SQL API) account

## Quick Start

Create a file `role.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureCosmosdbSqlRoleDefinition
metadata:
  name: app-reader
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureCosmosdbSqlRoleDefinition.app-reader
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: my-app-cosmos
      fieldPath: status.outputs.cosmosdb_account_id
  roleName: app-reader
  assignableScopes:
    - valueFrom:
        kind: AzureCosmosdbAccount
        name: my-app-cosmos
        fieldPath: status.outputs.cosmosdb_account_id
  permissions:
    - dataActions:
        - Microsoft.DocumentDB/databaseAccounts/readMetadata
        - Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/items/read
        - Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/executeQuery
```

Deploy:

```shell
planton apply -f role.yaml
```

Before creating a custom role, check whether a built-in fits: Data Reader (`...0001`) and Data Contributor (`...0002`) exist in every account and are assigned directly by ID from an AzureCosmosdbSqlRoleAssignment -- no definition resource needed. Author a custom definition for everything between those two.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `role_definition_id` | The fully-scoped ARM id -- exactly what an AzureCosmosdbSqlRoleAssignment's `roleDefinitionId` consumes |
| `role_definition_guid` | The definition's GUID resource name (pinned or generated) |
| `role_name` | The display name as deployed |
| `cosmosdb_account_name` | The account/definition pair, without a second reference |

## Related Resources

- [Azure Cosmos DB Account](/docs/catalog/azure/azurecosmosdbaccount) -- the parent account
- [Azure Cosmos DB SQL Role Assignment](/docs/catalog/azure/azurecosmosdbsqlroleassignment) -- the grant binding this role to a principal
