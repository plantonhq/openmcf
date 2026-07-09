---
title: "Read-Only Role"
description: "This preset creates a custom read-only role covering the full read surface -- point reads, SQL queries, and the change feed -- plus the metadata access every Cosmos SDK needs before its first data..."
type: "preset"
rank: "01"
presetSlug: "01-read-only-role"
componentSlug: "cosmos-db-sql-role-definition"
componentTitle: "Cosmos DB SQL Role Definition"
provider: "azure"
icon: "package"
order: 1
---

# Read-Only Role

This preset creates a custom read-only role covering the full read
surface -- point reads, SQL queries, and the change feed -- plus the
metadata access every Cosmos SDK needs before its first data
operation. It mirrors the built-in Data Reader; author it as a custom
role when you want the definition under your own governance (your
naming, your assignable-scope narrowing) rather than Microsoft's.

## When to Use

- Read-side services, analytics jobs, and dashboards that must never
  mutate data
- A base to narrow further (drop `executeQuery` for point-read-only
  materializers, drop `readChangeFeed` for request-response services)
- Teams standardizing on custom roles for every grant so audits read
  uniformly

## Key Configuration Choices

- **`readMetadata` is included deliberately** -- SDKs list databases,
  containers, and partition-key ranges before any read; a role without
  it fails clients in confusing ways
- **Account-wide assignable scope** -- assignments of this role can
  target the account, any database, or any container; narrow the LIST
  to database paths to pre-authorize where it may ever be granted
- **Allow-only semantics** -- Cosmos data-plane RBAC has no carve-outs;
  this role is read-only because it lists only read actions

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name (a GLOBAL_DOCUMENT_DB account) | Your Cosmos composition |

## Downstream Wiring

Grants bind this role through its fully-scoped ID output:

```yaml
# On an AzureCosmosdbSqlRoleAssignment
roleDefinitionId:
  valueFrom:
    kind: AzureCosmosdbSqlRoleDefinition
    name: my-read-only-role
    fieldPath: status.outputs.role_definition_id
```
