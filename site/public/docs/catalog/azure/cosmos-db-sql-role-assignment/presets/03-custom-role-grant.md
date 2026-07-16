---
title: "Custom Role Grant"
description: "This preset completes the custom-RBAC composition: an AzureCosmosdbSqlRoleDefinition's fully-scoped ID flows into the grant by reference, and the definition, the grant, the identity, and the account..."
type: "preset"
rank: "03"
presetSlug: "03-custom-role-grant"
componentSlug: "cosmos-db-sql-role-assignment"
componentTitle: "Cosmos DB SQL Role Assignment"
provider: "azure"
icon: "package"
order: 3
---

# Custom Role Grant

This preset completes the custom-RBAC composition: an
AzureCosmosdbSqlRoleDefinition's fully-scoped ID flows into the grant
by reference, and the definition, the grant, the identity, and the
account all deploy from one manifest set. Use it whenever the
built-ins don't express the permission set -- the definition carries
WHAT, this grant adds WHO and WHERE.

## When to Use

- Binding any custom role (writer-without-delete, read-only-no-query,
  metadata-only) to a workload identity
- Composing a team's whole data-access posture -- roles plus grants --
  as reviewable manifests
- Environments where every grant must trace to an org-authored,
  version-controlled role definition

## Key Configuration Choices

- **`roleDefinitionId` by reference** -- the definition's
  `role_definition_id` output IS the fully-scoped ID the grant needs;
  no ID assembly, no drift between the two resources
- **Database-scoped here** -- the grant's scope must sit at or below
  one of the definition's `assignableScopes`; Azure enforces the
  pairing at apply
- **Rebinding is the one in-place update** -- pointing this grant at a
  different definition updates it; changing principal or scope
  replaces the grant record

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name | Your Cosmos composition |
| `<role-definition-resource-name>` | The AzureCosmosdbSqlRoleDefinition to bind | Your RBAC composition |
| `<identity-resource-name>` | The AzureUserAssignedIdentity to grant | Your identity composition |
| `<subscription-id>` / `<resource-group-name>` / `<account-name>` / `<database-name>` | The coordinates of the scope path | The account's `cosmosdb_account_id` output |
