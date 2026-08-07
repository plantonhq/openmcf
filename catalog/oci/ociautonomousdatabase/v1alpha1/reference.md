# OciAutonomousDatabase

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciAutonomousDatabaseSpec defines the specification for an Oracle Cloud
Infrastructure Autonomous Database.

Autonomous Database is OCI's flagship fully managed database service
supporting multiple workload types:

  - **OLTP** (Autonomous Transaction Processing): mixed transactional workloads
  - **DW** (Autonomous Data Warehouse): analytic and reporting workloads
  - **AJD** (Autonomous JSON Database): JSON document storage and queries
  - **APEX**: Oracle APEX low-code application development
  - **LH** (Lakehouse): data lake analytics combining external and managed data

The database can run on shared (serverless) or dedicated Exadata
infrastructure. Serverless databases specify storage in terabytes;
dedicated databases may use gigabytes for finer granularity.

Two compute models are available: ECPU (recommended) and OCPU (legacy).
Auto-scaling can independently scale compute (up to 3x provisioned) and
storage on demand.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.dbName` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.dbWorkload` | `enum` |  |  |  |
| `spec.dbVersion` | `string` |  |  |  |
| `spec.databaseEdition` | `enum` |  |  |  |
| `spec.licenseModel` | `enum` |  |  |  |
| `spec.characterSet` | `string` |  |  |  |
| `spec.ncharacterSet` | `string` |  |  |  |
| `spec.computeModel` | `enum` |  |  |  |
| `spec.computeCount` | `float` |  |  |  |
| `spec.dataStorageSizeInTbs` | `int32` |  |  |  |
| `spec.dataStorageSizeInGb` | `int32` |  |  |  |
| `spec.isAutoScalingEnabled` | `bool` |  |  |  |
| `spec.isAutoScalingForStorageEnabled` | `bool` |  |  |  |
| `spec.adminPassword` | `string` (sensitive) |  |  |  |
| `spec.secretId` | `string \| valueFrom` |  |  |  |
| `spec.secretVersionNumber` | `int32` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` |  |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.privateEndpointLabel` | `string` |  |  |  |
| `spec.privateEndpointIp` | `string` |  |  |  |
| `spec.whitelistedIps` | `[]string` |  |  |  |
| `spec.isMtlsConnectionRequired` | `bool` |  |  |  |
| `spec.isAccessControlEnabled` | `bool` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  |  |
| `spec.vaultId` | `string \| valueFrom` |  |  |  |
| `spec.isDedicated` | `bool` |  |  |  |
| `spec.isFreeTier` | `bool` |  |  |  |
| `spec.isDevTier` | `bool` |  |  |  |
| `spec.autonomousContainerDatabaseId` | `string \| valueFrom` |  |  |  |
| `spec.backupRetentionPeriodInDays` | `int32` |  |  |  |
| `spec.isLocalDataGuardEnabled` | `bool` |  |  |  |
| `spec.autonomousMaintenanceScheduleType` | `enum` |  |  |  |
| `spec.customerContacts` | `[]CustomerContact` |  |  |  |
| `spec.customerContacts[].email` | `string` | yes |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the autonomous database will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.dbName

`string` · required

The database name. Must be alphanumeric, begin with a letter, and be
unique within the tenancy. Maximum 30 characters. Cannot be changed
after creation.

- rule: {"string":{"minLen":"1","maxLen":"30","pattern":"^[a-zA-Z][a-zA-Z0-9]*$"}}

### spec.displayName

`string`

Human-readable display name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.dbWorkload

`enum`

The workload type of the autonomous database. Determines internal
storage layout, query optimizer behavior, and available features.

Allowed values (use exactly as shown):

- `db_workload_unspecified`
- `oltp` -- Autonomous Transaction Processing -- optimized for mixed OLTP workloads.
- `dw` -- Autonomous Data Warehouse -- optimized for analytic and reporting workloads.
- `ajd` -- Autonomous JSON Database -- optimized for JSON document storage and queries.
- `apex` -- Oracle APEX Application Development -- optimized for low-code app development.
- `lh` -- Lakehouse -- optimized for data lake analytics combining external and managed data.

### spec.dbVersion

`string`

Oracle Database version (e.g. "19c", "23ai", "26ai").
When omitted, the latest available version is used.

### spec.databaseEdition

`enum`

The Oracle Database edition.

Allowed values (use exactly as shown):

- `database_edition_unspecified`
- `standard_edition` -- Standard Edition -- suitable for most workloads with core Oracle features.
- `enterprise_edition` -- Enterprise Edition -- includes advanced features like partitioning, compression, and advanced security.

### spec.licenseModel

`enum`

The license model. AJD and APEX workloads always use LICENSE_INCLUDED
regardless of this setting.

Allowed values (use exactly as shown):

- `license_model_unspecified`
- `bring_your_own_license` -- Bring Your Own License -- apply existing Oracle Database licenses.
- `license_included` -- License Included -- licensing cost is included in the service price.

### spec.characterSet

`string`

Character set for the database. Defaults to AL32UTF8 when omitted.
Cannot be changed after creation.

### spec.ncharacterSet

`string`

National character set. Defaults to AL16UTF16 when omitted.
Valid values: AL16UTF16, UTF8. Cannot be changed after creation.

### spec.computeModel

`enum`

The compute model. ECPU is Oracle's current recommended model.
OCPU is the legacy model still supported for backward compatibility.

Allowed values (use exactly as shown):

- `compute_model_unspecified`
- `ecpu` -- Elastic Compute Processing Units -- Oracle's current recommended model.
- `ocpu` -- Oracle Compute Processing Units -- legacy model, still supported.

### spec.computeCount

`float` · optional (explicit presence)

