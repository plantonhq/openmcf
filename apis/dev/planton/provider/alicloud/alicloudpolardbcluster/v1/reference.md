# AliCloudPolardbCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudPolardbClusterSpec defines the configuration for an Alibaba Cloud
PolarDB cluster with bundled databases, accounts, and account privileges.

PolarDB is Alibaba Cloud's cloud-native relational database built on a
shared-storage, compute-storage-separated architecture. It supports MySQL,
PostgreSQL, and Oracle compatibility modes. Unlike RDS, PolarDB uses a
cluster model where one primary node and one or more read-only replicas
share a single distributed storage layer.

This component is separate from AliCloudRdsInstance (per DD06) because
PolarDB has a fundamentally different resource schema: cluster-based
topology, node-class sizing, separate endpoint management, and different
Terraform/Pulumi resource types.

The bundled flow (per DD07 composite bundling):
  1. Create the PolarDB cluster in the specified VSwitch.
  2. Create databases within the cluster.
  3. Create accounts with passwords.
  4. Grant account privileges on specific databases.

PolarDB automatically provisions a primary endpoint and a cluster endpoint.
The primary endpoint's connection string and port are exposed as stack
outputs. Custom endpoints (read-only distribution, connection pooling) are
intentionally excluded -- they can be managed as a separate resource if
needed.

Provider resources:
  Terraform: alicloud_polardb_cluster + alicloud_polardb_database + alicloud_polardb_account + alicloud_polardb_account_privilege
  Pulumi:    polardb.Cluster + polardb.Database + polardb.Account + polardb.AccountPrivilege

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudPolardbCluster
metadata:
  name: alicloudpolardbcluster-demo
