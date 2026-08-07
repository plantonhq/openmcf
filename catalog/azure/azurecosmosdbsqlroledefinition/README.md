# AzureCosmosdbSqlRoleDefinition

A Cosmos DB SQL (NoSQL) API role definition inside an
AzureCosmosdbAccount: a named, reusable bundle of DATA-PLANE
permissions that AzureCosmosdbSqlRoleAssignment resources bind to
Microsoft Entra principals. Cosmos DB carries its own RBAC system,
separate from ARM RBAC -- an ARM role (even Owner) manages the account
but grants no access to the documents inside it. Together with the
account's key-authentication switch, this system is Cosmos DB's
keyless story.

## When to Use

Use AzureCosmosdbSqlRoleDefinition when neither built-in data role
fits:

- **Read-only with no query** -- point reads and change feed for a
  materializer, without `executeQuery`
- **Write-but-never-delete** -- create/replace/upsert items for an
  ingest workload, deletion reserved for operators
- **Metadata-only** -- a monitoring probe that lists databases and
  containers but touches no documents

Assigning a BUILT-IN role (Data Reader
`00000000-0000-0000-0000-000000000001`, Data Contributor
`00000000-0000-0000-0000-000000000002`) needs NO definition resource --
built-ins already exist in every account; reference their well-known
IDs directly from an AzureCosmosdbSqlRoleAssignment.

## Key Configuration

- `cosmosdb_account_id` -- the parent account, referenced from an
  AzureCosmosdbAccount's output (fixed at creation); must be a
  GLOBAL_DOCUMENT_DB (SQL API) account
- `role_name` -- the display name, unique among the account's
  definitions; renaming is an in-place update
- `assignable_scopes` -- WHERE assignments may be created: the account
  (the default reference) or database/container paths under it; scopes
  above the account are not enforceable in Cosmos data-plane RBAC
- `permissions[].data_actions` -- WHAT the role allows, as
  `Microsoft.DocumentDB/databaseAccounts/...` operation patterns;
  include `readMetadata` in practically every role (SDKs read metadata
  before any data operation). Cosmos supports ALLOW rules only -- no
  carve-outs
- `role_definition_id` -- optional pinned GUID; omit to generate one

## Composition

```yaml
cosmosdbAccountId:
  valueFrom:
    kind: AzureCosmosdbAccount
    name: app-cosmos
    fieldPath: status.outputs.cosmosdb_account_id
```

Grants consume the definition's `role_definition_id` output with zero
translation:

```yaml
# On an AzureCosmosdbSqlRoleAssignment
roleDefinitionId:
  valueFrom:
    kind: AzureCosmosdbSqlRoleDefinition
    name: app-writer-role
    fieldPath: status.outputs.role_definition_id
```

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
