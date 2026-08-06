# OciMysqlDbSystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciMysqlDbSystemSpec defines the specification for an Oracle Cloud
Infrastructure MySQL HeatWave Database System -- a fully managed MySQL
database service with integrated HeatWave in-memory analytics acceleration.

MySQL HeatWave runs on dedicated compute shapes in a chosen availability
domain and subnet. High Availability mode provisions three instances across
fault domains with automatic failover.

This component manages the DB System resource itself. HeatWave cluster
and replication channels are separate OCI resources with independent
lifecycles and are not bundled here.

Excluded from v1:
  - source block (BACKUP, PITR, IMPORTURL) -- only fresh creation supported
  - shutdown_type, state -- operational lifecycle controls
  - access_mode, database_mode -- runtime toggles
  - security_attributes -- ZPR artifacts
  - backup_policy.copy_policies -- cross-region backup copy
  - backup_policy.soft_delete -- newer feature
  - maintenance.maintenance_disabled_windows -- advanced scheduling

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.shapeName` | `string` | yes |  |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.adminUsername` | `string` |  |  |  |
| `spec.adminPassword` | `string` (sensitive) |  |  |  |
| `spec.mysqlVersion` | `string` |  |  |  |
| `spec.configurationId` | `string \| valueFrom` |  |  |  |
| `spec.isHighlyAvailable` | `bool` |  |  |  |
| `spec.hostnameLabel` | `string` |  |  |  |
| `spec.ipAddress` | `string` |  |  |  |
| `spec.faultDomain` | `string` |  |  |  |
| `spec.port` | `int32` |  |  |  |
| `spec.portX` | `int32` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.crashRecovery` | `string` |  |  |  |
| `spec.databaseManagement` | `string` |  |  |  |
| `spec.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.dataStorage` | `DataStorage` |  |  |  |
| `spec.dataStorage.dataStorageSizeInGb` | `int32` |  |  |  |
| `spec.dataStorage.isAutoExpandStorageEnabled` | `bool` |  |  |  |
| `spec.dataStorage.maxStorageSizeInGbs` | `int32` |  |  |  |
| `spec.backupPolicy` | `BackupPolicy` |  |  |  |
| `spec.backupPolicy.isEnabled` | `bool` |  |  |  |
| `spec.backupPolicy.retentionInDays` | `int32` |  |  |  |
| `spec.backupPolicy.windowStartTime` | `string` |  |  |  |
| `spec.backupPolicy.pitrPolicy` | `PitrPolicy` |  |  |  |
| `spec.backupPolicy.pitrPolicy.isEnabled` | `bool` |  |  |  |
| `spec.maintenance` | `Maintenance` |  |  |  |
| `spec.maintenance.windowStartTime` | `string` | yes |  |  |
| `spec.maintenance.maintenanceScheduleType` | `enum` |  |  |  |
| `spec.maintenance.versionPreference` | `enum` |  |  |  |
| `spec.maintenance.versionTrackPreference` | `enum` |  |  |  |
| `spec.deletionPolicy` | `DeletionPolicy` |  |  |  |
| `spec.deletionPolicy.automaticBackupRetention` | `string` |  |  |  |
| `spec.deletionPolicy.finalBackup` | `string` |  |  |  |
| `spec.deletionPolicy.isDeleteProtected` | `bool` |  |  |  |
| `spec.encryptData` | `EncryptData` |  |  |  |
| `spec.encryptData.keyGenerationType` | `enum` |  |  |  |
| `spec.encryptData.keyId` | `string \| valueFrom` |  |  |  |
| `spec.secureConnections` | `SecureConnections` |  |  |  |
| `spec.secureConnections.certificateGenerationType` | `enum` |  |  |  |
| `spec.secureConnections.certificateId` | `string \| valueFrom` |  |  |  |
| `spec.customerContacts` | `[]CustomerContact` |  |  |  |
| `spec.customerContacts[].email` | `string` | yes |  |  |
| `spec.readEndpoint` | `ReadEndpoint` |  |  |  |
| `spec.readEndpoint.isEnabled` | `bool` |  |  |  |
| `spec.readEndpoint.excludeIps` | `[]string` |  |  |  |
| `spec.readEndpoint.readEndpointHostnameLabel` | `string` |  |  |  |
| `spec.readEndpoint.readEndpointIpAddress` | `string` |  |  |  |
| `spec.databaseConsole` | `DatabaseConsole` |  |  |  |
| `spec.databaseConsole.status` | `enum` |  |  |  |
| `spec.databaseConsole.port` | `int32` |  |  |  |
| `spec.rest` | `Rest` |  |  |  |
| `spec.rest.configuration` | `string` |  |  |  |
| `spec.rest.port` | `int32` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the MySQL DB System will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.availabilityDomain

`string` · required

Availability domain for the primary (read/write) endpoint.
Example: "Uocm:PHX-AD-1". Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.shapeName

`string` · required

Compute shape for the DB System. Determines CPU, memory, and network
bandwidth. Examples: "MySQL.VM.Standard.E4.1.8GB",
"MySQL.VM.Standard.E4.4.64GB", "MySQL.HeatWave.VM.Standard".

- rule: {"string":{"minLen":"1"}}

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the DB System will be placed.
Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.adminUsername

`string`

Administrative username for the database. Changing this forces recreation.
When omitted, the provider default applies.

### spec.adminPassword

`string` · sensitive

Administrative password for the database. Must be 8-32 characters and
contain at least one numeric, one lowercase, one uppercase, and one
special character. Changing this forces recreation.

### spec.mysqlVersion

`string`

MySQL version identifier (e.g. "8.0.36", "9.1.0").
When omitted, the latest available version is used.
Changing this forces recreation.

### spec.configurationId

`string | valueFrom`

OCID of a MySQL Configuration to apply to the DB System.
Configurations define MySQL server variable settings (buffer pool size,
max connections, etc.). When omitted, the default configuration for
the selected shape is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.isHighlyAvailable

`bool` · optional (explicit presence)

When true, enables High Availability. Three instances are provisioned
across different fault domains with automatic failover. Standby
instances are not directly accessible.

### spec.hostnameLabel

`string`

Hostname for the primary endpoint. Combined with the subnet's DNS
domain to form the FQDN for database access.

### spec.ipAddress

`string`

Specific private IP address for the primary endpoint within the
subnet. When omitted, OCI auto-assigns an available IP.
Changing this forces recreation.

### spec.faultDomain

`string`

Fault domain for the primary (read/write) endpoint.
Example: "FAULT-DOMAIN-1". Changing this forces recreation.

### spec.port

`int32`

TCP port for the MySQL protocol on the primary endpoint.
Default: 3306. Changing this forces recreation.

### spec.portX

`int32`

TCP port for the X Protocol (MySQL Shell, connectors) on the
primary endpoint. Default: 33060. Changing this forces recreation.

### spec.description

`string`

User-provided description of the DB System.

### spec.crashRecovery

`string`

Controls InnoDB crash recovery mechanisms (Redo Logs, Double Write
Buffer, Binary Log syncing). Disabling improves write performance
but risks data loss on unexpected failure. Values: "ENABLED",
"DISABLED".

### spec.databaseManagement

`string`

Enables monitoring via the OCI Database Management service.
Values: "ENABLED", "DISABLED".

### spec.nsgIds

`[]string | valueFrom`

OCIDs of network security groups for the DB System VNIC.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.dataStorage

`DataStorage`

Data storage configuration. Uses the modern data_storage block
(the legacy top-level data_storage_size_in_gb is deprecated).

### spec.dataStorage.dataStorageSizeInGb

`int32`

Initial data volume size in gigabytes. Supported values depend on
the shape (typically 50 GB minimum).

### spec.dataStorage.isAutoExpandStorageEnabled

`bool` · optional (explicit presence)

When true, storage automatically expands when usage nears the limit.

### spec.dataStorage.maxStorageSizeInGbs

`int32`

Maximum storage size in gigabytes for auto-expansion.
Only effective when is_auto_expand_storage_enabled is true.
Range: 32768-131072 depending on initial size.

### spec.backupPolicy

`BackupPolicy`

Automatic backup configuration.

### spec.backupPolicy.isEnabled

`bool` · optional (explicit presence)

Whether automatic backups are enabled.

### spec.backupPolicy.retentionInDays

`int32`

Number of days to retain automatic backups.

### spec.backupPolicy.windowStartTime

`string`

Start of the 30-minute daily backup window in RFC3339 time format
(e.g. "03:00"). When omitted, OCI selects the window.

### spec.backupPolicy.pitrPolicy

`PitrPolicy`

Point-in-time recovery configuration.

### spec.backupPolicy.pitrPolicy.isEnabled

`bool` · optional (explicit presence)

Whether point-in-time recovery is enabled. Requires automatic
backups to be enabled.

### spec.maintenance

`Maintenance`

Maintenance window configuration.

### spec.maintenance.windowStartTime

`string` · required

Start of the maintenance window. Format: "{day-of-week} {time-of-day}"
(e.g. "mon 10:00"). Required when maintenance is configured.

- rule: {"string":{"minLen":"1"}}

### spec.maintenance.maintenanceScheduleType

`enum`

Maintenance schedule type.

Allowed values (use exactly as shown):

- `maintenance_schedule_type_unspecified`
- `early` -- Receive patches earlier than the regular schedule.
- `regular` -- Receive patches on the standard Oracle schedule.

### spec.maintenance.versionPreference

`enum`

Version preference for automatic upgrades.

Allowed values (use exactly as shown):

- `version_preference_unspecified`
- `oldest` -- Keep the oldest supported version as long as possible.
- `second_newest` -- Use the second-newest available version.
- `newest` -- Always use the newest available version.

### spec.maintenance.versionTrackPreference

`enum`

Version track preference controlling which release stream to follow.

Allowed values (use exactly as shown):

- `version_track_preference_unspecified`
- `long_term_support` -- Long-Term Support releases only.
- `innovation` -- Innovation releases (frequent, feature-rich).
- `follow` -- Follow the OCI-recommended track (default behavior).

### spec.deletionPolicy

`DeletionPolicy`

Deletion safety configuration.

### spec.deletionPolicy.automaticBackupRetention

`string`

What to do with automatic backups on deletion.
Values: "DELETE", "RETAIN".

### spec.deletionPolicy.finalBackup

`string`

Whether to create a final backup before deletion.
Values: "REQUIRE_FINAL_BACKUP", "SKIP_FINAL_BACKUP".

### spec.deletionPolicy.isDeleteProtected

`bool` · optional (explicit presence)

When true, the DB System cannot be deleted. Must be set to false
before deletion is possible.

### spec.encryptData

`EncryptData`

Data-at-rest encryption configuration.

### spec.encryptData.keyGenerationType

`enum`

Key generation strategy. "system" uses Oracle-managed keys.
"byok" (Bring Your Own Key) requires key_id to be set.

Allowed values (use exactly as shown):

- `key_generation_type_unspecified`
- `system` -- Oracle-managed encryption keys.
- `byok` -- Customer-managed encryption keys (Bring Your Own Key).

### spec.encryptData.keyId

`string | valueFrom`

OCID of the customer-managed encryption key. Required when
key_generation_type is byok.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.secureConnections

`SecureConnections`

TLS certificate configuration for client connections.

### spec.secureConnections.certificateGenerationType

`enum`

Certificate generation strategy. "system" uses Oracle-managed
certificates. "byoc" (Bring Your Own Certificate) requires
certificate_id to be set.

Allowed values (use exactly as shown):

- `certificate_generation_type_unspecified`
- `system_cert` -- Oracle-managed TLS certificates.
- `byoc` -- Customer-managed TLS certificates (Bring Your Own Certificate).

### spec.secureConnections.certificateId

`string | valueFrom`

OCID of the customer-managed certificate. Required when
certificate_generation_type is byoc.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.customerContacts

`[]CustomerContact`

Customer contact email addresses for operational notifications
(maintenance windows, critical alerts). Maximum 10 contacts.

- rule: {"repeated":{"maxItems":"10"}}

### spec.customerContacts[].email

`string` · required

Email address for operational notifications.

- rule: {"string":{"minLen":"1"}}

### spec.readEndpoint

`ReadEndpoint`

Read-only endpoint configuration for read scaling.

### spec.readEndpoint.isEnabled

`bool` · optional (explicit presence)

Whether the read endpoint is enabled.

### spec.readEndpoint.excludeIps

`[]string`

IP addresses to exclude from serving read requests.

### spec.readEndpoint.readEndpointHostnameLabel

`string`

Hostname for the read endpoint. Combined with the subnet's DNS
domain to form the read endpoint FQDN.

### spec.readEndpoint.readEndpointIpAddress

`string`

Specific private IP address for the read endpoint.
When omitted, OCI auto-assigns an available IP.

### spec.databaseConsole

`DatabaseConsole`

MySQL Database Console (web-based management UI) configuration.

### spec.databaseConsole.status

`enum`

Whether the database console is enabled or disabled.

Allowed values (use exactly as shown):

- `database_console_status_unspecified`
- `enabled` -- Database console is enabled.
- `disabled` -- Database console is disabled.

### spec.databaseConsole.port

`int32`

Port for the database console. Valid values: 443 or 1024-65535.

### spec.rest

`Rest`

MySQL REST API service configuration.

### spec.rest.configuration

`string`

REST API configuration mode.

### spec.rest.port

`int32`

Port for the REST API service. Valid values: 443 or 1024-65535.

## Validation Rules

- `encrypt_data_key_required`: encrypt_data.key_id is required when key_generation_type is byok
- `secure_connections_cert_required`: secure_connections.certificate_id is required when certificate_generation_type is byoc

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciMysqlDbSystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.db_system_id` | `string` | OCID of the MySQL DB System. |
| `status.outputs.endpoint_hostname` | `string` | Hostname of the primary (read/write) endpoint. |
| `status.outputs.endpoint_ip_address` | `string` | Private IP address of the primary (read/write) endpoint. |
| `status.outputs.endpoint_port` | `int32` | TCP port of the primary (read/write) endpoint. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
