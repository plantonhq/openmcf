# AliCloudMongodbInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudMongodbInstanceSpec defines the configuration for an Alibaba Cloud
MongoDB replica-set instance.

MongoDB is a managed NoSQL document database used for content management,
IoT, real-time analytics, and mobile backends. Alibaba Cloud's ApsaraDB
for MongoDB supports replica-set deployments with configurable replication
factors, multi-zone high availability, and read-only replicas.

This component wraps a single `alicloud_mongodb_instance` resource
(replica-set mode). Sharding deployments use a separate Terraform resource
(`alicloud_mongodb_sharding_instance`) and are not covered here.

Provider resources:
  Terraform: alicloud_mongodb_instance
  Pulumi:    mongodb.Instance

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudMongodbInstance
metadata:
  name: alicloudmongodbinstance-demo
spec:
  region: cn-hangzhou
  engineVersion: "7.0"
  dbInstanceClass: dds.mongo.mid
  dbInstanceStorage: 20
  accountPassword: DemoP@ssw0rd123
  vswitchId:
    value: vsw-demo123
  replicationFactor: 3
  storageEngine: WiredTiger
  securityIpList:
    - "10.0.0.0/8"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.engineVersion` | `string` | yes |  |  |
| `spec.dbInstanceClass` | `string` | yes |  |  |
| `spec.dbInstanceStorage` | `int32` | yes |  |  |
| `spec.accountPassword` | `string` (sensitive) | yes |  |  |
| `spec.dbInstanceName` | `string` |  |  |  |
| `spec.zoneId` | `string` |  |  |  |
| `spec.secondaryZoneId` | `string` |  |  |  |
| `spec.hiddenZoneId` | `string` |  |  |  |
| `spec.replicationFactor` | `int32` |  | `3` |  |
| `spec.readonlyReplicas` | `int32` |  |  |  |
| `spec.storageEngine` | `string` |  | `WiredTiger` |  |
| `spec.storageType` | `string` |  |  |  |
| `spec.provisionedIops` | `int32` |  |  |  |
| `spec.instanceChargeType` | `string` |  | `PostPaid` |  |
| `spec.securityIpList` | `[]string` |  |  |  |
| `spec.securityGroupId` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.sslAction` | `string` |  |  |  |
| `spec.tdeStatus` | `string` |  |  |  |
| `spec.encryptionKey` | `string` |  |  |  |
| `spec.encrypted` | `bool` |  |  |  |
| `spec.cloudDiskEncryptionKey` | `string` |  |  |  |
| `spec.maintainStartTime` | `string` |  |  |  |
| `spec.maintainEndTime` | `string` |  |  |  |
| `spec.backupTime` | `string` |  |  |  |
| `spec.backupPeriod` | `[]string` |  |  |  |
| `spec.parameters` | `map<string, string>` |  |  |  |
| `spec.dbInstanceReleaseProtection` | `bool` |  |  |  |
| `spec.period` | `int32` |  |  |  |
| `spec.autoRenew` | `bool` |  |  |  |
| `spec.autoRenewDuration` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the MongoDB instance will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the MongoDB instance is placed. The instance inherits
the VPC and availability zone from the VSwitch.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.engineVersion

`string` · required

MongoDB engine version. This is a critical architectural choice --
different versions support different features and wire protocols.

- rule: engine_version must be one of: 4.0, 4.2, 4.4, 5.0, 6.0, 7.0
- rule: {"required":true}

### spec.dbInstanceClass

`string` · required

Instance specification class that determines CPU, memory, and IOPS.
Examples: "dds.mongo.mid", "dds.mongo.standard", "mongo.x8.medium",
"mongo.x8.large", "mongo.x8.xlarge".
See Alibaba Cloud documentation for the full list of instance classes.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.dbInstanceStorage

`int32` · required

Storage space in GB. The valid range depends on the instance class and
storage type. Typical values: 10, 20, 50, 100, 200, 500, 1000, 2000.

- rule: {"required":true,"int32":{"gte":1}}

### spec.accountPassword

`string` · required · sensitive

Root account password. 8-32 characters, must contain at least three of:
uppercase letters, lowercase letters, digits, special characters.

- rule: {"required":true,"string":{"minLen":"8","maxLen":"32"}}

### spec.dbInstanceName

`string`

Instance display name. 2-256 characters.
If omitted, defaults to the metadata.name.

- rule: db_instance_name must be between 2 and 256 characters when set

### spec.zoneId

`string`

Primary availability zone ID. If omitted, Alibaba Cloud selects a zone
from the VSwitch's region.
Examples: "cn-hangzhou-a", "cn-shanghai-b".

### spec.secondaryZoneId

`string`

Secondary availability zone for the standby node. When set together with
zone_id and hidden_zone_id, creates a three-zone high-availability
deployment where each replica-set member runs in a different AZ.

### spec.hiddenZoneId

`string`

Availability zone for the hidden node (the third replica-set member used
for elections and backups). Required for full three-zone HA deployments.

### spec.replicationFactor

`int32` · optional (explicit presence)