Number of compute units (ECPUs or OCPUs depending on compute_model).
Minimum 2 ECPUs for ECPU model. Minimum varies by workload for OCPU.

### spec.dataStorageSizeInTbs

`int32`

Maximum storage in terabytes. Used for serverless deployments.
Mutually exclusive with data_storage_size_in_gb.

### spec.dataStorageSizeInGb

`int32`

Maximum storage in gigabytes. Used for dedicated Exadata deployments
that need finer granularity than whole terabytes.
Mutually exclusive with data_storage_size_in_tbs.

### spec.isAutoScalingEnabled

`bool` · optional (explicit presence)

When true, CPU auto-scaling is enabled, allowing the database to use
up to three times the provisioned compute count during demand spikes.

### spec.isAutoScalingForStorageEnabled

`bool` · optional (explicit presence)

When true, storage auto-scaling is enabled, automatically expanding
storage when usage reaches the configured threshold.

### spec.adminPassword

`string` · sensitive

The administrator password for the database. Must be 12 to 30 characters,
contain at least one uppercase, one lowercase, and one numeric character.
Cannot contain the word "admin" or the double-quote character.
Mutually exclusive with secret_id.

### spec.secretId

`string | valueFrom`

OCID of a Vault secret containing the administrator password.
Use this instead of admin_password for production environments where
passwords should not appear in manifests.
Mutually exclusive with admin_password.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.secretVersionNumber

`int32`

Version number of the Vault secret. When omitted, the latest version
is used. Only applicable when secret_id is set.

### spec.subnetId

`string | valueFrom`

OCID of the subnet for private endpoint access. When set, the database
is only accessible from the specified subnet (and peered networks),
disabling public secure access endpoints.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.nsgIds

`[]string | valueFrom`

OCIDs of network security groups to associate with the private endpoint.
Maximum 5 NSGs. Only applicable when subnet_id is set.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.privateEndpointLabel

`string`

DNS label prefix for the private endpoint. Combined with the subnet's
DNS domain to form the FQDN for private access.

### spec.privateEndpointIp

`string`

Specific private IP address for the private endpoint within the subnet.
When omitted, OCI auto-assigns an available IP.

### spec.whitelistedIps

`[]string`

Client IP access control list. Each entry can be an IP address, a CIDR
block, or a VCN OCID (allowing all IPs in that VCN). When subnet_id is
set, this list further restricts access within the private network.

### spec.isMtlsConnectionRequired

`bool` · optional (explicit presence)

When true, only mutual TLS (mTLS) connections are allowed. When false,
both TLS and mTLS connections are accepted.

### spec.isAccessControlEnabled

`bool` · optional (explicit presence)

When true, enables database-level access control.
Applicable for Exadata Cloud@Customer deployments.

### spec.kmsKeyId

`string | valueFrom`

OCID of the KMS key for Transparent Data Encryption (TDE). When omitted,
Oracle-managed encryption is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.vaultId

`string | valueFrom`

OCID of the OCI Vault containing the KMS key. Required when kms_key_id
is set.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.isDedicated

`bool` · optional (explicit presence)

When true, provisions the database on dedicated Exadata infrastructure
(requires autonomous_container_database_id). When false or omitted,
uses shared (serverless) infrastructure. Cannot be changed after creation.

### spec.isFreeTier

`bool` · optional (explicit presence)

When true, provisions an Always Free autonomous database with limited
compute and storage. Cannot scale and is automatically reclaimed after
extended inactivity.

### spec.isDevTier

`bool` · optional (explicit presence)

When true, provisions a Developer tier database with reduced cost
suitable for development and testing workloads.

### spec.autonomousContainerDatabaseId

`string | valueFrom`

OCID of the autonomous container database for dedicated infrastructure
deployments. Required when is_dedicated is true.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.backupRetentionPeriodInDays

`int32`

Number of days to retain automatic backups.
When omitted (zero), the service default applies.

### spec.isLocalDataGuardEnabled

`bool` · optional (explicit presence)

When true, enables local Autonomous Data Guard for high availability.
A standby database is automatically provisioned in a different
availability domain within the same region.

### spec.autonomousMaintenanceScheduleType

`enum`

Maintenance schedule type controlling when patching occurs.

Allowed values (use exactly as shown):

- `maintenance_schedule_type_unspecified`
- `early` -- Early -- receive patches earlier than the regular schedule.
- `regular` -- Regular -- receive patches on the standard Oracle schedule.

### spec.customerContacts

`[]CustomerContact`

Customer contact email addresses for operational notifications
such as maintenance windows and critical alerts.

### spec.customerContacts[].email

`string` · required

Email address for operational notifications.

- rule: {"string":{"minLen":"1"}}

## Validation Rules

- `storage_size_mutual_exclusivity`: only one of data_storage_size_in_tbs or data_storage_size_in_gb may be set
- `admin_credential_mutual_exclusivity`: only one of admin_password or secret_id may be set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciAutonomousDatabase, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.autonomous_database_id` | `string` | OCID of the autonomous database. |
| `status.outputs.connection_string_high` | `string` | High-priority connection string for latency-sensitive workloads. Provides the highest level of service with the most resources. |
| `status.outputs.connection_string_medium` | `string` | Medium-priority connection string for typical application workloads. |
| `status.outputs.connection_string_low` | `string` | Low-priority connection string for batch and background workloads. Uses fewer resources and may queue requests when the database is busy. |
| `status.outputs.service_console_url` | `string` | URL of the OCI Service Console for this database. |
| `status.outputs.private_endpoint` | `string` | Private endpoint IP address. Empty when the database is not configured with a private endpoint (subnet_id not set). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciVaultSecret | `spec.rotationConfig.targetSystemDetails.adbId` | `status.outputs.autonomous_database_id` |

## See Also

- [Overview](../README.md)
