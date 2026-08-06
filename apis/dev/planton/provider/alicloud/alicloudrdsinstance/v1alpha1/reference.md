# AliCloudRdsInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudRdsInstanceSpec defines the configuration for an Alibaba Cloud RDS
(Relational Database Service) instance with bundled databases, accounts, and
account privileges.

RDS is a managed relational database service that supports multiple engines
(MySQL, PostgreSQL, SQL Server, MariaDB, PPAS) through a single resource
type. The engine is selected via the `engine` field, following the
provider-authentic pattern where `alicloud_db_instance` handles all engines
with the same schema (per DD02).

This component bundles the instance with its databases, accounts, and account
privileges (per DD07 composite bundling) because an RDS instance without at
least one database and account is incomplete for application use.

The bundled flow:
  1. Create the RDS instance in the specified VSwitch.
  2. Create databases within the instance.
  3. Create accounts with passwords.
  4. Grant account privileges on specific databases.

The intranet connection string and port are exposed as stack outputs. Public
internet endpoints (alicloud_db_connection) are intentionally excluded --
exposing a database to the internet is a security-sensitive decision that
should be managed separately.

Provider resources:
  Terraform: alicloud_db_instance + alicloud_db_database + alicloud_rds_account + alicloud_db_account_privilege
  Pulumi:    rds.Instance + rds.Database + rds.RdsAccount + rds.AccountPrivilege

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudRdsInstance
metadata:
  name: alicloudrdsinstance-demo
