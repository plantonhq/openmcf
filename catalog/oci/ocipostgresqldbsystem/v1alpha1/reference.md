# OciPostgresqlDbSystem

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciPostgresqlDbSystemSpec defines the specification for an Oracle Cloud
Infrastructure PostgreSQL Database System -- a fully managed PostgreSQL
service running on dedicated compute shapes with configurable storage
durability, flexible OCPU/memory sizing, and built-in backup policies.

PostgreSQL DB Systems support read replicas via instance_count (2+ nodes),
regional or AD-local storage durability, and IOPS performance tiers.

This component manages the DB System resource itself. Backups and
configurations are separate OCI resources with independent lifecycles
and are not bundled here.

Excluded from v1:
  - source block (BACKUP restore) -- only fresh creation supported
  - patch_operations -- operational concern for replica management
  - apply_config -- operational concern for config change behavior
  - backup_policy.copy_policy -- cross-region backup copy
  - system_type (top-level) -- computed, rarely user-specified
  - defined_tags, system_tags -- managed by platform via freeform_tags

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.dbVersion` | `string` | yes |  |  |
| `spec.shape` | `string` | yes |  |  |
| `spec.instanceOcpuCount` | `int32` |  |  |  |
| `spec.instanceMemorySizeInGbs` | `int32` |  |  |  |
| `spec.instanceCount` | `int32` |  |  |  |
| `spec.networkDetails` | `NetworkDetails` | yes |  |  |
| `spec.networkDetails.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.networkDetails.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.networkDetails.isReaderEndpointEnabled` | `bool` |  |  |  |
| `spec.networkDetails.primaryDbEndpointPrivateIp` | `string` |  |  |  |
| `spec.storageDetails` | `StorageDetails` | yes |  |  |
| `spec.storageDetails.isRegionallyDurable` | `bool` |  |  |  |
| `spec.storageDetails.availabilityDomain` | `string` |  |  |  |
| `spec.storageDetails.iops` | `int64` |  |  |  |
| `spec.credentials` | `Credentials` |  |  |  |
| `spec.credentials.username` | `string` | yes |  |  |
| `spec.credentials.passwordDetails` | `PasswordDetails` | yes |  |  |
| `spec.credentials.passwordDetails.passwordType` | `enum` |  |  |  |
| `spec.credentials.passwordDetails.password` | `string` (sensitive) |  |  |  |
| `spec.credentials.passwordDetails.secretId` | `string \| valueFrom` |  |  |  |
| `spec.credentials.passwordDetails.secretVersion` | `string` |  |  |  |
| `spec.managementPolicy` | `ManagementPolicy` |  |  |  |
| `spec.managementPolicy.backupPolicy` | `BackupPolicy` |  |  |  |
| `spec.managementPolicy.backupPolicy.kind` | `enum` |  |  |  |
| `spec.managementPolicy.backupPolicy.backupStart` | `string` |  |  |  |
| `spec.managementPolicy.backupPolicy.retentionDays` | `int32` |  |  |  |
| `spec.managementPolicy.backupPolicy.daysOfTheMonth` | `[]int32` |  |  |  |
| `spec.managementPolicy.backupPolicy.daysOfTheWeek` | `[]string` |  |  |  |
| `spec.managementPolicy.maintenanceWindowStart` | `string` |  |  |  |
| `spec.configId` | `string \| valueFrom` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.instancesDetails` | `[]InstanceDetails` |  |  |  |
| `spec.instancesDetails[].displayName` | `string` |  |  |  |
| `spec.instancesDetails[].description` | `string` |  |  |  |
| `spec.instancesDetails[].privateIp` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the PostgreSQL DB System will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.dbVersion

`string` · required

PostgreSQL major version (e.g. "14", "15", "16").
Only the major version is stored; minor versions are managed by OCI.
Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.shape

`string` · required

Compute shape for the DB System instances. The provider automatically
prefixes "PostgreSQL." if not present. For flexible shapes, set
instance_ocpu_count and instance_memory_size_in_gbs separately.
Example: "VM.Standard.E4.Flex".

- rule: {"string":{"minLen":"1"}}

### spec.instanceOcpuCount

`int32`

Number of OCPUs allocated to each database instance. Used with
flexible shapes (e.g. VM.Standard.E4.Flex). Updatable.

### spec.instanceMemorySizeInGbs

`int32`

Amount of memory in gigabytes allocated to each database instance.
Used with flexible shapes. Updatable.

### spec.instanceCount

`int32`

Number of database instances (nodes). 1 creates a standalone system;
2 or more creates read replicas for read scaling and availability.

### spec.networkDetails

`NetworkDetails` · required

Network configuration for the DB System. Required.

- rule: {"required":true}

### spec.networkDetails.subnetId

`string | valueFrom` · required

OCID of the subnet where the DB System will be placed.
Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.networkDetails.nsgIds

`[]string | valueFrom`

OCIDs of network security groups for the DB System instances.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.networkDetails.isReaderEndpointEnabled

