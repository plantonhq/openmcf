# AzureMssqlFailoverGroup -- Design Research

## The Resource

An Azure SQL Failover Group (`Microsoft.Sql/servers/failoverGroups`)
replicates databases from a primary logical server to partner servers in
other regions and provides failover-following listener endpoints. The
component maps onto `azurerm_mssql_failover_group` (azurerm v4.x,
`internal/services/mssql/mssql_failover_group_resource.go`), parity-verified
against pulumi-azure v6 (`mssql.FailoverGroup`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `server_id` | `server_id` | FK → AzureMssqlServer, ForceNew |
| `partner_server` (block, MinItems 1) | `partner_servers` | Each `{ server_id }`; location/role are provider-computed |
| `databases` (set) | `database_ids` | FK → AzureMssqlDatabase; must live on the primary |
| `read_write_endpoint_failover_policy` (block, MaxItems 1) | `read_write_endpoint_failover_policy` | mode + grace_minutes |
| `readonly_endpoint_failover_policy_enabled` | `readonly_endpoint_failover_policy_enabled` | Optional+Computed; the provider sends Disabled when unset |
| `tags` | `tags` | User tags merged over Planton-derived tags |
| `id` (computed) | `failover_group_id` output | |

## Key Design Decisions

- **The grace-minutes pairing is a CEL, mirroring the provider's
  CustomizeDiff.** azurerm enforces at apply that `grace_minutes` ≥ 60 for
  `Automatic` and is unset for `Manual`. This is expressible as a
  message-level CEL on the policy sub-message (it dereferences only scalar
  fields, not a StringValueOrRef), so it is front-loaded to validation.
- **partner_servers is a message list, not bare refs.** azurerm's
  `partner_server` block also carries provider-computed `location` and
  `role`; modeling it as a message leaves room for those to surface as
  outputs later without a breaking change, and reads as the topology it is.
- **Listener endpoints are composed, not read back.** azurerm exports only
  the group `id`; the read-write and read-only listener FQDNs are
  deterministic from the group name (`{name}.database.windows.net` /
  `{name}.secondary.database.windows.net`), so both modules compose them as
  outputs -- the value applications actually need.
- **Databases optional but meaningful.** The set may be empty (an empty
  group whose databases are added later), matching azurerm; the spec allows
  it while the docs note a DR group with no databases protects nothing.

## Composition Seams

The group consumes `AzureMssqlServer` (primary via `server_id`, partners via
`partner_servers[].server_id`) and `AzureMssqlDatabase` (`database_ids`). Its
`read_write_listener_endpoint` output is the connection target downstream
applications and connection strings reference.

## Live E2E

Live dual-engine E2E uses the shared MSSQL server fixture (westus3) as the
primary, a scenario-local partner server in eastus, and a scenario-local S0
database on the primary. Logical servers provision in minutes and are free;
the database is the only billed fixture. The group is verified via the ARM
API and torn down cleanly.
