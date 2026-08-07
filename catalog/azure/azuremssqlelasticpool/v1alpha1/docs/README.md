# AzureMssqlElasticPool -- Design Research

## The Resource

An Azure SQL elastic pool (`Microsoft.Sql/servers/elasticPools`) is a
shared-compute container: member databases draw from its eDTUs/vCores
instead of carrying their own SKU. The component maps onto
`azurerm_mssql_elasticpool` (azurerm v4.x,
`internal/services/mssql/mssql_elasticpool_resource.go`),
parity-verified against pulumi-azure v6 (`mssql.ElasticPool`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `pool_name` | Required, ForceNew |
| `server_name` + `resource_group_name` | derived from `server_id` | The spec carries ONE authoritative parent FK (AzureMssqlServer.server_id); both engines derive name + RG from the ARM id (the parent-derivation precedent) |
| `location` | `region` | Explicit and required -- ARM requires the pool's location and it MUST match the server's; a data-source lookup would break offline plan proofs, so the field is explicit with the must-match contract documented (ARM rejects a mismatch) |
| `sku.name` / `sku.tier` / `sku.family` / `sku.capacity` | `sku_name` + `capacity` | tier and family are PURE FUNCTIONS of the name (verified in azurerm's own validation helper), so both engines derive them -- a name/tier/family mismatch is unrepresentable; azurerm's CustomizeDiff that exists to catch mismatches has nothing left to catch |
| `per_database_settings` | `per_database_settings` | Required; min ≤ max CEL |
| `max_size_gb` XOR `max_size_bytes` | same | azurerm's ConflictsWith → XOR CEL |
| `zone_redundant` | same | |
| `enclave_type` | `enclave_type` enum | Pool-wide (every member database must share it) |
| `license_type` | `license_type` enum | vCore pools only (CEL) |
| `high_availability_replica_count` | same | Hyperscale pools only (CEL), 0-4 |
| `maintenance_configuration_name` | same | default SQL_Default; member databases inherit it |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Decomposition Decisions

- **The pool is a first-class kind**: independent billing container,
  many-per-server, FK-referenced by every pooled database
  (`AzureMssqlDatabase.elastic_pool_id`). It cannot fold into the server
  (many-per-parent) nor into the database (one pool serves many).

## Recorded Skips (with reasons)

- **Gen4 hardware families** (`GP_Gen4`, `BC_Gen4`) -- retired hardware
  Azure no longer provisions; the closed sku vocabulary carries only
  provisionable families. Widening the vocabulary is additive if a
  region ever resurrects them.

## Design Decisions

- **`sku_name` + `capacity` instead of azurerm's four-field sku block.**
  azurerm requires name+tier+family+capacity and then validates that
  tier and family match the name; both are derivable (its own
  `getTierFromName`/`getFamilyFromName` helpers prove it), so the spec
  takes the two real inputs and the modules derive the rest.
- **`region` is explicit rather than looked up from the server.** A
  data-source lookup keeps the spec minimal but breaks offline `tofu
  plan` proofs (the plan would need the server to exist) and adds a
  hidden runtime dependency; the explicit field carries the documented
  must-match-server contract that ARM enforces.

## Operational Behavior Worth Knowing

- **Pools create in ~2-4 minutes.**
- **A pool cannot be destroyed while databases are members** -- detach
  or destroy the databases first (the composed teardown order handles
  this when both are Planton-managed).
- **Per-database `min_capacity` × database count cannot exceed the
  pool's capacity** -- ARM rejects the overflow at membership time.

## Composition

- `server_id` → `AzureMssqlServer.status.outputs.server_id`
- `elastic_pool_id` output ← referenced by
  `AzureMssqlDatabase.elastic_pool_id` (with `sku_name: ElasticPool`)
