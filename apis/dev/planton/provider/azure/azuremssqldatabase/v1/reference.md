# AzureMssqlDatabase

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureMssqlDatabaseSpec** defines the configuration for creating an
Azure SQL Database on an AzureMssqlServer logical server.

In Azure SQL's logical-server model the DATABASE is the unit of compute
and billing: each database carries its own SKU (DTU or vCore),
storage ceiling, availability posture, backup policy, and encryption --
or joins an AzureMssqlElasticPool to share the pool's compute instead.

**SKU model** (`sku_name`):
- DTU tiers (bundled compute/IO): "Basic", "S0"-"S12" (Standard),
  "P1"-"P15" (Premium)
- vCore tiers: "{TIER}_{FAMILY}_{VCORES}" -- "GP_Gen5_2" (General
  Purpose), "BC_Gen5_4" (Business Critical), "HS_Gen5_2" (Hyperscale)
- Serverless (auto-pause, per-second billing): "GP_S_Gen5_1",
  "HS_S_Gen5_2"
- Data Warehouse: "DW100c"-"DW30000c"; Stretch: "DS100"-"DS2000"
- "ElasticPool": the database's compute comes from the pool referenced
  by elastic_pool_id
- "Free": one free database per subscription

**Lifecycle modes** (`create_mode`): a fresh database (DEFAULT), a copy,
a readable geo/named secondary, a point-in-time restore, recovery from a
geo-replicated backup, restore of a dropped database, or restore from a
long-term-retention backup. Each mode consumes its matching source
field; all are fixed at creation.

