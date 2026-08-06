# AzureMysqlFlexibleServer -- Design Research

## The Resource

An Azure Database for MySQL Flexible Server
(`Microsoft.DBforMySQL/flexibleServers`) is Azure's managed MySQL:
per-server compute/storage sizing, zone-redundant high availability, a
Microsoft Entra administrator, customer-managed-key encryption, read
replicas, and point-in-time restore. The component maps onto
`azurerm_mysql_flexible_server` (azurerm v4.x,
`internal/services/mysql/mysql_flexible_server_resource.go`) and its child
resources, parity-verified against pulumi-azure v6 (`mysql.FlexibleServer*`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `server_name` | Required, ForceNew, globally unique (DNS name) |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `create_mode` | `create_mode` enum | DEFAULT / POINT_IN_TIME_RESTORE / REPLICA / GEO_RESTORE; unspecified not sent (same as azurerm's omitted default) |
| `source_server_id` | `source_server_id` | FK → this kind's own `server_id` output -- the replica/restore seam; mode-paired by CEL |
| `point_in_time_restore_time_in_utc` | same | RFC-3339; required for POINT_IN_TIME_RESTORE, forbidden elsewhere (CEL) -- GEO_RESTORE restores the latest geo-backup and takes no timestamp |
| `replication_role` | `replication_role` enum | Only legal value NONE (replica promotion); day-2 only -- Azure rejects it at create |
| `administrator_login` / `administrator_password` | same | Required for a fresh server (CEL); replicas/restores inherit; password auth can never be disabled on MySQL; login reserved-name rule mirrored |
| `version` | `version` | "5.7" / "8.0.21" / "8.4", default "8.0.21"; only sent for a fresh server (replicas/restores inherit); downgrade = ForceNew |
| `sku_name` | `sku_name` | Pattern-validated {TIER}\_Standard\_{SIZE}; required for DEFAULT (CEL), a replica left unset inherits |
| `storage` block | `storage` message | `size_gb` 20-16384 (shrink = ForceNew), `iops` 360-48000 XOR `io_scaling_enabled` (azurerm create-time check → CEL), `auto_grow_enabled` (Azure default TRUE -- opposite of PostgreSQL), `log_on_disk_enabled` |
| `zone` | `zone` | "1"/"2"/"3"; post-create changes only via planned failover |
| `high_availability` | `high_availability` | Mode enum (ZONE_REDUNDANT / SAME_ZONE) + standby zone |
| `maintenance_window` | `maintenance_window` | day/hour/minute; absence = system-managed window |
| `backup_retention_days` | same | 1-35 (PostgreSQL's floor is 7), default 7 |
| `geo_redundant_backup_enabled` | same | Plain bool, ForceNew |
| `public_network_access` | `public_network_access` enum | String enum in azurerm (not PostgreSQL's bool); unset lets Azure derive (Enabled publicly, Disabled when injected); ENABLED-on-injected rejected by CEL |
| `delegated_subnet_id` / `private_dns_zone_id` | same | FK → AzureSubnet / AzurePrivateDnsZone; pairing enforced by CEL; BOTH ForceNew on MySQL (the DNS zone is not ForceNew on PostgreSQL) |
| `identity` | `user_assigned_identity_ids` | MySQL supports UserAssigned ONLY (no system-assigned flavor), so the spec models the identity list directly instead of a type-plus-ids message |
| `customer_managed_key` | `customer_managed_key` | Key FK → AzureKeyVaultKey `versionless_id` (rotation propagates); MySQL adds a geo-backup key + identity pair (RequiredWith → CEL); primary identity required via CEL |
| `tags` | `tags` | User tags merged over Planton-derived tags |

Child resources folded into the spec (see Decomposition):

| azurerm | spec |
|---|---|
| `azurerm_mysql_flexible_database` | `databases[]` (name/charset/collation, all ForceNew) |
| `azurerm_mysql_flexible_server_firewall_rule` | `firewall_rules[]` (IPv4 range allowlist) |
| `azurerm_mysql_flexible_server_configuration` | `server_parameters` map (user overrides; delete = reset to default) |
| `azurerm_mysql_flexible_server_active_directory_administrator` | `aad_administrator` (singular -- MySQL supports exactly one; requires a user-assigned identity, unlike PostgreSQL's principal-only grant) |

## Decomposition Decisions

- **Databases, firewall rules, server parameters, and the Entra
  administrator FOLD into the server spec.** None has an independent
  lifecycle apart from its server, none is FK-referenced by any other
  kind, and all are configured as one administrative unit with the
  server. Each is still its own Azure sub-resource on both engines, so
  per-entry changes never touch the server.
- **A replica is another `AzureMysqlFlexibleServer`** (create_mode
  REPLICA + `source_server_id` referencing the primary's `server_id`
  output) -- a primary-plus-replicas topology composes in one manifest
  set with no dedicated replica kind. The `replica_capacity` output
  reports the source's remaining replica budget.

## Recorded Skips (with reasons)

- **`administrator_password_wo` / `administrator_password_wo_version`** --
  Terraform's write-only-attribute ergonomic for the same password value;
  duplicating the sensitive `administrator_password` field with no pulumi
  analog would add a second way to say one thing.
- **`customer_managed_key.managed_hsm_key_id`** -- deprecated in azurerm
  4.x (folds into `key_vault_key_id` in 5.0); the spec's single
  `key_vault_key_id` already accepts a Managed HSM key's data-plane URI.
- **Create-mode `Update`** -- azurerm's internal vehicle for
  login-change calls, not a state a user declares.

## Design Decisions

- **The storage profile is a message, not flat fields.** azurerm models
  MySQL storage as a block with two mutually exclusive IOPS shapes
  (provisioned `iops` vs `io_scaling_enabled`); the message keeps the
  XOR local and mirrors azurerm's create-time check as a message CEL.
- **`public_network_access` is a tri-state enum, not a bool.** MySQL's
  azurerm contract is a string with a derived default; modeling unset
  explicitly lets Azure derive DISABLED for injected servers instead of
  the module guessing.
- **The identity list is `user_assigned_identity_ids` directly.** MySQL
  supports no system-assigned identity, so a type-plus-ids identity
  message (PostgreSQL's shape) would carry a dead enum; the list of UAI
  FKs is the honest surface.
- **The Entra administrator is singular and identity-backed.** MySQL
  admits exactly one AAD administrator and requires a user-assigned
  identity to validate directory tokens; the CEL requires that identity
  to be attached to the server. `object_id` defaults to referencing a
  UAI's CLIENT id because MySQL validates managed-identity tokens
  against the application ID (documented on the field).
- **`version` is sent only for fresh servers on both engines.** Replicas
  and restores inherit the source's version; materializing the spec
  default ("8.0.21") onto a replica would fight the service.

## Operational Behavior Worth Knowing

- **Creates run ~10-15 minutes**; HA roughly doubles it. Destroys are
  fast (~2 min).
- **Static server parameters need a restart**: Azure applies the value
  but reports "pending restart". Neither engine restarts automatically --
  restart is a control-plane action outside declarative state.
- **`administrator_login` is immutable once set**; the password rotates
  in place.
- **Zone changes are failovers**: after creation, `zone` and
  `standby_availability_zone` can only change by swapping them together,
  which Azure executes as a planned failover.
- **Burstable SKUs** (`B_Standard_*`) support neither high availability
  nor read replicas (`replica_capacity` reports 0) -- the cheapest tier
  is a genuinely different capability envelope.
- **CMK identity ordering**: the user-assigned identity must hold
  wrap/unwrap on the key's vault BEFORE the server create starts, and
  the key's vault must have purge protection enabled.
- **5.7 → 8.0.21 upgrades in place** (irreversible); any other version
  transition replaces the server.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `delegated_subnet_id` → `AzureSubnet.status.outputs.subnet_id` (subnet
  delegated to `Microsoft.DBforMySQL/flexibleServers`)
- `private_dns_zone_id` → `AzurePrivateDnsZone.status.outputs.zone_id`
- `source_server_id` → another server's `status.outputs.server_id`
- `user_assigned_identity_ids[]` / `customer_managed_key.*_identity_id` /
  `aad_administrator.identity_id` →
  `AzureUserAssignedIdentity.status.outputs.identity_id`
- `customer_managed_key.key_vault_key_id` →
  `AzureKeyVaultKey.status.outputs.versionless_id`
- `aad_administrator.object_id` →
  `AzureUserAssignedIdentity.status.outputs.client_id` (MySQL validates
  managed-identity tokens against the CLIENT id, not the principal id)
- `server_id` output is consumed by
  `AzurePrivateEndpoint.private_connection_resource_id` and
  replica/restore servers' `source_server_id`
- `fqdn` + `administrator_login` outputs are what applications build
  connection strings from; `replica_capacity` sizes replica topologies