spec:
  region: cn-hangzhou
  dbType: MySQL
  dbVersion: "8.0"
  dbNodeClass: polar.mysql.x4.large
  vswitchId:
    value: vsw-demo123
  dbNodeCount: 2
  description: Demo PolarDB cluster
  securityIps:
    - "10.0.0.0/8"
  databases:
    - dbName: appdb
      characterSetName: utf8mb4
  accounts:
    - accountName: app_user
      accountPassword: DemoP@ssw0rd123
      privileges:
        - dbNames: [appdb]
          accountPrivilege: ReadWrite
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.dbType` | `string` | yes |  |  |
| `spec.dbVersion` | `string` | yes |  |  |
| `spec.dbNodeClass` | `string` | yes |  |  |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.dbNodeCount` | `int32` |  | `2` |  |
| `spec.description` | `string` |  |  |  |
| `spec.payType` | `string` |  | `PostPaid` |  |
| `spec.period` | `int32` |  |  |  |
| `spec.renewalStatus` | `string` |  |  |  |
| `spec.autoRenewPeriod` | `int32` |  |  |  |
| `spec.zoneId` | `string` |  |  |  |
| `spec.securityIps` | `[]string` |  |  |  |
| `spec.securityGroupIds` | `[]string` |  |  |  |
| `spec.maintainTime` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.creationCategory` | `string` |  |  |  |
| `spec.subCategory` | `string` |  |  |  |
| `spec.storageType` | `string` |  |  |  |
| `spec.storageSpace` | `int32` |  |  |  |
| `spec.tdeStatus` | `string` |  |  |  |
| `spec.encryptionKey` | `string` |  |  |  |
| `spec.deletionLock` | `int32` |  |  |  |
| `spec.collectorStatus` | `string` |  |  |  |
| `spec.backupRetentionPolicyOnClusterDeletion` | `string` |  |  |  |
| `spec.parameters` | `[]AliCloudPolardbParameter` |  |  |  |
| `spec.parameters[].name` | `string` | yes |  |  |
| `spec.parameters[].value` | `string` | yes |  |  |
| `spec.databases` | `[]AliCloudPolardbDatabase` |  |  |  |
| `spec.databases[].dbName` | `string` | yes |  |  |
| `spec.databases[].characterSetName` | `string` |  |  |  |
| `spec.databases[].dbDescription` | `string` |  |  |  |
| `spec.databases[].collate` | `string` |  |  |  |
| `spec.databases[].ctype` | `string` |  |  |  |
| `spec.accounts` | `[]AliCloudPolardbAccount` |  |  |  |
| `spec.accounts[].accountName` | `string` | yes |  |  |
| `spec.accounts[].accountPassword` | `string` (sensitive) | yes |  |  |
| `spec.accounts[].accountType` | `string` |  | `Normal` |  |
| `spec.accounts[].accountDescription` | `string` |  |  |  |
| `spec.accounts[].privileges` | `[]AliCloudPolardbAccountPrivilege` |  |  |  |
| `spec.accounts[].privileges[].dbNames` | `[]string` | yes |  |  |
| `spec.accounts[].privileges[].accountPrivilege` | `string` |  | `ReadOnly` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the PolarDB cluster will be created.
Must match the region of the VSwitch.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dbType

`string` · required

Database engine type. PolarDB supports three compatibility modes, each
with its own set of valid versions and node classes.

- rule: db_type must be one of: MySQL, PostgreSQL, Oracle
- rule: {"required":true}

### spec.dbVersion

`string` · required

Database engine version. Must be a version supported by the selected
db_type. Examples: "8.0", "5.7", "5.6" (MySQL), "14", "11" (PostgreSQL),
"11" (Oracle).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dbNodeClass

`string` · required

PolarDB node instance class that determines CPU and memory per node.
Examples: "polar.mysql.x4.large", "polar.mysql.x4.xlarge",
"polar.pg.x4.large", "polar.o.x4.large".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the PolarDB cluster is placed. The cluster inherits
the VPC and availability zone from the VSwitch.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.dbNodeCount

`int32` · optional (explicit presence)

Number of nodes in the cluster. Includes 1 primary node and (N-1) read
replicas. Minimum 1, maximum 16. Default: 2 (1 primary + 1 read replica).

- default: `2`
- rule: db_node_count must be between 1 and 16

### spec.description

`string`

Cluster description. 2-256 characters.

- rule: description must be between 2 and 256 characters when set

### spec.payType

`string` · optional (explicit presence)

Billing method. "PostPaid" for pay-as-you-go, "PrePaid" for subscription.
Default: "PostPaid"

- default: `PostPaid`
- rule: pay_type must be one of: PostPaid, PrePaid

### spec.period

`int32` · optional (explicit presence)

Subscription period in months. Only applicable when pay_type is "PrePaid".
Valid values: 1-9, 12, 24, 36.

- rule: period must be one of: 1-9, 12, 24, 36

### spec.renewalStatus

`string` · optional (explicit presence)

Auto-renewal behavior for PrePaid clusters. Ignored for PostPaid clusters.
"AutoRenewal" -- automatically renew before expiration.
"Normal" -- send expiration notification but do not auto-renew.
"NotRenewal" -- no notification, no auto-renew (default).

- rule: renewal_status must be one of: AutoRenewal, Normal, NotRenewal

### spec.autoRenewPeriod

`int32` · optional (explicit presence)

Auto-renewal period in months. Only applicable when renewal_status is
"AutoRenewal". Valid values: 1, 2, 3, 6, 12, 24, 36. Default: 1.

- rule: auto_renew_period must be one of: 1, 2, 3, 6, 12, 24, 36

### spec.zoneId

`string`

Availability zone ID for the primary node.
If omitted, Alibaba Cloud selects a zone from the VSwitch's region.
Examples: "cn-hangzhou-a", "cn-shanghai-b".

### spec.securityIps

`[]string`

IP addresses or CIDR blocks allowed to access the PolarDB cluster.
Default: ["127.0.0.1"] (no access). Set to ["0.0.0.0/0"] to allow all
(not recommended for production).

### spec.securityGroupIds

`[]string`

VPC security group IDs to associate with the PolarDB cluster.
You can add a maximum of three security groups.

### spec.maintainTime

`string`

Maintenance window in UTC. PolarDB may perform minor version upgrades
or patches during this window.
Format: "HH:00Z-HH:00Z" (must be a 1-hour interval on the hour,
e.g., "02:00Z-03:00Z").

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the cluster is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the PolarDB cluster resource.

### spec.creationCategory

`string` · optional (explicit presence)

PolarDB edition that determines the storage architecture.
"Normal" -- Enterprise Edition with shared distributed storage.
"Basic" -- single-node edition for dev/test.
"ArchiveNormal" -- archival storage edition.
"NormalMultimaster" -- multi-master edition (MySQL only).
"SENormal" -- Standard Edition with local ESSD storage.

- rule: creation_category must be one of: Normal, Basic, ArchiveNormal, NormalMultimaster, SENormal

### spec.subCategory

`string` · optional (explicit presence)

Sub-category for the cluster. Determines the resource allocation model.
"Exclusive" -- dedicated resources (higher performance).
"General" -- shared resources (lower cost).
Only applicable to MySQL clusters.

- rule: sub_category must be one of: Exclusive, General

### spec.storageType

`string` · optional (explicit presence)

Storage type. Determines the underlying storage medium.
Enterprise Edition: "PSL5" (performance level 5), "PSL4" (performance level 4).
Standard Edition: "ESSDPL0", "ESSDPL1", "ESSDPL2", "ESSDPL3" (ESSD performance levels),
"ESSDAUTOPL" (auto performance level).

- rule: storage_type must be one of: PSL5, PSL4, ESSDPL0, ESSDPL1, ESSDPL2, ESSDPL3, ESSDAUTOPL

### spec.storageSpace

`int32` · optional (explicit presence)

Storage space in GB. Only applicable for Standard Edition clusters
(creation_category = "SENormal") where storage is pre-allocated.
Enterprise Edition clusters use auto-scaling shared storage and do not
need this field. Valid range: 20-100000 GB.

- rule: storage_space must be between 20 and 100000 when set

### spec.tdeStatus

`string` · optional (explicit presence)

Transparent Data Encryption (TDE) status. TDE encrypts data at rest
at the storage level. Once enabled, TDE cannot be disabled.
"Enabled" turns on TDE. "Disabled" is the default.

- rule: tde_status must be one of: Enabled, Disabled

### spec.encryptionKey

`string`

KMS encryption key ID for TDE. Required when tde_status is "Enabled".
Cannot be changed after TDE is enabled.

### spec.deletionLock

`int32` · optional (explicit presence)

Deletion protection lock. Prevents accidental cluster deletion via
console or API. Set to 1 to enable, 0 to disable.

- rule: deletion_lock must be 0 or 1

### spec.collectorStatus

`string` · optional (explicit presence)

SQL audit log collector status.
"Enable" turns on audit logging for compliance and troubleshooting.
"Disabled" turns it off.

- rule: collector_status must be one of: Enable, Disabled

### spec.backupRetentionPolicyOnClusterDeletion

`string` · optional (explicit presence)

Backup retention policy when the cluster is deleted.
"ALL" -- retain all backups.
"LATEST" -- retain only the most recent backup.
"NONE" -- delete all backups immediately.

- rule: backup_retention_policy_on_cluster_deletion must be one of: ALL, LATEST, NONE

### spec.parameters

`[]AliCloudPolardbParameter`

Cluster parameter overrides. Each parameter sets a specific database
engine configuration value. Refer to Alibaba Cloud documentation for
valid parameter names and value ranges per engine.

### spec.parameters[].name

`string` · required

Parameter name (e.g., "loose_innodb_buffer_pool_size", "max_connections").

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.parameters[].value

`string` · required

Parameter value as a string. The database engine interprets the value
according to the parameter's type.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.databases

`[]AliCloudPolardbDatabase`

Databases to create within the PolarDB cluster.

### spec.databases[].dbName

`string` · required

Database name. Must start with a letter, consist of lowercase letters,
numbers, and underscores, and be no more than 64 characters.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.databases[].characterSetName

`string`

Character set for the database. Default: "utf8" for MySQL, varies for
PostgreSQL and Oracle. Common values: "utf8", "utf8mb4", "gbk" (MySQL),
"UTF8", "SQL_ASCII" (PostgreSQL/Oracle).

### spec.databases[].dbDescription

`string`

Human-readable description of the database. 2-256 characters.

### spec.databases[].collate

`string`

Collation rules for the database. Only applicable to PolarDB for
PostgreSQL and Oracle compatibility modes. Must be compatible with
the character_set_name. Default: "C".

### spec.databases[].ctype

`string`

Character type setting. Only applicable to PolarDB for PostgreSQL and
Oracle compatibility modes. Must be compatible with the
character_set_name. Default: "C".

### spec.accounts

`[]AliCloudPolardbAccount`

Database accounts to create within the PolarDB cluster.

### spec.accounts[].accountName

`string` · required

Account login name. Must start with a lowercase letter, end with a
letter or number, consist of lowercase letters, numbers, or underscores,
and be 2-16 characters long.

- rule: {"required":true,"string":{"minLen":"2"}}

### spec.accounts[].accountPassword

`string` · required · sensitive

Account password. Must be 8-32 characters and contain at least three of:
uppercase letters, lowercase letters, digits, and special characters
(!@#$%^&*()_+-=).

- rule: {"required":true,"string":{"minLen":"8"}}

### spec.accounts[].accountType

`string` · optional (explicit presence)

Account type.
"Normal" -- standard account with explicit database privileges.
"Super" -- superuser with full privileges on all databases.
Default: "Normal"

- default: `Normal`
- rule: account_type must be one of: Normal, Super

### spec.accounts[].accountDescription

`string`

Human-readable description of the account's purpose.

### spec.accounts[].privileges

`[]AliCloudPolardbAccountPrivilege`

Database privileges to grant to this account. Each entry grants a
privilege level on a set of databases. Only applicable when account_type
is "Normal"; Super accounts have implicit full access.

### spec.accounts[].privileges[].dbNames

`[]string` · required

Database names to grant access to. Each name must match a database
defined in the spec's databases list.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.accounts[].privileges[].accountPrivilege

`string` · optional (explicit presence)

Privilege level to grant on the specified databases.
"ReadOnly" -- SELECT access.
"ReadWrite" -- SELECT, INSERT, UPDATE, DELETE access.
"DDLOnly" -- DDL operations only (CREATE, ALTER, DROP).
"DMLOnly" -- DML operations only (INSERT, UPDATE, DELETE).
Default: "ReadOnly"

- default: `ReadOnly`
- rule: account_privilege must be one of: ReadOnly, ReadWrite, DDLOnly, DMLOnly

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudPolardbCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The PolarDB cluster ID assigned by Alibaba Cloud (e.g., "pc-xxxxx"). |
| `status.outputs.connection_string` | `string` | The primary endpoint connection string for the PolarDB cluster. Applications connect to this endpoint for read-write operations. |
| `status.outputs.port` | `string` | The database service port (e.g., "3306" for MySQL, "5432" for PostgreSQL/Oracle). |
| `status.outputs.database_ids` | `map<string, string>` | Map of database names to their IDs within the PolarDB cluster. Keys are the names specified in spec.databases[].db_name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
