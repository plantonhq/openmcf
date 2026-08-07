# AzureCosmosdbSqlRoleDefinition -- Design Research

## The Resource

A Cosmos DB SQL (NoSQL) API role definition
(`Microsoft.DocumentDB/databaseAccounts/sqlRoleDefinitions`) is a named
bundle of DATA-PLANE permissions in Cosmos DB's own RBAC system --
separate from ARM RBAC, whose roles manage the account but grant no
access to the documents inside it. The component maps onto
`azurerm_cosmosdb_sql_role_definition` (azurerm v4.x,
`internal/services/cosmos/cosmosdb_sql_role_definition_resource.go`,
driven through the management-plane `2024-08-15/rbacs` API),
parity-verified against pulumi-azure v6 (`cosmosdb.SqlRoleDefinition`).

## Field Mapping (azurerm -> spec)

| azurerm | spec | Notes |
|---|---|---|
| `resource_group_name` + `account_name` | `cosmosdb_account_id` | azurerm addresses Cosmos RBAC resources by the (resource group, account, GUID) trio; the spec models a single parent ARM-id reference and both modules parse the trio from it identically -- no redundant, contradictable state. ForceNew |
| `name` | `role_name` | Required; the role's DISPLAY name (ARM property `roleName`), unique among the account's definitions; updatable in place. Named `role_name` because the GUID resource name is a separate coordinate (`role_definition_id`) |
| `type` | `type` | Optional closed enum (CUSTOM_ROLE / BUILT_IN_ROLE mapping to ARM's CustomRole / BuiltInRole); unspecified sends nothing so azurerm's CustomRole default applies identically on both engines. ForceNew |
| `assignable_scopes` | `assignable_scopes` | Required (min 1), repeated StringValueOrRef defaulting to the account's ARM ID -- the account-wide posture composes with zero literals; database/container paths are literals composed on the account ID (references cannot append suffixes). Scopes above the account are not enforceable (documented; Azure rejects them at apply). Updatable |
| `permissions[].data_actions` | `permissions[].data_actions` | Required (min 1 block, min 1 action each, non-empty strings) -- mirrors the provider's Required flags. Cosmos supports ALLOW rules only: no not_data_actions, no control-plane actions in this RBAC system. Updatable |
| `role_definition_id` | `role_definition_id` | Optional pinned UUID (spec CEL); unset lets the provider generate one at create time. ForceNew |

No tags field: ARM does not support tags on Cosmos child resources, so
the platform's identity tags live on the account.

## Decomposition Decisions

- **First-class kind, not a fold**: role definitions are many-per-account
  with independent lifecycles, and they are FK-referenced -- an
  AzureCosmosdbSqlRoleAssignment binds a definition by its resource ID.
  The same split verdict as the ARM RBAC pair.
- **The parent is a single ARM-id FK**: the account's
  `cosmosdb_account_id` output is the one authoritative reference; the
  resource-group and account names azurerm wants are derived, never
  asked for twice.
- **The `role_definition_id` OUTPUT carries the fully-scoped resource
  ID** (`{account-id}/sqlRoleDefinitions/{guid}`) -- exactly what an
  assignment's `role_definition_id` field consumes, so the composition
  seam needs zero translation. The bare GUID is exported separately as
  `role_definition_guid`.

## Apply-Time Contracts (provider source-diff)

- The provider locks on the account name around create/update/delete --
  Cosmos serializes control-plane writes per account. Module-internal
  behavior; nothing to model.
- When `role_definition_id` is unset the provider generates a UUID
  client-side before the ARM call. Both engines inherit this from their
  shared provider code.
- Assignable scopes above the account and non-existent role paths are
  rejected by ARM at apply -- cross-resource contracts that stay
  apply-time (a spec CEL cannot see the account ID inside a reference;
  the known protovalidate constraint).

## Recorded Skips (with reasons)

Nothing skipped: `azurerm_cosmosdb_sql_role_definition` exposes exactly
the parent addressing, type, name, assignable scopes, and permissions
surface, and the spec models all of it.

## Related Surface Deliberately Not Modeled Here

- **Mongo data-plane RBAC** (`azurerm_cosmosdb_mongo_role_definition` +
  `azurerm_cosmosdb_mongo_user_definition`) is a different API family
  on MONGO_DB accounts with its own inheritance and user model -- a
  candidate pair evaluated on its own merit when the Mongo RBAC story
  is taken up.

## Operational Behavior Worth Knowing

- **Include `readMetadata` in practically every role** -- Cosmos SDKs
  read database/container metadata and partition-key ranges before any
  data operation; a role without it fails clients in confusing ways.
- **Built-ins need no definition resource**: Data Reader
  (`00000000-0000-0000-0000-000000000001`) and Data Contributor
  (`...0002`) already exist in every account -- assign them by
  well-known ID. Author a custom definition only when neither fits.
- **This RBAC surface is SQL-API-only**: role definitions exist on
  GLOBAL_DOCUMENT_DB accounts; Mongo/Cassandra/Gremlin/Table accounts
  carry their own mechanisms.
- **Permissions are allow-only unions**: blocks are additive, and there
  is no carve-out mechanism -- express "everything except X" by listing
  everything-but-X.
- **Renaming is safe**: assignments track the role by GUID; the display
  name is an in-place update.

## Composition

- `cosmosdb_account_id` -> `AzureCosmosdbAccount.status.outputs.cosmosdb_account_id`
- `assignable_scopes[]` -> the same account output (account-wide), or
  literal `{account-id}/dbs/{db}[/colls/{container}]` paths
- `role_definition_id` output <- `AzureCosmosdbSqlRoleAssignment.role_definition_id`
  (the grant binding this role to a principal)