Number of replica-set nodes. Includes the primary, secondary, and hidden
nodes. Higher values increase read capacity and fault tolerance.
Default: 3

- default: `3`
- rule: replication_factor must be one of: 1, 3, 5, 7

### spec.readonlyReplicas

`int32` · optional (explicit presence)

Number of read-only replicas (0-5). Read replicas handle read traffic
to reduce load on the primary. Only available for replica-set instances
with replication_factor >= 3.

- rule: {"int32":{"lte":5,"gte":0}}

### spec.storageEngine

`string` · optional (explicit presence)

Storage engine. WiredTiger is the default and recommended engine for
most workloads. RocksDB is optimized for write-heavy workloads but is
deprecated in newer MongoDB versions.
Default: "WiredTiger"

- default: `WiredTiger`
- rule: storage_engine must be one of: WiredTiger, RocksDB

### spec.storageType

`string` · optional (explicit presence)

Storage type. Determines the underlying disk technology.
"cloud_essd1" (PL1), "cloud_essd2" (PL2), "cloud_essd3" (PL3) offer
increasing IOPS tiers. "cloud_auto" auto-scales IOPS. "local_ssd"
uses local NVMe disks for lowest latency.

- rule: storage_type must be one of: cloud_essd1, cloud_essd2, cloud_essd3, cloud_auto, local_ssd

### spec.provisionedIops

`int32` · optional (explicit presence)

Provisioned IOPS for cloud storage types. Only applicable when
storage_type is a cloud_essd variant or cloud_auto.

### spec.instanceChargeType

`string` · optional (explicit presence)

Billing method. "PostPaid" for pay-as-you-go, "PrePaid" for subscription.
Default: "PostPaid"

- default: `PostPaid`
- rule: instance_charge_type must be one of: PostPaid, PrePaid

### spec.securityIpList

`[]string`

IP addresses or CIDR blocks allowed to access the MongoDB instance.
Default: ["127.0.0.1"] (no access). Set to ["0.0.0.0/0"] to allow all
(not recommended for production).

### spec.securityGroupId

`string`

ECS security group ID to associate with the MongoDB instance.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the instance is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the MongoDB instance.

### spec.sslAction

`string` · optional (explicit presence)

SSL encryption action for client connections.
"Open" enables SSL. "Close" disables SSL. "Update" rotates the
SSL certificate. Once set, SSL cannot be fully removed (only toggled).

- rule: ssl_action must be one of: Open, Close, Update

### spec.tdeStatus

`string` · optional (explicit presence)

Transparent Data Encryption (TDE) status. Once enabled, TDE cannot be
disabled. Encrypts data at rest at the storage engine level.
Mutually exclusive with `encrypted` and `cloud_disk_encryption_key`.

- rule: tde_status must be enabled

### spec.encryptionKey

`string`

Custom KMS key ID for TDE encryption. When set with tde_status = "enabled",
uses a customer-managed key instead of the service-managed default.

### spec.encrypted

`bool` · optional (explicit presence)

Cloud disk encryption. Encrypts the underlying cloud disk at the
infrastructure layer. Mutually exclusive with `tde_status`.

### spec.cloudDiskEncryptionKey

`string`

KMS key ID for cloud disk encryption. Only applicable when encrypted
is true. Mutually exclusive with `tde_status` and `encryption_key`.

### spec.maintainStartTime

`string`

Maintenance window start time in UTC format (e.g., "02:00Z").
Alibaba Cloud may perform minor version upgrades or patches during
this window.

### spec.maintainEndTime

`string`

Maintenance window end time in UTC format (e.g., "06:00Z").

### spec.backupTime

`string`

Backup time window in UTC. Format: "HH:00Z-HH:00Z" (1-hour window).
Example: "02:00Z-03:00Z".

### spec.backupPeriod

`[]string`

Backup days of the week. Each entry is a day name.
Examples: ["Monday", "Wednesday", "Friday"].

### spec.parameters

`map<string, string>`

MongoDB engine parameters as key-value pairs. Allows tuning engine
settings such as operationProfiling.slowOpThresholdMs.
Refer to Alibaba Cloud documentation for valid parameter names.

### spec.dbInstanceReleaseProtection

`bool` · optional (explicit presence)

Enable release protection to prevent accidental deletion via the
console or API.

### spec.period

`int32` · optional (explicit presence)

Subscription period in months. Only applicable when instance_charge_type
is "PrePaid". Valid values: 1-9, 12, 24, 36.

- rule: period must be one of: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36

### spec.autoRenew

`bool` · optional (explicit presence)

Auto-renewal for PrePaid instances. Ignored for PostPaid instances.

### spec.autoRenewDuration

`int32` · optional (explicit presence)

Auto-renewal period in months (1-12). Only applicable when auto_renew
is true and instance_charge_type is "PrePaid".

- rule: {"int32":{"lte":12,"gte":1}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudMongodbInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The MongoDB instance ID assigned by Alibaba Cloud (e.g., "dds-xxxxx"). |
| `status.outputs.replica_set_name` | `string` | The replica set name. Applications use this in MongoDB connection strings (e.g., "?replicaSet=<name>") to enable automatic failover. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](../README.md)
