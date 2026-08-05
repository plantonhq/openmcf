# OciDbSystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciDbSystemSpec defines the specification for an Oracle Cloud Infrastructure
Database System -- a managed Oracle Database running on Virtual Machine or
Bare Metal infrastructure.

A DB System provisions the underlying compute and storage, a DB Home
containing the Oracle Database software, and an initial database instance.
These three layers (system, home, database) form an inseparable unit at
creation time and are modeled as nested messages rather than separate
resources.

This component supports fresh creation only (source=NONE). Clone and
restore scenarios (source=DB_BACKUP, DATABASE, DB_SYSTEM) require
different field sets and may be added in a future version.

Deprecated/excluded fields:
  - db_workload: deprecated by Oracle since November 2023
  - compute_model/compute_count: Exadata Cloud@Customer specific; VM/BM
    shapes use cpu_core_count
  - source and clone/restore fields: different operational workflow
  - backup_destination_details: advanced backup routing (DBRS default
    covers most cases)

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.shape` | `string` | yes |  |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.sshPublicKeys` | `[]string` | yes |  |  |
| `spec.hostname` | `string` | yes |  |  |
| `spec.cpuCoreCount` | `int32` |  |  |  |
| `spec.databaseEdition` | `enum` |  |  |  |
| `spec.licenseModel` | `enum` |  |  |  |
| `spec.dataStorageSizeInGb` | `int32` |  |  |  |
| `spec.dataStoragePercentage` | `int32` |  |  |  |
| `spec.diskRedundancy` | `enum` |  |  |  |
| `spec.nodeCount` | `int32` |  |  |  |
| `spec.domain` | `string` |  |  |  |
| `spec.clusterName` | `string` |  |  |  |
| `spec.faultDomains` | `[]string` |  |  |  |
| `spec.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.backupSubnetId` | `string \| valueFrom` |  |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.backupNetworkNsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.kmsKeyVersionId` | `string` |  |  |  |
| `spec.timeZone` | `string` |  |  |  |
| `spec.sparseDiskgroup` | `bool` |  |  |  |
| `spec.storageVolumePerformanceMode` | `enum` |  |  |  |
| `spec.privateIp` | `string` |  |  |  |
| `spec.dataCollectionOptions` | `DataCollectionOptions` |  |  |  |
| `spec.dataCollectionOptions.isDiagnosticsEventsEnabled` | `bool` |  |  |  |
| `spec.dataCollectionOptions.isHealthMonitoringEnabled` | `bool` |  |  |  |
| `spec.dataCollectionOptions.isIncidentLogsEnabled` | `bool` |  |  |  |
| `spec.dbSystemOptions` | `DbSystemOptions` |  |  |  |
| `spec.dbSystemOptions.storageManagement` | `enum` |  |  |  |
| `spec.maintenanceWindowDetails` | `MaintenanceWindowDetails` |  |  |  |
| `spec.maintenanceWindowDetails.preference` | `enum` |  |  |  |
| `spec.maintenanceWindowDetails.patchingMode` | `enum` |  |  |  |
| `spec.maintenanceWindowDetails.leadTimeInWeeks` | `int32` |  |  |  |
| `spec.maintenanceWindowDetails.months` | `[]string` |  |  |  |
| `spec.maintenanceWindowDetails.weeksOfMonth` | `[]int32` |  |  |  |
| `spec.maintenanceWindowDetails.daysOfWeek` | `[]string` |  |  |  |
| `spec.maintenanceWindowDetails.hoursOfDay` | `[]int32` |  |  |  |
| `spec.maintenanceWindowDetails.customActionTimeoutInMins` | `int32` |  |  |  |
| `spec.maintenanceWindowDetails.isCustomActionTimeoutEnabled` | `bool` |  |  |  |
| `spec.maintenanceWindowDetails.isMonthlyPatchingEnabled` | `bool` |  |  |  |
| `spec.dbHome` | `DbHome` | yes |  |  |
| `spec.dbHome.dbVersion` | `string` |  |  |  |
| `spec.dbHome.displayName` | `string` |  |  |  |
| `spec.dbHome.databaseSoftwareImageId` | `string \| valueFrom` |  |  |  |
| `spec.dbHome.database` | `Database` | yes |  |  |
| `spec.dbHome.database.adminPassword` | `string` (sensitive) | yes |  |  |
| `spec.dbHome.database.dbName` | `string` | yes |  |  |
| `spec.dbHome.database.characterSet` | `string` |  |  |  |
| `spec.dbHome.database.ncharacterSet` | `string` |  |  |  |
| `spec.dbHome.database.pdbName` | `string` |  |  |  |
| `spec.dbHome.database.dbDomain` | `string` |  |  |  |
| `spec.dbHome.database.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.dbHome.database.kmsKeyVersionId` | `string` |  |  |  |
| `spec.dbHome.database.vaultId` | `string \| valueFrom` |  |  |  |
| `spec.dbHome.database.dbBackupConfig` | `DbBackupConfig` |  |  |  |
| `spec.dbHome.database.dbBackupConfig.autoBackupEnabled` | `bool` |  |  |  |
| `spec.dbHome.database.dbBackupConfig.autoBackupWindow` | `string` |  |  |  |
| `spec.dbHome.database.dbBackupConfig.recoveryWindowInDays` | `int32` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the DB System will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name for the DB System shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.availabilityDomain

`string` · required

Availability domain where the DB System will be placed.
Example: "Uocm:PHX-AD-1". Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.shape

`string` · required

Compute shape for the DB System nodes. Determines CPU architecture,
core count range, and memory. Examples: "VM.Standard2.4",
"VM.Standard.E4.Flex", "BM.DenseIO2.52".

- rule: {"string":{"minLen":"1"}}

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the DB System will be placed.
Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.sshPublicKeys

`[]string` · required

SSH public keys for administrative access to the DB System nodes.
At least one key is required. Keys must be in OpenSSH format.

- rule: {"repeated":{"minItems":"1"}}

### spec.hostname

`string` · required

Hostname for the DB System. Must be unique within the subnet.
Combined with the domain to form the FQDN. Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.cpuCoreCount

`int32`

Number of CPU cores to enable. For VM shapes, this is the number
of OCPUs. For BM shapes, this is the total core count. When omitted,
the shape's default is used.

### spec.databaseEdition

`enum`

Oracle Database edition. Determines available features and licensing
tiers. Exadata and 2-node RAC require enterprise_edition_extreme_performance.
Changing this forces recreation.

Allowed values (use exactly as shown):

- `database_edition_unspecified`
- `standard_edition`
- `enterprise_edition`
- `enterprise_edition_high_performance`
- `enterprise_edition_extreme_performance`

### spec.licenseModel

`enum`

License model for the DB System. Determines whether existing Oracle
licenses are applied (BYOL) or new licenses are included in the price.

Allowed values (use exactly as shown):

- `license_model_unspecified`
- `bring_your_own_license`
- `license_included`

### spec.dataStorageSizeInGb

`int32`

Initial data storage size in gigabytes. Required for VM DB Systems.
Supported values depend on the shape (typically 256, 512, 1024, 2048,
4096, 6144, 8192, 10240, 12288, 14336, 16384, 20480, 24576, 28672,
32768, 36864, 40960).

### spec.dataStoragePercentage

`int32`

Percentage of total storage allocated to data (as opposed to recovery).
Valid values: 40 or 80. Only applicable for BM DB Systems.
Changing this forces recreation.

### spec.diskRedundancy

`enum`

Disk redundancy level. "normal" provides 2-way mirroring; "high"
provides 3-way mirroring. Only applicable for BM DB Systems.
Changing this forces recreation.

Allowed values (use exactly as shown):

- `disk_redundancy_unspecified`
- `normal`
- `high`

### spec.nodeCount

`int32`

Number of DB System nodes. Use 1 for a single-node system or 2 for
a 2-node RAC cluster. Changing this forces recreation.

### spec.domain

`string`

Network domain name for the DB System. Defaults to the subnet's
domain if not specified. Changing this forces recreation.

### spec.clusterName

`string`

Cluster name for the DB System. Used for RAC clusters.
Maximum 11 characters. Changing this forces recreation.

### spec.faultDomains

`[]string`

Fault domains for distributing DB System nodes. Use for RAC to
place nodes in different fault domains for high availability.
Changing this forces recreation.

### spec.nsgIds

`[]string | valueFrom`

OCIDs of network security groups to apply to the DB System VNIC.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.backupSubnetId

`string | valueFrom`

OCID of the backup subnet for the DB System. Required for RAC
configurations. Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.backupNetworkNsgIds

`[]string | valueFrom`

OCIDs of network security groups for the backup VNIC.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.kmsKeyId

`string | valueFrom`

OCID of the KMS key for encrypting the DB System's data at rest.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.kmsKeyVersionId

`string`

OCID of the specific KMS key version. When omitted, the current
key version is used.

### spec.timeZone

`string`

Time zone for the DB System. Example: "UTC", "US/Pacific".
Changing this forces recreation.

### spec.sparseDiskgroup

`bool` · optional (explicit presence)

When true, configures a sparse disk group on the BM DB System.
Only applicable for BM shapes. Changing this forces recreation.

### spec.storageVolumePerformanceMode

`enum`

Storage volume performance mode. "balanced" is suitable for most
workloads; "high_performance" provides lower latency at higher cost.
Changing this forces recreation.

Allowed values (use exactly as shown):

- `storage_volume_performance_mode_unspecified`
- `balanced`
- `high_performance`

### spec.privateIp

`string`

Specific private IP address for the DB System. When omitted,
OCI auto-assigns an available IP from the subnet.
Changing this forces recreation.

### spec.dataCollectionOptions

`DataCollectionOptions`

Data collection options controlling OCI diagnostic telemetry.

### spec.dataCollectionOptions.isDiagnosticsEventsEnabled

`bool` · optional (explicit presence)

Whether to enable diagnostic event collection.

### spec.dataCollectionOptions.isHealthMonitoringEnabled

`bool` · optional (explicit presence)

Whether to enable health monitoring.

### spec.dataCollectionOptions.isIncidentLogsEnabled

`bool` · optional (explicit presence)

Whether to enable incident log collection.

### spec.dbSystemOptions

`DbSystemOptions`

DB System options controlling storage management strategy.

### spec.dbSystemOptions.storageManagement

`enum`

Storage management strategy. ASM (Automatic Storage Management) is
the default for most shapes. LVM (Logical Volume Management) is
available for single-node VM systems only.
Changing this forces recreation.

Allowed values (use exactly as shown):

- `storage_management_unspecified`
- `asm`
- `lvm`

### spec.maintenanceWindowDetails

`MaintenanceWindowDetails`

Maintenance window scheduling. When omitted or preference is
no_preference, OCI selects the maintenance window automatically.

### spec.maintenanceWindowDetails.preference

`enum`

Maintenance scheduling preference.

Allowed values (use exactly as shown):

- `preference_unspecified`
- `no_preference`
- `custom_preference`

### spec.maintenanceWindowDetails.patchingMode

`enum`

Patching strategy. Rolling applies patches one node at a time
(zero downtime for RAC). Nonrolling patches all nodes simultaneously.

Allowed values (use exactly as shown):

- `patching_mode_unspecified`
- `rolling`
- `nonrolling`

### spec.maintenanceWindowDetails.leadTimeInWeeks

`int32`

Weeks of advance notice before the maintenance window opens.

### spec.maintenanceWindowDetails.months

`[]string`

Months when maintenance is allowed. Example: ["JANUARY", "APRIL",
"JULY", "OCTOBER"] for quarterly patching.

### spec.maintenanceWindowDetails.weeksOfMonth

`[]int32`

Weeks of the month (1-4) when maintenance is allowed.

### spec.maintenanceWindowDetails.daysOfWeek

`[]string`

Days of the week when maintenance is allowed. Example: ["MONDAY"].

### spec.maintenanceWindowDetails.hoursOfDay

`[]int32`

Hours of the day (0-23) when maintenance may start.

### spec.maintenanceWindowDetails.customActionTimeoutInMins

`int32`

Custom timeout for patching actions in minutes (0-120).

### spec.maintenanceWindowDetails.isCustomActionTimeoutEnabled

`bool` · optional (explicit presence)

Whether the custom action timeout is enabled.

### spec.maintenanceWindowDetails.isMonthlyPatchingEnabled

`bool` · optional (explicit presence)

Whether monthly patching is enabled.

### spec.dbHome

`DbHome` · required

DB Home configuration. Contains the Oracle Database software version
and the initial database to create. Exactly one DB Home is required
at DB System creation time.

- rule: {"required":true}

### spec.dbHome.dbVersion

`string`

Oracle Database version. Example: "19.0.0.0", "21.0.0.0".
Mutually exclusive with database_software_image_id.
Changing this forces recreation.

### spec.dbHome.displayName

`string`

Human-readable name for the DB Home.

### spec.dbHome.databaseSoftwareImageId

`string | valueFrom`

OCID of a custom database software image to use instead of a
standard db_version. Mutually exclusive with db_version.
Changing this forces recreation.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.dbHome.database

`Database` · required

The initial database to create within this DB Home.

- rule: {"required":true}

### spec.dbHome.database.adminPassword

`string` · required · sensitive

Administrator password for the SYS and SYSTEM users. Must be 2 to 30
characters, contain at least one uppercase, one lowercase, and one
numeric character. Cannot contain the double-quote character.
This value is not returned by the API after creation.

- rule: {"string":{"minLen":"2"}}

### spec.dbHome.database.dbName

`string` · required

Database name. Must be alphanumeric, begin with a letter, and be
at most 8 characters for single-node or 30 characters otherwise.
Changing this forces recreation.

- rule: {"string":{"minLen":"1","maxLen":"30","pattern":"^[a-zA-Z][a-zA-Z0-9]*$"}}

### spec.dbHome.database.characterSet

`string`

Character set for the database. Defaults to AL32UTF8 when omitted.
Changing this forces recreation.

### spec.dbHome.database.ncharacterSet

`string`

National character set. Defaults to AL16UTF16 when omitted.
Valid values: AL16UTF16, UTF8. Changing this forces recreation.

### spec.dbHome.database.pdbName

`string`

Pluggable database name. Must begin with a letter and contain only
alphanumeric characters. Changing this forces recreation.

### spec.dbHome.database.dbDomain

`string`

Database domain. Defaults to the DB System's domain when omitted.
Changing this forces recreation.

### spec.dbHome.database.kmsKeyId

`string | valueFrom`

OCID of the KMS key for Transparent Data Encryption (TDE).
When omitted, Oracle-managed encryption is used.
Changing this forces recreation.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.dbHome.database.kmsKeyVersionId

`string`

OCID of the specific KMS key version for TDE.
Changing this forces recreation.

### spec.dbHome.database.vaultId

`string | valueFrom`

OCID of the OCI Vault containing the TDE encryption key.
Required when kms_key_id is set. Changing this forces recreation.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.dbHome.database.dbBackupConfig

`DbBackupConfig`

Automatic backup configuration for the database.

### spec.dbHome.database.dbBackupConfig.autoBackupEnabled

`bool` · optional (explicit presence)

Whether automatic backups are enabled.

### spec.dbHome.database.dbBackupConfig.autoBackupWindow

`string`

Preferred backup window. Values: SLOT_ONE through SLOT_TWELVE
representing 2-hour windows starting from 00:00-02:00 UTC.
Only applied when auto_backup_enabled is true.

### spec.dbHome.database.dbBackupConfig.recoveryWindowInDays

`int32`

Number of days to retain automatic backups (1-60).
Only applied when auto_backup_enabled is true.

## Validation Rules

- `db_home_version_mutual_exclusivity`: db_home.db_version and db_home.database_software_image_id are mutually exclusive
- `db_home_version_required`: either db_home.db_version or db_home.database_software_image_id must be provided

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciDbSystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.db_system_id` | `string` | OCID of the DB System. |
| `status.outputs.db_home_id` | `string` | OCID of the first DB Home created with the DB System. |
| `status.outputs.database_id` | `string` | OCID of the initial database created within the first DB Home. |
| `status.outputs.listener_port` | `int32` | TCP port on which the database listener accepts connections. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |
| `spec.backupSubnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.backupNetworkNsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
