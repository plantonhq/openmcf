# AzureCosmosdbSqlRoleAssignment

A Cosmos DB SQL (NoSQL) API role assignment inside an
AzureCosmosdbAccount: the grant binding a Cosmos data-plane role
(built-in or AzureCosmosdbSqlRoleDefinition) to a Microsoft Entra
principal at an account, database, or container scope. Cosmos DB
carries its own RBAC system, separate from ARM RBAC -- an ARM role
(even Owner) manages the account but grants no access to the documents
inside it. With the account's key authentication disabled, these
grants are the ONLY way clients connect: the fully keyless posture.

## When to Use

Use AzureCosmosdbSqlRoleAssignment for every Entra identity that
touches Cosmos data:

- **A workload's managed identity** -- the dominant composition: the
  app authenticates keylessly and this grant scopes exactly what it
  may do
- **A CI or migration principal** -- time-bounded, auditable data
  access without sharing account keys
- **An operator or group** -- container-scoped read access for
  debugging, without key custody

## Key Configuration

- `cosmosdb_account_id` -- the account the grant lives in (fixed at
  creation); must be a GLOBAL_DOCUMENT_DB (SQL API) account
- `role_definition_id` -- WHAT is permitted: a custom
  AzureCosmosdbSqlRoleDefinition by reference (the default), or a
  built-in by its well-known ID composed on the account -- Data Reader
  `{account-id}/sqlRoleDefinitions/00000000-0000-0000-0000-000000000001`,
  Data Contributor `...0002`. Rebinding is the one in-place update
- `principal_id` -- WHO: the Entra OBJECT ID (not the client ID);
  defaults to referencing an AzureUserAssignedIdentity's principal_id
- `scope` -- WHERE: the account (the default reference), or a literal
  `{account-id}/dbs/{db}[/colls/{container}]` path; permissions
  inherit downward, so prefer the narrowest scope that works
- `name` -- optional pinned GUID; omit to generate one

## Composition

```yaml
# Grant a custom role to a workload identity across one container
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: app-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
roleDefinitionId:
  valueFrom:
    kind: AzureCosmosdbSqlRoleDefinition
    name: app-writer-role
    fieldPath: status.outputs.role_definition_id
principalId:
  valueFrom:
    kind: AzureUserAssignedIdentity
    name: app-identity
    fieldPath: status.outputs.principal_id
scope:
  value: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.DocumentDB/databaseAccounts/{account}/dbs/app-data/colls/orders
```

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
