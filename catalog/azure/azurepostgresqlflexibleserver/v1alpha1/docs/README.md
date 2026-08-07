# AzurePostgresqlFlexibleServer -- Design Research

## The Resource

An Azure Database for PostgreSQL Flexible Server
(`Microsoft.DBforPostgreSQL/flexibleServers`) is Azure's managed PostgreSQL:
per-server compute/storage sizing, zone-redundant high availability,
Microsoft Entra authentication, customer-managed-key encryption, and
point-in-time restore. The component maps onto
`azurerm_postgresql_flexible_server` (azurerm v4.x,
`internal/services/postgres/postgresql_flexible_server_resource.go`) and its
child resources, parity-verified against pulumi-azure v6
(`postgresql.FlexibleServer*`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `server_name` | Required, ForceNew, globally unique (DNS name); azurerm's real charset allows a leading digit |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `create_mode` | `create_mode` enum | DEFAULT / POINT_IN_TIME_RESTORE / REPLICA / GEO_RESTORE / REVIVE_DROPPED; unspecified not sent (same as azurerm's omitted default) |
| `source_server_id` | `source_server_id` | FK → this kind's own `server_id` output -- the replica/restore seam; mode-paired by CEL |
| `point_in_time_restore_time_in_utc` | same | RFC-3339; required for the two restore modes, forbidden elsewhere (CEL) |
| `replication_role` | `replication_role` enum | Only legal value NONE (replica promotion); day-2 only -- Azure rejects it at create |
| `administrator_login` / `administrator_password` | same | Conditional per CEL: required for a fresh password-auth server, forbidden when password auth is off; login reserved-name rule mirrored |
| `version` | `version` | "11"-"18", default "16"; only sent for a fresh server (replicas/restores inherit); downgrade = ForceNew |
| `sku_name` | `sku_name` | Pattern-validated {TIER}\_Standard\_{SIZE}; required for DEFAULT (CEL), a replica left unset inherits |
| `storage_mb` | `storage_mb` | Closed 12-value ladder; shrink = ForceNew |
| `storage_tier` | `storage_tier` enum | P4-P80; the size→valid-tier matrix mirrored as one message CEL |
| `auto_grow_enabled` | same | Plain bool, Azure default false |
| `zone` | `zone` | "1"/"2"/"3"; post-create changes only via planned failover |
| `high_availability` | `high_availability` | Mode enum (ZONE_REDUNDANT / SAME_ZONE) + standby zone (create-only) |
| `maintenance_window` | `maintenance_window` | day/hour/minute; absence = system-managed window |
| `backup_retention_days` | same | 7-35, default 7 |
| `geo_redundant_backup_enabled` | same | Plain bool, ForceNew |
| `public_network_access_enabled` | same | optional bool default true -- modeled explicitly (the prior spec derived it from subnet presence; azurerm's real contract has both dials) |
| `delegated_subnet_id` / `private_dns_zone_id` | same | FK → AzureSubnet / AzurePrivateDnsZone; pairing + public-access-off enforced by CEL |
| `authentication` | `authentication` | password (default true) + Entra + tenant (falls back to the deploying credential's tenant on both engines) |
| `identity` | `identity` | SYSTEM_ASSIGNED / USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED; ids FK → AzureUserAssignedIdentity |
| `customer_managed_key` | `customer_managed_key` | Key FK → AzureKeyVaultKey `versionless_id` (rotation propagates); geo-backup pair CEL-coupled; ForceNew |
| `cluster` | `cluster` | Elastic cluster (PG 17+, DEFAULT mode only -- both CEL-enforced); size 1-20 grows only |
| `tags` | `tags` | User tags merged over Planton-derived tags |

Child resources folded into the spec (see Decomposition):

| azurerm | spec |
|---|---|
| `azurerm_postgresql_flexible_server_database` | `databases[]` (name/charset/collation, all ForceNew) |
| `azurerm_postgresql_flexible_server_firewall_rule` | `firewall_rules[]` (IPv4 range allowlist) |
| `azurerm_postgresql_flexible_server_configuration` | `server_parameters` map (user overrides; delete = reset to default) |
| `azurerm_postgresql_flexible_server_active_directory_administrator` | `aad_administrators[]` (object id FK → UAI principal id) |

## Decomposition Decisions

- **Databases, firewall rules, server parameters, and Entra administrators
  FOLD into the server spec.** None has an independent lifecycle apart from
  its server, none is FK-referenced by any other kind, and all are
  configured as one administrative unit with the server. Each is still its
  own Azure sub-resource on both engines, so per-entry changes never touch
  the server.
- **A replica is another `AzurePostgresqlFlexibleServer`** (create_mode
  REPLICA + `source_server_id` referencing the primary's `server_id`
  output) -- a primary-plus-replicas topology composes in one manifest set
  with no dedicated replica kind.

## Recorded Skips (with reasons)

- **`administrator_password_wo` / `administrator_password_wo_version`** --
  Terraform's write-only-attribute ergonomic for the same password value;
  duplicating the sensitive `administrator_password` field with no pulumi
  analog would add a second way to say one thing.
- **`azurerm_postgresql_flexible_server_virtual_endpoint`** -- a stable
  connection name that spans a source↔replica PAIR and follows failover;
  it is a two-server topology object, not a property of one server.
  Adoption backlog: revisit as a standalone kind when replica topologies
  become a chart scenario.
- **`azurerm_postgresql_flexible_server_backup`** -- triggers a one-time
  on-demand backup; an imperative action, not declarative state.
- **Create-mode `Update`** -- azurerm's internal vehicle for
  login-change calls, not a state a user declares.

## Design Decisions

- **`public_network_access_enabled` modeled explicitly.** The prior spec
  derived it from subnet presence "to eliminate contradiction risk" --
  but azurerm's real contract has two independent dials, and the invented
  coupling made a no-public-access, no-VNet server (private endpoints
  only) inexpressible. CEL enforces the one true constraint (VNet
  injection requires public access off) and leaves the rest of the matrix
  open.
- **Entra auth + AAD administrators modeled in full** (the prior spec
  hard-coded password-only auth). An Entra-only server omits credentials
  entirely -- CEL mirrors azurerm's create-time contract in both
  directions.
- **The size→tier matrix is spec validation, not module logic.** azurerm
  validates storage tiers in CustomizeDiff; the equivalent lives in one
  message CEL so both engines reject an invalid tier before any cloud
  call.
- **`version` is sent only for fresh servers on both engines.** Replicas
  and restores inherit the source's version; materializing the spec
  default ("16") onto a replica would fight the service.
- **The admin-login reserved-name rule uses a field CEL** mirroring
  azurerm's `AdminUsernames` validator (reserved names + the `pg_`
  prefix), so the error surfaces at validation time, not after a
  15-minute create.

## Operational Behavior Worth Knowing

- **Creates run ~8-15 minutes**; HA roughly doubles it. Destroys are
  fast (~2 min).
- **Static server parameters need a restart**: Azure applies the value
  but reports "pending restart". Neither engine restarts automatically --
  restart is a control-plane action outside declarative state.
- **`administrator_login` is immutable once set** (azurerm forces
  replacement); the password rotates in place.
- **Zone changes are failovers**: after creation, `zone` and
  `standby_availability_zone` can only change by swapping them together,
  which Azure executes as a planned failover.
- **Burstable SKUs** (`B_Standard_*`) support neither high availability
  nor read replicas -- the cheapest tier is a genuinely different
  capability envelope.
- **CMK identity ordering**: the user-assigned identity must hold
  wrap/unwrap on the key's vault BEFORE the server create starts, and the
  key's vault must have purge protection enabled.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `delegated_subnet_id` → `AzureSubnet.status.outputs.subnet_id` (subnet
  delegated to `Microsoft.DBforPostgreSQL/flexibleServers`)
- `private_dns_zone_id` → `AzurePrivateDnsZone.status.outputs.zone_id`
- `source_server_id` → another server's `status.outputs.server_id`
- `identity.identity_ids[]` / `customer_managed_key.*_identity_id` →
  `AzureUserAssignedIdentity.status.outputs.identity_id`
- `customer_managed_key.key_vault_key_id` →
  `AzureKeyVaultKey.status.outputs.versionless_id`
- `aad_administrators[].object_id` →
  `AzureUserAssignedIdentity.status.outputs.principal_id`
- `server_id` output is consumed by
  `AzurePrivateEndpoint.private_connection_resource_id` and replica/restore
  servers' `source_server_id`
- `fqdn` + `administrator_login` outputs are what applications build
  connection strings from; `identity_principal_id` is the
  `AzureRoleAssignment` seam for the system-assigned identity
