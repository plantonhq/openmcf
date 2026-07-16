---
title: "Database-Scoped Role"
description: "This preset narrows WHERE a role may ever be granted, not just what it allows: its single assignable scope is one database's path, so an assignment of this role at the account level -- or in any..."
type: "preset"
rank: "03"
presetSlug: "03-database-scoped-role"
componentSlug: "cosmos-db-sql-role-definition"
componentTitle: "Cosmos DB SQL Role Definition"
provider: "azure"
icon: "package"
order: 3
---

# Database-Scoped Role

This preset narrows WHERE a role may ever be granted, not just what it
allows: its single assignable scope is one database's path, so an
assignment of this role at the account level -- or in any other
database -- is rejected by Azure at apply. Use it to give a team full
data access to THEIR database in a shared account, with the boundary
enforced by the definition itself rather than by review discipline on
every grant.

## When to Use

- Shared Cosmos accounts where each team owns one database and grants
  must not be able to escape it
- Pre-authorizing the blast radius of a powerful role (this one carries
  the full item surface) at definition time
- Compliance postures where "who could ever be granted what, where"
  must be answerable from the definitions alone

## Key Configuration Choices

- **The assignable scope is a literal database path** -- composed on
  the account's ARM ID (`{account-id}/dbs/{database}`); references
  cannot append path suffixes, so sub-account scopes are always
  literals
- **`items/*` is safe HERE because the scope is narrow** -- the
  wildcard grants every item operation including delete, contained to
  the one database this role can be assigned in
- **Assignments still pick their own scope** -- at or below the
  assignable scope: the whole database, or a single container in it

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name | Your Cosmos composition |
| `<subscription-id>` / `<resource-group-name>` / `<account-name>` | The coordinates of the account's ARM ID | The account's `cosmosdb_account_id` output |
| `<database-name>` | The database this role is confined to | Your AzureCosmosdbSqlDatabase |

## Downstream Wiring

```yaml
# On an AzureCosmosdbSqlRoleAssignment -- the grant's scope must sit
# at or below the definition's assignable scope
scope:
  value: /subscriptions/<subscription-id>/resourceGroups/<resource-group-name>/providers/Microsoft.DocumentDB/databaseAccounts/<account-name>/dbs/<database-name>
roleDefinitionId:
  valueFrom:
    kind: AzureCosmosdbSqlRoleDefinition
    name: my-database-scoped-role
    fieldPath: status.outputs.role_definition_id
```
