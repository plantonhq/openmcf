# Azure Cosmos DB SQL Role Assignment

Grants a Cosmos DB data-plane role to a Microsoft Entra principal at a scope inside one SQL (NoSQL) API account. Cosmos DB carries its own RBAC system, separate from ARM RBAC — this assignment is how an Entra identity (a workload's managed identity, a CI principal, an operator) gets data access. With the account's local (key) authentication disabled, grants like this are the ONLY way clients connect: the fully keyless posture.

Every grant has three coordinates: WHAT (`roleDefinitionId` — a built-in role by well-known ID, or a custom AzureCosmosdbSqlRoleDefinition by reference), WHO (`principalId` — the Entra OBJECT ID), and WHERE (`scope` — the account, one database, or one container; permissions inherit downward).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cosmos DB SQL Role Assignment** -- a GUID-identified grant record binding the role to the principal at the scope

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **A Cosmos DB account** speaking the SQL (NoSQL) API.
- **The principal's OBJECT ID** — for a Planton-managed identity, reference the AzureUserAssignedIdentity's `principal_id` output; for users, groups, or external service principals, the object ID from Entra. The object ID is NOT the application (client) ID — an assignment authored with a client ID deploys successfully and grants nothing.
- **For custom roles**: an AzureCosmosdbSqlRoleDefinition whose assignable scopes cover the grant's scope.

## Deploy

### Console

Open the deployment store, find **Azure Cosmos DB SQL Role Assignment**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **workload-data-contributor** preset in the [Presets](#presets) tab for the common workload grant.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbSqlRoleAssignment
metadata:
  name: orders-api-data-contributor
  org: acme-corp
  env: prod
spec:
  cosmosdbAccountId:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: orders-db
      fieldPath: status.outputs.cosmosdb_account_id
  roleDefinitionId:
    value: /subscriptions/<sub>/resourceGroups/<rg>/providers/Microsoft.DocumentDB/databaseAccounts/<account>/sqlRoleDefinitions/00000000-0000-0000-0000-000000000002
  principalId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: orders-api-identity
      fieldPath: status.outputs.principal_id
  scope:
    valueFrom:
      kind: AzureCosmosdbAccount
      name: orders-db
      fieldPath: status.outputs.cosmosdb_account_id
```

```shell
planton apply -f cosmosdb-sql-role-assignment.yaml
```

This grants the built-in Data Contributor account-wide to the workload's managed identity. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the account, the identity, a custom role, and this grant compose in one InfraPipeline: the pipeline resolves `cosmosdb_account_id`, `principal_id`, and `role_definition_id` into the grant in dependency order.

## Key Configuration

**Role** -- Built-in roles (Data Reader `…0001`, Data Contributor `…0002`) exist in every account and are addressed by their well-known GUID composed on the account's ARM ID as a literal. Custom roles reference an AzureCosmosdbSqlRoleDefinition's output with zero translation. Rebinding the role is the grant's ONE in-place update.

**Scope** -- The account (the default reference), one database (`{account-id}/dbs/{db}`), or one container (`{account-id}/dbs/{db}/colls/{container}`). Sub-account scopes are literals composed on the account ID. Prefer the narrowest scope that satisfies the use case; the scope must sit at or below one of the role definition's assignable scopes.

**Principal** -- The Entra OBJECT ID. Referencing a managed identity's `principal_id` output makes the object-ID-vs-client-ID mistake unrepresentable.

**Immutability** -- Account, principal, scope, and the pinned GUID all replace the assignment; only the role binding updates in place — each grant record stays an auditable fact.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| AzureCosmosdbAccount | `cosmosdbAccountId` | `status.outputs.cosmosdb_account_id` |
| AzureCosmosdbSqlRoleDefinition | `roleDefinitionId` | `status.outputs.role_definition_id` |
| AzureUserAssignedIdentity | `principalId` | `status.outputs.principal_id` |
| AzureCosmosdbAccount | `scope` | `status.outputs.cosmosdb_account_id` |

### What This Component Produces

| Output | Description | Consumed By |
|--------|-------------|-------------|
| `role_assignment_id` | The fully-scoped ARM ID of the grant record | Audit tooling |
| `role_assignment_guid` | The assignment's GUID resource name | Automation that cites the grant |
| `cosmosdb_account_name` | The parent account's name | Automation that needs the account/grant pair |