spec:
  region: cn-hangzhou
  engine: MySQL
  engineVersion: "8.0"
  instanceType: rds.mysql.s2.large
  instanceStorage: 50
  vswitchId:
    value: vsw-demo123
  instanceName: demo-mysql
  category: HighAvailability
  dbInstanceStorageType: cloud_essd
  securityIps:
    - "10.0.0.0/8"
  databases:
    - name: appdb
      characterSet: utf8mb4
  accounts:
    - accountName: app_user
      accountPassword: DemoP@ssw0rd123
      privileges:
        - databaseNames: [appdb]
          privilege: ReadWrite
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.engine` | `string` | yes |  |  |
| `spec.engineVersion` | `string` | yes |  |  |
| `spec.instanceType` | `string` | yes |  |  |
| `spec.instanceStorage` | `int32` | yes |  |  |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.instanceName` | `string` |  |  |  |
| `spec.instanceChargeType` | `string` |  | `Postpaid` |  |
| `spec.category` | `string` |  | `HighAvailability` |  |
| `spec.dbInstanceStorageType` | `string` |  |  |  |
| `spec.zoneId` | `string` |  |  |  |
| `spec.zoneIdSlaveA` | `string` |  |  |  |
| `spec.securityIps` | `[]string` |  |  |  |
| `spec.securityGroupIds` | `[]string` |  |  |  |
| `spec.monitoringPeriod` | `int32` |  |  |  |
| `spec.maintainTime` | `string` |  |  |  |
| `spec.deletionProtection` | `bool` |  |  |  |
| `spec.sslAction` | `string` |  |  |  |
| `spec.tdeStatus` | `string` |  |  |  |
| `spec.encryptionKey` | `string` |  |  |  |
| `spec.autoRenew` | `bool` |  |  |  |
| `spec.autoRenewPeriod` | `int32` |  |  |  |
| `spec.period` | `int32` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.parameters` | `[]AliCloudRdsParameter` |  |  |  |
| `spec.parameters[].name` | `string` | yes |  |  |
| `spec.parameters[].value` | `string` | yes |  |  |
| `spec.databases` | `[]AliCloudRdsDatabase` |  |  |  |
| `spec.databases[].name` | `string` | yes |  |  |
| `spec.databases[].characterSet` | `string` |  |  |  |
| `spec.databases[].description` | `string` |  |  |  |
| `spec.accounts` | `[]AliCloudRdsAccount` |  |  |  |
| `spec.accounts[].accountName` | `string` | yes |  |  |
| `spec.accounts[].accountPassword` | `string` (sensitive) | yes |  |  |
| `spec.accounts[].accountType` | `string` |  | `Normal` |  |
| `spec.accounts[].accountDescription` | `string` |  |  |  |
| `spec.accounts[].privileges` | `[]AliCloudRdsAccountPrivilege` |  |  |  |
| `spec.accounts[].privileges[].databaseNames` | `[]string` | yes |  |  |
| `spec.accounts[].privileges[].privilege` | `string` |  | `ReadOnly` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the RDS instance will be created.
Must match the region of the VSwitch.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.engine

`string` · required

Database engine. Determines which relational database software runs on the
instance. Each engine has its own set of valid versions, instance types,
and character sets.

- rule: engine must be one of: MySQL, PostgreSQL, SQLServer, MariaDB, PPAS
- rule: {"required":true}

### spec.engineVersion

`string` · required

Database engine version. Must be a version supported by the selected
engine. Examples: "8.0", "5.7" (MySQL), "16.0", "15.0" (PostgreSQL),
"2019_ent" (SQL Server), "10.3" (MariaDB).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceType

`string` · required

RDS instance class that determines CPU, memory, and I/O capacity.
Examples: "rds.mysql.s2.large", "rds.pg.s2.large",
"rds.mssql.s2.large", "rds.mysql.t1.small".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.instanceStorage

`int32` · required

Storage size in GB. The valid range depends on the engine, instance type,
and storage type. Typical ranges: 20-6000 GB for cloud_essd.

- rule: {"required":true,"int32":{"gt":0}}

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the RDS instance is placed. The instance inherits the
VPC and availability zone from the VSwitch.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.instanceName

`string`

RDS instance name. 2-256 characters.
If omitted, defaults to the metadata.name.

- rule: instance_name must be between 2 and 256 characters when set

### spec.instanceChargeType

`string` · optional (explicit presence)

Billing method. "Postpaid" for pay-as-you-go, "Prepaid" for subscription.
Default: "Postpaid"

- default: `Postpaid`
- rule: instance_charge_type must be one of: Postpaid, Prepaid

### spec.category

`string` · optional (explicit presence)

Instance category that determines the deployment architecture.
"Basic" -- single node, suitable for dev/test.
"HighAvailability" -- primary + standby with automatic failover.
"AlwaysOn" -- SQL Server AlwaysOn Availability Group clusters.
"Finance" -- three-node enterprise edition for financial scenarios.
"cluster" -- MySQL NDB cluster mode.
Default: "HighAvailability"

- default: `HighAvailability`
- rule: category must be one of: Basic, HighAvailability, AlwaysOn, Finance, cluster

### spec.dbInstanceStorageType

`string` · optional (explicit presence)

Storage type for the instance's data disks.
"local_ssd" -- local SSD (only for specific instance types).
"cloud_ssd" -- standard cloud SSD.
"cloud_essd" -- enhanced SSD (recommended for most workloads).
"cloud_essd2" -- ESSD PL2 (higher IOPS).
"cloud_essd3" -- ESSD PL3 (highest IOPS).

- rule: db_instance_storage_type must be one of: local_ssd, cloud_ssd, cloud_essd, cloud_essd2, cloud_essd3

### spec.zoneId

`string`

Availability zone ID for the primary instance.
If omitted, Alibaba Cloud selects a zone from the VSwitch's region.
Examples: "cn-hangzhou-a", "cn-shanghai-b".

### spec.zoneIdSlaveA

`string`

Availability zone ID for the standby instance in HA deployments.
Must differ from zone_id for cross-AZ high availability.
Only applicable when category is "HighAvailability" or "Finance".

### spec.securityIps

`[]string`

IP addresses or CIDR blocks allowed to access the RDS instance.
Default: ["127.0.0.1"] (no access). Set to ["0.0.0.0/0"] to allow all
(not recommended for production).

### spec.securityGroupIds

`[]string`

VPC security group IDs to associate with the RDS instance.
Provides network-level access control in addition to the IP whitelist.

### spec.monitoringPeriod

`int32` · optional (explicit presence)

Monitoring data collection interval in seconds.
Valid values: 5, 10, 60, 300. Lower values provide finer granularity
but may increase monitoring costs.

- rule: monitoring_period must be one of: 5, 10, 60, 300

### spec.maintainTime

`string`

Maintenance window in UTC. The RDS service may perform minor version
upgrades or patches during this window.
Format: "HH:mmZ-HH:mmZ" (e.g., "02:00Z-06:00Z").

### spec.deletionProtection

`bool` · optional (explicit presence)

Enable deletion protection to prevent accidental deletion via the
console or API.

### spec.sslAction

`string` · optional (explicit presence)

SSL/TLS configuration for encrypted client connections.
"Open" enables SSL and generates server certificates.
"Close" disables SSL (default).

- rule: ssl_action must be one of: Open, Close

### spec.tdeStatus

`string` · optional (explicit presence)

Transparent Data Encryption (TDE) status. TDE encrypts data at rest
at the storage level. Once enabled, TDE cannot be disabled.
"Enabled" turns on TDE. "Disabled" is the default.

- rule: tde_status must be one of: Enabled, Disabled

### spec.encryptionKey

`string`

KMS key ID for disk encryption. When set, the instance's storage is
encrypted using this customer-managed key.

### spec.autoRenew

`bool` · optional (explicit presence)

Auto-renewal for Prepaid instances. Ignored for Postpaid instances.

### spec.autoRenewPeriod

`int32` · optional (explicit presence)

Auto-renewal period in months (1-12). Only applicable when auto_renew
is true and instance_charge_type is "Prepaid".

- rule: {"int32":{"lte":12,"gte":1}}

### spec.period

`int32` · optional (explicit presence)

Subscription period in months. Only applicable when instance_charge_type
is "Prepaid". Valid values: 1-9, 12, 24, 36.

- rule: period must be one of: 1-9, 12, 24, 36

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the instance is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the RDS instance resource.

### spec.parameters

`[]AliCloudRdsParameter`

Database parameter overrides. Each parameter sets a specific database
engine configuration value, overriding the default from the parameter
group. Refer to Alibaba Cloud documentation for valid parameter names
and value ranges per engine.

### spec.parameters[].name

`string` · required

Parameter name (e.g., "innodb_buffer_pool_size", "max_connections",
"shared_buffers").

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.parameters[].value

`string` · required

Parameter value as a string. The database engine interprets the value
according to the parameter's type.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.databases

`[]AliCloudRdsDatabase`

Databases to create within the RDS instance.

### spec.databases[].name

`string` · required

Database name. Must match the pattern [a-z][a-z0-9_-]*[a-z0-9] and
be unique within the instance.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.databases[].characterSet

`string`

Character set for the database. Defaults vary by engine:
MySQL: "utf8mb4", PostgreSQL: "UTF8", SQL Server: "Chinese_PRC_CI_AS",
MariaDB: "utf8mb4". Set explicitly to override the engine default.

### spec.databases[].description

`string`

Human-readable description of the database.

### spec.accounts

`[]AliCloudRdsAccount`

Database accounts to create within the RDS instance.

### spec.accounts[].accountName

`string` · required

Account login name. Must match [a-z][a-z0-9_]{0,61}[a-z0-9].

- rule: {"required":true,"string":{"minLen":"2"}}

### spec.accounts[].accountPassword

`string` · required · sensitive

Account password. Must meet the Alibaba Cloud password complexity
requirements: 8-32 characters, containing at least three of uppercase,
lowercase, digits, and special characters.

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

`[]AliCloudRdsAccountPrivilege`

Database privileges to grant to this account. Each entry grants a
privilege level on a set of databases. Only applicable when account_type
is "Normal"; Super accounts have implicit full access.

### spec.accounts[].privileges[].databaseNames

`[]string` · required

Database names to grant access to. Each name must match a database
defined in the spec's databases list.

- rule: {"required":true,"repeated":{"minItems":"1"}}

### spec.accounts[].privileges[].privilege

`string` · optional (explicit presence)

Privilege level to grant on the specified databases.
"ReadOnly" -- SELECT access.
"ReadWrite" -- SELECT, INSERT, UPDATE, DELETE access.
"DDLOnly" -- DDL operations only (CREATE, ALTER, DROP).
"DMLOnly" -- DML operations only (INSERT, UPDATE, DELETE).
"DBOwner" -- full database ownership (SQL Server only).
Default: "ReadOnly"

- default: `ReadOnly`
- rule: privilege must be one of: ReadOnly, ReadWrite, DDLOnly, DMLOnly, DBOwner

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudRdsInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The RDS instance ID assigned by Alibaba Cloud (e.g., "rm-xxxxx"). |
| `status.outputs.connection_string` | `string` | The intranet (VPC-internal) connection endpoint for the RDS instance. Applications within the same VPC use this endpoint to connect. |
| `status.outputs.port` | `string` | The database service port (e.g., "3306" for MySQL, "5432" for PostgreSQL, "1433" for SQL Server). |
| `status.outputs.database_ids` | `map<string, string>` | Map of database names to their IDs within the RDS instance. Keys are the names specified in spec.databases[].name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