**Conditional replacement** (ARM's contract, enforced at update):
changing `sku_name` between Hyperscale and non-Hyperscale families
replaces the database, as does changing `enclave_type` -- plan
migrations accordingly.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlDatabase
metadata:
  name: test-mssql-db
spec:
  serverId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Sql/servers/test-mssql-server
  databaseName: orders
  # A serverless sku exercises the auto-pause/min-capacity dials.
  skuName: GP_S_Gen5_2
  autoPauseDelayInMinutes: 60
  minCapacity: 0.5
  maxSizeGb: 32
  collation: SQL_Latin1_General_CP1_CI_AS
  # Exercises the license, redundancy, and enclave enum mappings.
  licenseType: LICENSE_INCLUDED
  storageAccountType: ZONE_REDUNDANT
  enclaveType: VBS
  zoneRedundant: false
  ledgerEnabled: true
  maintenanceConfigurationName: SQL_Default
  # Exercises the database-scoped CMK with rotation.
  userAssignedIdentityIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/db-uai
  transparentDataEncryptionKeyVaultKeyId:
    value: https://test-vault.vault.azure.net/keys/db-tde/0123456789abcdef0123456789abcdef
  transparentDataEncryptionKeyAutomaticRotationEnabled: true
  # Exercises both retention sub-APIs.
  shortTermRetentionPolicy:
    retentionDays: 14
    backupIntervalInHours: 24
  longTermRetentionPolicy:
    weeklyRetention: P5W
    monthlyRetention: P12M
    yearlyRetention: P5Y
    weekOfYear: 26
  # Exercises the database-scoped Defender policy with the bool -> wire
  # string mapping for email_account_admins.
  threatDetectionPolicy:
    state: ENABLED
    disabledAlerts:
      - Sql_Injection
    emailAccountAdmins: true
    emailAddresses:
      - secops@contoso.com
    retentionDays: 30
  tags:
    team: data
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serverId` | `string \| valueFrom` | yes |  | AzureMssqlServer (`status.outputs.server_id`) |
| `spec.databaseName` | `string` | yes |  |  |
| `spec.skuName` | `string` |  |  |  |
| `spec.elasticPoolId` | `string \| valueFrom` |  |  | AzureMssqlElasticPool (`status.outputs.elastic_pool_id`) |
| `spec.maxSizeGb` | `double` |  |  |  |
| `spec.collation` | `string` | yes | `SQL_Latin1_General_CP1_CI_AS` |  |
| `spec.licenseType` | `enum` |  |  |  |
| `spec.autoPauseDelayInMinutes` | `int32` |  |  |  |
| `spec.minCapacity` | `double` |  |  |  |
| `spec.readReplicaCount` | `int32` |  |  |  |
| `spec.readScale` | `bool` |  |  |  |
| `spec.zoneRedundant` | `bool` |  |  |  |
| `spec.ledgerEnabled` | `bool` |  |  |  |
| `spec.enclaveType` | `enum` |  |  |  |
| `spec.maintenanceConfigurationName` | `string` |  |  |  |
| `spec.createMode` | `enum` |  |  |  |
| `spec.creationSourceDatabaseId` | `string \| valueFrom` |  |  | AzureMssqlDatabase (`status.outputs.database_id`) |
| `spec.secondaryType` | `enum` |  |  |  |
| `spec.restorePointInTime` | `string` |  |  |  |
| `spec.recoverDatabaseId` | `string` |  |  |  |
| `spec.recoveryPointId` | `string` |  |  |  |
| `spec.restoreDroppedDatabaseId` | `string` |  |  |  |
| `spec.restoreLongTermRetentionBackupId` | `string` |  |  |  |
| `spec.storageAccountType` | `enum` |  |  |  |
| `spec.geoBackupEnabled` | `bool` |  | `true` |  |
| `spec.sampleName` | `string` |  |  |  |
| `spec.userAssignedIdentityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.transparentDataEncryptionEnabled` | `bool` |  | `true` |  |
| `spec.transparentDataEncryptionKeyVaultKeyId` | `string \| valueFrom` |  |  | AzureKeyVaultKey (`status.outputs.key_id`) |
| `spec.transparentDataEncryptionKeyAutomaticRotationEnabled` | `bool` |  |  |  |
| `spec.import` | `AzureMssqlDatabaseImport` |  |  |  |
| `spec.import.storageUri` | `string` | yes |  |  |
| `spec.import.storageKey` | `string \| valueFrom` (sensitive) | yes |  |  |
| `spec.import.storageKeyType` | `enum` | yes |  |  |
| `spec.import.administratorLogin` | `string` | yes |  |  |
| `spec.import.administratorLoginPassword` | `string \| valueFrom` (sensitive) | yes |  |  |
| `spec.import.authenticationType` | `enum` | yes |  |  |
| `spec.import.storageAccountId` | `string` |  |  |  |
| `spec.shortTermRetentionPolicy` | `AzureMssqlDatabaseShortTermRetentionPolicy` |  |  |  |
| `spec.shortTermRetentionPolicy.retentionDays` | `int32` | yes |  |  |
| `spec.shortTermRetentionPolicy.backupIntervalInHours` | `int32` |  |  |  |
| `spec.longTermRetentionPolicy` | `AzureMssqlDatabaseLongTermRetentionPolicy` |  |  |  |
| `spec.longTermRetentionPolicy.weeklyRetention` | `string` |  |  |  |
| `spec.longTermRetentionPolicy.monthlyRetention` | `string` |  |  |  |
| `spec.longTermRetentionPolicy.yearlyRetention` | `string` |  |  |  |
| `spec.longTermRetentionPolicy.weekOfYear` | `int32` |  |  |  |
| `spec.threatDetectionPolicy` | `AzureMssqlDatabaseThreatDetectionPolicy` |  |  |  |
| `spec.threatDetectionPolicy.state` | `enum` | yes |  |  |
| `spec.threatDetectionPolicy.disabledAlerts` | `[]string` |  |  |  |
| `spec.threatDetectionPolicy.emailAccountAdmins` | `bool` |  |  |  |
| `spec.threatDetectionPolicy.emailAddresses` | `[]string` |  |  |  |
| `spec.threatDetectionPolicy.retentionDays` | `int32` |  |  |  |
| `spec.threatDetectionPolicy.storageEndpoint` | `string` |  |  |  |
| `spec.threatDetectionPolicy.storageAccountAccessKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.serverId

`string | valueFrom` · required

The logical server the database is created on, by ARM ID.
References an AzureMssqlServer's server_id output; the server's
resource group and location are derived from it. Fixed at creation.

- references: AzureMssqlServer (`status.outputs.server_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMssqlServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.databaseName

`string` · required

The database name, unique within the server: at most 128 characters,
no '<>*%&:\/?' or control characters, not ending with '.' or ' '.
Changing the name replaces the database.

- rule: database_name cannot contain <>*%&:\/? or end with '.' or ' '
- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.skuName

`string`

The compute SKU -- the database's most consequential dial (capacity
and cost). See the message comment for the vocabulary. Set
"ElasticPool" when the database joins a pool via elastic_pool_id.
Unset lets Azure apply its default (GP_S_Gen5_2 serverless).
Downgrading OUT of Hyperscale replaces the database.

- rule: sku_name must be a DTU tier (Basic, S0-S12, P1-P15), a vCore tier (GP_/BC_/HS_ with optional _S_ serverless, e.g. GP_Gen5_2, GP_S_Gen5_1), DW/DS, ElasticPool, or Free

### spec.elasticPoolId

`string | valueFrom`

The elastic pool the database joins, by ARM ID. References an
AzureMssqlElasticPool's elastic_pool_id output. A pooled database
must set sku_name to "ElasticPool" and must NOT set
maintenance_configuration_name (it inherits the pool's).

- references: AzureMssqlElasticPool (`status.outputs.elastic_pool_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMssqlElasticPool, name: <that resource's name>, fieldPath: status.outputs.elastic_pool_id}} -- a bare string does not parse

### spec.maxSizeGb

`double` · optional (explicit presence)

The database's storage ceiling in gigabytes (0.5-4096 for
vCore/DTU; Hyperscale grows elastically and ignores it). Fractional
values are legal ARM sizes (Basic tops out at 2, S0 at 250). Unset
applies the SKU's default.

- rule: {"double":{"lte":4096,"gte":0.1}}

### spec.collation

`string` · required · optional (explicit presence)

The database collation (sort order and string comparison).
Unspecified applies "SQL_Latin1_General_CP1_CI_AS" -- Azure SQL's
default. Changing collation replaces the database.

- default: `SQL_Latin1_General_CP1_CI_AS`
- rule: {"string":{"minLen":"1"}}

### spec.licenseType

`enum`

Azure Hybrid Benefit: BASE_PRICE brings your own SQL Server license
(with Software Assurance) for up to 55% savings; LICENSE_INCLUDED
pays as you go. vCore tiers only. Unset lets Azure default
(LicenseIncluded).

Allowed values (use exactly as shown):

- `azure_mssql_database_license_type_unspecified` -- Not specified: Azure defaults to LicenseIncluded.
- `BASE_PRICE` -- Bring your own SQL Server license with Software Assurance.
- `LICENSE_INCLUDED` -- Pay-as-you-go, license included in the hourly rate.

### spec.autoPauseDelayInMinutes

`int32` · optional (explicit presence)

Serverless only: minutes of inactivity before the database
auto-pauses (60-10080), or -1 to disable auto-pause. Requires a
GP_S_/HS_S_ sku.

- rule: auto_pause_delay_in_minutes must be -1 (disabled) or 60-10080

### spec.minCapacity

`double` · optional (explicit presence)

Serverless only: the minimum vCores the database always keeps warm
(0.25-40, fractional allowed). Requires a GP_S_/HS_S_ sku.

- rule: {"double":{"lte":40,"gte":0.25}}

### spec.readReplicaCount

`int32` · optional (explicit presence)

Hyperscale only: how many readable HA replicas back the database
(0-4). Each replica serves read-intent connections and shortens
failover.

- rule: {"int32":{"lte":4,"gte":0}}

### spec.readScale

`bool`

Premium / Business Critical only: route read-intent connections
(ApplicationIntent=ReadOnly) to a readable secondary replica.

### spec.zoneRedundant

`bool`

Spread the database's replicas across availability zones. Supported
on Premium, Business Critical, Hyperscale, and zone-capable General
Purpose regions.

### spec.ledgerEnabled

`bool`

Ledger: every table becomes cryptographically verifiable
(tamper-evident). Fixed at creation.

### spec.enclaveType

`enum`

The confidential-computing enclave databases run queries in. VBS
enables Always Encrypted with secure enclaves. Changing the enclave
type replaces the database; Hyperscale added VBS support only for
newly created databases.

Allowed values (use exactly as shown):

- `azure_mssql_database_enclave_type_unspecified` -- Not specified: no enclave configured.
- `VBS` -- Virtualization-based security enclaves (Always Encrypted with secure enclaves).
- `DEFAULT_ENCLAVE` -- ARM's explicit "Default" (no enclave) -- distinct from unspecified so an update can actively clear a previously set enclave.

### spec.maintenanceConfigurationName

`string`

The maintenance window Azure patches the database in (e.g.
"SQL_Default", "SQL_EastUS_DB_1", "SQL_EastUS_DB_2"). Region-specific
vocabularies; unset applies SQL_Default. A pooled database inherits
the pool's window and must leave this unset.

### spec.createMode

`enum`

How the database comes into existence. Unspecified means DEFAULT (a
fresh, empty database). Each mode consumes its matching source
field; all fixed at creation.

Allowed values (use exactly as shown):

- `azure_mssql_database_create_mode_unspecified` -- Not specified: DEFAULT (a fresh, empty database).
- `DEFAULT` -- A fresh, empty database.
- `COPY` -- A transactionally consistent copy of the source database.
- `SECONDARY` -- A readable secondary replica of the source (geo or named, per secondary_type).
- `ONLINE_SECONDARY` -- A readable secondary created while the source stays online (ARM's OnlineSecondary variant).
- `POINT_IN_TIME_RESTORE` -- Restore the source database to restore_point_in_time.
- `RECOVERY` -- Recover from a geo-replicated backup (recover_database_id or recovery_point_id).
- `RESTORE` -- Restore a dropped database (restore_dropped_database_id).
- `RESTORE_LONG_TERM_RETENTION_BACKUP` -- Restore from a long-term-retention backup (restore_long_term_retention_backup_id).

### spec.creationSourceDatabaseId

`string | valueFrom`

For COPY, SECONDARY, ONLINE_SECONDARY, and POINT_IN_TIME_RESTORE:
the source database, by ARM ID -- references another
AzureMssqlDatabase's database_id output, so copy/DR topologies
compose in one manifest set. Fixed at creation.

- references: AzureMssqlDatabase (`status.outputs.database_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMssqlDatabase, name: <that resource's name>, fieldPath: status.outputs.database_id}} -- a bare string does not parse

### spec.secondaryType

`enum`

For SECONDARY/ONLINE_SECONDARY: GEO creates a geo-replicated DR
secondary; NAMED creates an in-region readable replica. Unset lets
Azure default to GEO.

Allowed values (use exactly as shown):

- `azure_mssql_database_secondary_type_unspecified` -- Not specified: Azure defaults to GEO.
- `GEO` -- A geo-replicated disaster-recovery secondary (usually cross-region).
- `NAMED` -- An in-region readable replica addressed by its own name.

### spec.restorePointInTime

`string`

For POINT_IN_TIME_RESTORE: the RFC-3339 UTC instant to restore the
source database to (must fall inside its retention window). Fixed at
creation.

- rule: restore_point_in_time must be an RFC-3339 UTC timestamp, e.g. 2026-07-01T08:30:00Z

### spec.recoverDatabaseId

`string`

For RECOVERY: the recoverable database (geo-replicated backup) to
recover, by ARM ID.

### spec.recoveryPointId

`string`

For RECOVERY: the geo-backup recovery point to recover from, by ARM
ID (alternative to recover_database_id).

### spec.restoreDroppedDatabaseId

`string`

For RESTORE: the DROPPED database to restore, by its restorable-
dropped-database ARM ID (see the server's restorable dropped
databases).

### spec.restoreLongTermRetentionBackupId

`string`

For RESTORE_LONG_TERM_RETENTION_BACKUP: the long-term-retention
backup to restore, by ARM ID.

### spec.storageAccountType

`enum`

Which storage the database's BACKUPS replicate to: GEO (paired
region -- the default), GEO_ZONE (zones + paired region), ZONE
(zones in-region), LOCAL (single copy in-region).

Allowed values (use exactly as shown):

- `azure_mssql_database_backup_storage_account_type_unspecified` -- Not specified: Azure defaults to geo-redundant.
- `GEO_REDUNDANT` -- Geo-redundant: backups replicate to the paired region (enables geo restore).
- `GEO_ZONE_REDUNDANT` -- Geo-zone-redundant: across zones AND to the paired region.
- `LOCALLY_REDUNDANT` -- Locally redundant: one region, one zone.
- `ZONE_REDUNDANT` -- Zone-redundant: across the region's availability zones.

### spec.geoBackupEnabled

`bool` · optional (explicit presence)

Data Warehouse (DW) SKUs only: whether geo-redundant backups are
taken. Azure's default is true; other tiers control backup
redundancy via storage_account_type.

- default: `true`

### spec.sampleName

`string`

Seed the database with a sample schema. The only sample Azure ships
is "AdventureWorksLT". Fixed at creation.

- rule: sample_name must be AdventureWorksLT (the only sample Azure ships)

### spec.userAssignedIdentityIds

`[]string | valueFrom`

User-assigned identities attached to the DATABASE (database-scoped
customer-managed keys unwrap through them), by ARM ID.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.transparentDataEncryptionEnabled

`bool` · optional (explicit presence)

Whether transparent data encryption is on. Azure's default is true;
only OFF for compliance regimes that bring their own file-level
encryption (Hyperscale cannot disable it).

- default: `true`

### spec.transparentDataEncryptionKeyVaultKeyId

`string | valueFrom`

Database-scoped TDE customer-managed key (overrides the server's
key for this database), by VERSIONED Key Vault key ID -- references
an AzureKeyVaultKey's key_id output. Requires an attached
user-assigned identity with wrap/unwrap access on the key's vault.

- references: AzureKeyVaultKey (`status.outputs.key_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.key_id}} -- a bare string does not parse

### spec.transparentDataEncryptionKeyAutomaticRotationEnabled

`bool`

Whether the database automatically re-encrypts with new versions of
the customer-managed key as it rotates. Requires the CMK. Azure's
default is false.

### spec.import

`AzureMssqlDatabaseImport`

Import a bacpac export into the database right after creation.
Incompatible with the non-default create modes.

### spec.import.storageUri

`string` · required

The blob URI of the .bacpac file (e.g.
https://account.blob.core.windows.net/bacpacs/app.bacpac).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.import.storageKey

`string | valueFrom` · required · sensitive

The storage credential that reads the bacpac: the account's access
key or a SAS token, per storage_key_type.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.import.storageKeyType

`enum` · required

What kind of credential storage_key carries.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_mssql_database_import_storage_key_type_unspecified` -- Not specified -- invalid; choose an explicit type.
- `SHARED_ACCESS_KEY` -- A shared-access-signature (SAS) token.
- `STORAGE_ACCESS_KEY` -- The storage account's access key.

### spec.import.administratorLogin

`string` · required

The login the import runs as -- the server's SQL administrator or an
Entra login, per authentication_type.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.import.administratorLoginPassword

`string | valueFrom` · required · sensitive

The password for administrator_login.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.import.authenticationType

`enum` · required

How the import authenticates to the server.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_mssql_database_import_authentication_type_unspecified` -- Not specified -- invalid; choose an explicit type.
- `SQL` -- SQL authentication (the server's SQL administrator).
- `AD_PASSWORD` -- Microsoft Entra (Azure AD) password authentication.

### spec.import.storageAccountId

`string`

The storage account hosting the bacpac, by ARM ID -- required when
the account is behind a firewall or VNet rules so the import service
can negotiate access.

### spec.shortTermRetentionPolicy

`AzureMssqlDatabaseShortTermRetentionPolicy`

Point-in-time-restore horizon: how many days of automatic backups
are kept and how often differential backups are taken. Unset applies
Azure's defaults (7 days, 12-hour differentials).

### spec.shortTermRetentionPolicy.retentionDays

`int32` · required

How many days automatic backups are retained (1-35) -- the PITR
window.

- rule: {"required":true,"int32":{"lte":35,"gte":1}}

### spec.shortTermRetentionPolicy.backupIntervalInHours

`int32` · optional (explicit presence)

Hours between differential backups: 12 or 24. Unset applies Azure's
default of 12.

- rule: backup_interval_in_hours must be 12 or 24

### spec.longTermRetentionPolicy

`AzureMssqlDatabaseLongTermRetentionPolicy`

Long-term backup retention: weekly/monthly/yearly full backups kept
beyond the PITR window, each as an ISO-8601 duration (e.g. "P1M",
"P5W", "P1Y"). Unset keeps no long-term backups.

- rule: at least one of weekly_retention, monthly_retention, or yearly_retention must be set

### spec.longTermRetentionPolicy.weeklyRetention

`string`

How long weekly full backups are kept, e.g. "P5W", "P1M". At least
one horizon must be non-zero.

- rule: weekly_retention must be an ISO-8601 duration, e.g. P5W

### spec.longTermRetentionPolicy.monthlyRetention

`string`

How long monthly full backups are kept, e.g. "P12M", "P1Y".

- rule: monthly_retention must be an ISO-8601 duration, e.g. P12M

### spec.longTermRetentionPolicy.yearlyRetention

`string`

How long yearly full backups are kept, e.g. "P5Y".

- rule: yearly_retention must be an ISO-8601 duration, e.g. P5Y

### spec.longTermRetentionPolicy.weekOfYear

`int32` · optional (explicit presence)

Which ISO week's full backup becomes the yearly backup (1-52).
Meaningful only with yearly_retention.

- rule: {"int32":{"lte":52,"gte":1}}

### spec.threatDetectionPolicy

`AzureMssqlDatabaseThreatDetectionPolicy`

Microsoft Defender threat detection scoped to this database
(overrides the server-scope policy). Unset inherits the server's
posture.

- rule: storage_endpoint and storage_account_access_key must be set together

### spec.threatDetectionPolicy.state

`enum` · required

Whether the database-scope policy is enforced.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_mssql_database_threat_detection_state_unspecified` -- Not specified -- invalid; choose an explicit state when the policy block is present.
- `ENABLED` -- Threat detection is on for this database.
- `DISABLED` -- Threat detection is configured but off.

### spec.threatDetectionPolicy.disabledAlerts

`[]string`

Alert classes to suppress (e.g. "Sql_Injection", "Access_Anomaly",
"Data_Exfiltration", "Unsafe_Action") -- ARM's wire vocabulary,
matching the server-scope policy's detectors.

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.threatDetectionPolicy.emailAccountAdmins

`bool`

Whether alert emails also go to the subscription administrators.

### spec.threatDetectionPolicy.emailAddresses

`[]string`

Additional email addresses that receive alerts.

- rule: {"repeated":{"items":{"string":{"minLen":"3"}}}}

### spec.threatDetectionPolicy.retentionDays

`int32` · optional (explicit presence)

How many days threat-detection audit records are retained in the
export storage account. 0 retains indefinitely.

- rule: {"int32":{"gte":0}}

### spec.threatDetectionPolicy.storageEndpoint

`string`

The blob-storage endpoint threat-detection audit records are
exported to. Set together with storage_account_access_key.

### spec.threatDetectionPolicy.storageAccountAccessKey

`string | valueFrom` · sensitive

The access key of the export storage account.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the database, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins.

## Validation Rules

- `mssql_db_pool_requires_elasticpool_sku`: a pooled database (elastic_pool_id set) must set sku_name to "ElasticPool"; a standalone database must not
- `mssql_db_pool_forbids_maintenance`: a pooled database inherits the pool's maintenance window -- maintenance_configuration_name must be unset when elastic_pool_id is set
- `mssql_db_serverless_dials_require_serverless_sku`: auto_pause_delay_in_minutes and min_capacity require a serverless sku (GP_S_/HS_S_)
- `mssql_db_read_replicas_require_hyperscale`: read_replica_count requires a Hyperscale sku (HS_)
- `mssql_db_source_matches_create_mode`: creation_source_database_id is required for COPY, SECONDARY, ONLINE_SECONDARY, and POINT_IN_TIME_RESTORE, and must be omitted for other create modes
- `mssql_db_restore_time_matches_create_mode`: restore_point_in_time is required for POINT_IN_TIME_RESTORE and must be omitted for other create modes
- `mssql_db_secondary_type_matches_create_mode`: secondary_type is only valid for SECONDARY and ONLINE_SECONDARY create modes
- `mssql_db_recovery_source_matches_create_mode`: recover_database_id/recovery_point_id require RECOVERY, restore_dropped_database_id requires RESTORE, and restore_long_term_retention_backup_id requires RESTORE_LONG_TERM_RETENTION_BACKUP
- `mssql_db_cmk_requires_identity`: transparent_data_encryption_key_vault_key_id requires at least one entry in user_assigned_identity_ids
- `mssql_db_cmk_rotation_requires_key`: transparent_data_encryption_key_automatic_rotation_enabled requires transparent_data_encryption_key_vault_key_id
- `mssql_db_import_requires_default_mode`: import (bacpac) is only valid on a fresh (DEFAULT) database

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMssqlDatabase, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.database_id` | `string` | The Azure Resource Manager ID of the database. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Sql/servers/{server}/databases/{name} Referenced by copy/secondary/restore databases (creation_source_database_id). |
| `status.outputs.database_name` | `string` | The name of the database -- the Database= segment of connection strings against the server's fqdn. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverId` | AzureMssqlServer | `status.outputs.server_id` |
| `spec.elasticPoolId` | AzureMssqlElasticPool | `status.outputs.elastic_pool_id` |
| `spec.creationSourceDatabaseId` | AzureMssqlDatabase | `status.outputs.database_id` |
| `spec.userAssignedIdentityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.transparentDataEncryptionKeyVaultKeyId` | AzureKeyVaultKey | `status.outputs.key_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMssqlDatabase | `spec.creationSourceDatabaseId` | `status.outputs.database_id` |
| AzureMssqlFailoverGroup | `spec.databaseIds` | `status.outputs.database_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
