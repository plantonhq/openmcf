# AzureCosmosdbSqlRoleAssignment -- Design Research

## The Resource

A Cosmos DB SQL (NoSQL) API role assignment
(`Microsoft.DocumentDB/databaseAccounts/sqlRoleAssignments`) is the
grant record in Cosmos DB's own RBAC system -- separate from ARM RBAC
-- binding a data-plane role to a Microsoft Entra principal at an
account, database, or container scope. The component maps onto
`azurerm_cosmosdb_sql_role_assignment` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_sql_role_assignment_resource.go`,
driven through the management-plane `2024-08-15/rbacs` API),
parity-verified against pulumi-azure v6 (`cosmosdb.SqlRoleAssignment`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` + `account_name` | `cosmosdb_account_id` | azurerm addresses Cosmos RBAC resources by the (resource group, account, GUID) trio; the spec models a single parent ARM-id reference and both modules parse the trio from it identically. ForceNew |
| `role_definition_id` | `role_definition_id` | Required StringValueOrRef defaulting to AzureCosmosdbSqlRoleDefinition's `role_definition_id` output (zero-translation custom-role composition); built-ins pass their well-known ID as a literal. The provider validates the ID's shape at plan time on Terraform; the Pulumi module enforces the same contract explicitly (see below). The one in-place update |
| `principal_id` | `principal_id` | Required StringValueOrRef defaulting to AzureUserAssignedIdentity's `principal_id` output -- the OBJECT-ID-not-client-ID trap documented on the field. ForceNew |
| `scope` | `scope` | Required StringValueOrRef defaulting to the account's ARM ID (account-wide grant); database/container paths are literals composed on the account ID (references cannot append suffixes). ForceNew |
| `name` | `name` | Optional pinned UUID (spec CEL); unset lets the provider generate one at create time. ForceNew |

No tags field: ARM does not support tags on Cosmos child resources, so
the platform's identity tags live on the account.

## Decomposition Decisions

- **First-class kind, not a fold**: grants are many-per-account (one
  per principal-role-scope triple), independently created and revoked,
  and a module must never own grants -- the same extraction verdict
  that shaped the ARM RBAC and Redis grant kinds.
- **The parent is a single ARM-id FK**: the account's
  `cosmosdb_account_id` output is the one authoritative reference; the
  resource-group and account names azurerm wants are derived, never
  asked for twice.
- **Built-in by literal, custom by reference**: the built-in roles
  exist in every account with well-known GUIDs, so their fully-scoped
  IDs pass as literal values; custom roles compose through the
  definition kind's `role_definition_id` output with zero translation.

## Apply-Time Contracts (provider source-diff)

- **Terraform validates `role_definition_id`'s shape at plan time**
  (`rbacs.ValidateSqlRoleDefinitionID`); the Pulumi SDK carries no
  equivalent, so the Pulumi module enforces the same
  `.../databaseAccounts/{account}/sqlRoleDefinitions/{guid}` contract
  before any ARM call -- both engines fail loudly and early on a
  malformed literal.
- The provider locks on the account name around create/update/delete --
  Cosmos serializes control-plane writes per account. Module-internal
  behavior; nothing to model.
- When `name` is unset the provider generates a UUID client-side before
  the ARM call. Both engines inherit this from their shared provider
  code.
- **Cross-resource contracts that stay apply-time** (a spec CEL cannot
  see inside references -- the known protovalidate constraint): the
  scope and the role definition must live inside THIS account; the
  scope must sit at or below one of the definition's assignable scopes;
  the principal must exist in the tenant's directory. Azure rejects
  violations with clear diagnostics.

## Recorded Skips (with reasons)

Nothing skipped: `azurerm_cosmosdb_sql_role_assignment` exposes exactly
the parent addressing, role binding, principal, scope, and GUID
surface, and the spec models all of it.

## Operational Behavior Worth Knowing

- **OBJECT ID, not client ID**: the most common grant mistake -- ARM
  accepts a client (application) ID and the assignment deploys, but no
  directory object carries it, so nothing is granted.
- **Rebinding is the one update**: changing `role_definition_id`
  updates the grant in place; changing the principal, scope, account,
  or pinned GUID replaces the record -- ARM's grant-record model.
- **Permissions inherit downward**: an account-scoped grant covers
  every database and container; prefer the narrowest scope the
  workload allows.
- **Grants take effect within seconds** -- Cosmos data-plane RBAC does
  not carry ARM RBAC's multi-minute propagation delays; a just-granted
  identity can usually connect immediately.
- **Keyless closure**: pair grants with the account's key-auth switch
  -- keys off plus explicit grants is the audit-friendly posture, and
  the account's `local_authentication` setting is where keys turn off.

## Composition

- `cosmosdb_account_id` -> `AzureCosmosdbAccount.status.outputs.cosmosdb_account_id`
- `role_definition_id` -> `AzureCosmosdbSqlRoleDefinition.status.outputs.role_definition_id`
  (custom roles), or the built-ins' well-known IDs as literals
- `principal_id` -> `AzureUserAssignedIdentity.status.outputs.principal_id`
- `scope` -> the same account output (account-wide), or literal
  `{account-id}/dbs/{db}[/colls/{container}]` paths
