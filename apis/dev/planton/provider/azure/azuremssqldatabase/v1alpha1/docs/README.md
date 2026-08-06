# AzureMssqlDatabase -- Design Research

## The Resource

An Azure SQL Database (`Microsoft.Sql/servers/databases`) is the unit of
compute and billing in Azure SQL's logical-server model: SKU, storage
ceiling, availability, backups, and encryption are all per-database. The
component maps onto `azurerm_mssql_database` (azurerm v4.x,
`internal/services/mssql/mssql_database_resource.go`), parity-verified
against pulumi-azure v6 (`mssql.Database`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `database_name` | Required, ForceNew; azurerm's charset/ending rule mirrored as a CEL |
| `server_id` | `server_id` | FK → AzureMssqlServer.server_id -- the parent seam |
| `sku_name` | `sku_name` | Full vocabulary (DTU/vCore/serverless/HS/DW/DS/ElasticPool/Free) as a pattern CEL; unset lets Azure compute its serverless default; conditional ForceNew on HS↔non-HS transitions (documented) |
| `elastic_pool_id` | `elastic_pool_id` | FK → AzureMssqlElasticPool.elastic_pool_id; pairing contracts verified in azurerm source: pooled ⇔ sku "ElasticPool", maintenance window must be unset (both CELs) |
| `max_size_gb` | `max_size_gb` | double 0.1-4096 -- fractional sizes are legal ARM values (the prior int32 "for simplicity" shape dropped them) |
| `collation` | `collation` | ForceNew, default SQL_Latin1_General_CP1_CI_AS |
| `license_type` | `license_type` enum | BasePrice/LicenseIncluded |
| `auto_pause_delay_in_minutes` / `min_capacity` | same | Serverless-only (CEL-gated to GP_S_/HS_S_); -1 disables auto-pause |
| `read_replica_count` | same | Hyperscale-only (CEL-gated), 0-4 |
| `read_scale` | same | Premium/BC read-intent routing |
| `zone_redundant` / `ledger_enabled` | same | ledger ForceNew |
| `enclave_type` | `enclave_type` enum | VBS / DEFAULT_ENCLAVE (ARM's explicit "Default", distinct from unspecified so an update can actively clear); changing = ForceNew |
| `maintenance_configuration_name` | same | Conflicts with elastic_pool_id (CEL) |
| `create_mode` + sources | `create_mode` enum + 6 source fields | Full azurerm vocabulary; every mode↔source pairing is a CEL |
| `secondary_type` | `secondary_type` enum | Geo/Named, secondary modes only (CEL) |
| `storage_account_type` | `storage_account_type` enum | Wire Geo/GeoZone/Local/Zone; spec names *_REDUNDANT (enum-value scoping: "GEO" collides with secondary_type's) |
| `geo_backup_enabled` | same | optional bool default true, DW SKUs only |
| `sample_name` | `sample_name` | Single-value vocabulary (AdventureWorksLT) kept as azurerm's field |
| `identity` | `user_assigned_identity_ids` | UserAssigned-only on databases -- the list is the honest surface |
| `transparent_data_encryption_*` (3 fields) | same | Database-scoped CMK is a VERSIONED AzureKeyVaultKey.key_id FK; rotation-requires-key + key-requires-identity CELs |
| `import` | `import` message | bacpac: storage key + password sensitive; DEFAULT-mode-only CEL |
| `short_term_retention_policy` | same | 1-35 days + 12/24h differential cadence |
| `long_term_retention_policy` | same | ISO-8601 horizons + week_of_year; at-least-one CEL (this azurerm version carries no immutable_backups_enabled -- the clone wins over docs) |
| `threat_detection_policy` | same | Database-scoped Defender; email_account_admins is an Enabled/Disabled STRING on this resource (unlike the server's bool) -- the spec bool maps to the wire vocabulary in both modules |
| `tags` | `tags` | User tags merged over Planton-derived tags |

## Decomposition Decisions

- **The database is a first-class kind, not a server field.** It has a
  ~40-field surface of its own, an independent lifecycle, is
  many-per-server, and is FK-referenced (copy/secondary/restore sources
  reference `database_id`; pools are joined per database). The prior
  bundled shape could not express any of the lifecycle modes, serverless,
  Hyperscale, retention, or CMK.
- **A copy/secondary/restored database is another AzureMssqlDatabase**
  whose `create_mode` + source reference the original -- DR topologies
  compose in manifests with no dedicated kinds.

## Recorded Skips (with reasons)

- **`azurerm_mssql_database_extended_auditing_policy`** -- the
  server-scope auditing policy (modeled on AzureMssqlServer) already
  covers every database; a per-database override is a niche posture.
  Revisit on demand.
- **`azurerm_mssql_database_vulnerability_assessment_rule_baseline`** --
  tuning knobs for the classic (storage-backed) vulnerability
  assessment, which itself is superseded by express VA on the server.

## Design Decisions

- **`sku_name` stays one string** (azurerm's own shape here) with a
  pattern CEL over the full vocabulary -- unlike the pool, where ARM
  wants name+tier+family and the module derives the redundant two.
- **The mode↔source matrix is CELs, not module logic** -- six source
  fields each pair with exactly one create mode, so a wrong combination
  fails at validation, not after a deploy.
- **`enclave_type` models ARM's explicit "Default"** as DEFAULT_ENCLAVE,
  distinct from unspecified, so clearing a previously set enclave is
  expressible.

## Operational Behavior Worth Knowing

- **Basic/S0 databases create in ~1-2 minutes**; Hyperscale and large
  vCore SKUs take longer. Destroys are fast.
- **Hyperscale is a one-way door in place** -- leaving HS replaces the
  database; plan migrations with a copy.
- **A paused serverless database resumes on first connection** (seconds
  to a minute) -- the client sees a transient connect error if it does
  not retry.
- **LTR restores create NEW databases**
  (`RESTORE_LONG_TERM_RETENTION_BACKUP`), never restore in place.

## Composition

- `server_id` → `AzureMssqlServer.status.outputs.server_id`
- `elastic_pool_id` → `AzureMssqlElasticPool.status.outputs.elastic_pool_id`
- `creation_source_database_id` → another
  `AzureMssqlDatabase.status.outputs.database_id`
- `user_assigned_identity_ids[]` →
  `AzureUserAssignedIdentity.status.outputs.identity_id`
- `transparent_data_encryption_key_vault_key_id` →
  `AzureKeyVaultKey.status.outputs.key_id` (versioned)
- `database_id` output ← referenced by copy/secondary/restore databases