`bool` · optional (explicit presence)

When true, a reader endpoint is created for distributing read
queries across replica instances.

### spec.networkDetails.primaryDbEndpointPrivateIp

`string`

Specific private IP address for the primary (read-write) endpoint.
When omitted, OCI auto-assigns an available IP from the subnet.
Changing this forces recreation.

### spec.storageDetails

`StorageDetails` · required

Storage configuration for the DB System. Required.

- rule: {"required":true}

### spec.storageDetails.isRegionallyDurable

`bool`

When true, data is replicated across multiple availability domains
for regional durability. When false, data resides in a single AD
and availability_domain must be specified. Changing this forces
recreation.

### spec.storageDetails.availabilityDomain

`string`

Availability domain for AD-local storage. Required when
is_regionally_durable is false. Ignored when is_regionally_durable
is true. Example: "Uocm:PHX-AD-1". Changing this forces recreation.

### spec.storageDetails.iops

`int64`

Guaranteed input/output storage requests per second (IOPS).
Determines the performance tier. See OCI documentation for
supported values per shape. Updatable.

### spec.credentials

`Credentials`

Initial database credentials. When omitted, the provider applies
defaults. The entire credentials block is immutable after creation
(changing it forces recreation).

### spec.credentials.username

`string` · required

Administrator username. Changing this forces recreation.

- rule: {"string":{"minLen":"1"}}

### spec.credentials.passwordDetails

`PasswordDetails` · required

Password configuration using either a plain-text password or
an OCI Vault secret reference.

- rule: {"required":true}

### spec.credentials.passwordDetails.passwordType

`enum`

Password type discriminator.

Allowed values (use exactly as shown):

- `password_type_unspecified`
- `plain_text` -- Plain-text password provided directly.
- `vault_secret` -- Password stored in an OCI Vault secret.

### spec.credentials.passwordDetails.password

`string` · sensitive

Database password in plain text. Required when password_type
is plain_text. Not returned by the API after creation.

### spec.credentials.passwordDetails.secretId

`string | valueFrom`

OCID of the OCI Vault secret containing the password.
Required when password_type is vault_secret.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.credentials.passwordDetails.secretVersion

`string`

Version of the Vault secret to use. When omitted, the latest
version is used.

### spec.managementPolicy

`ManagementPolicy`

Management policy controlling backup schedule and maintenance window.

### spec.managementPolicy.backupPolicy

`BackupPolicy`

Backup policy for the DB System.

### spec.managementPolicy.backupPolicy.kind

`enum`

Backup schedule kind.

Allowed values (use exactly as shown):

- `backup_kind_unspecified`
- `daily` -- Daily backups at the specified backup_start time.
- `weekly` -- Weekly backups on specified days_of_the_week.
- `monthly` -- Monthly backups on specified days_of_the_month.
- `none` -- Backups disabled.

### spec.managementPolicy.backupPolicy.backupStart

`string`

Hour of the day (UTC) when the backup starts.
Required for daily, weekly, and monthly kinds.

### spec.managementPolicy.backupPolicy.retentionDays

`int32`

Number of days to retain backups after the DB System is deleted.

### spec.managementPolicy.backupPolicy.daysOfTheMonth

`[]int32`

Days of the month when backups run (1-28). Required for monthly.

- rule: {"repeated":{"maxItems":"28"}}

### spec.managementPolicy.backupPolicy.daysOfTheWeek

`[]string`

Days of the week when backups run (e.g. "MONDAY", "FRIDAY").
Required for weekly.

### spec.managementPolicy.maintenanceWindowStart

`string`

Start of the maintenance window in UTC. Format:
"{day-of-week} {time-of-day}" (e.g. "tue 02:00:00").

### spec.configId

`string | valueFrom`

OCID of a PostgreSQL configuration to apply. Configurations define
server parameter settings (shared_buffers, max_connections, etc.).
When omitted, the default configuration for the shape is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.description

`string`

User-provided description of the DB System.

### spec.instancesDetails

`[]InstanceDetails`

Per-instance configuration details. When specified, the list size
must match instance_count. Allows pinning display names, descriptions,
and private IP addresses per instance. Changing this forces recreation.

### spec.instancesDetails[].displayName

`string`

Display name for this database instance node.

### spec.instancesDetails[].description

`string`

Description of this database instance node.

### spec.instancesDetails[].privateIp

`string`

Specific private IP address for this instance within the subnet.
When omitted, OCI auto-assigns an available IP.

## Validation Rules

- `plain_text_password_required`: credentials.password_details.password is required when password_type is plain_text
- `vault_secret_id_required`: credentials.password_details.secret_id is required when password_type is vault_secret

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciPostgresqlDbSystem, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.db_system_id` | `string` | OCID of the PostgreSQL DB System. |
| `status.outputs.primary_db_endpoint_private_ip` | `string` | Private IP address of the primary (read-write) endpoint. |
| `status.outputs.admin_username` | `string` | Administrator username (computed after creation). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.networkDetails.subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.networkDetails.nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
