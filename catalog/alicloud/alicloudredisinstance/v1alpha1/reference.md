# AliCloudRedisInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudRedisInstanceSpec defines the configuration for an Alibaba Cloud
Redis (KVStore) instance.

Redis is a managed in-memory key-value store used for caching, session
management, real-time analytics, and message brokering. Alibaba Cloud's
KVStore service supports both Redis and Memcache engines through the same
resource type, with Redis being the default and dominant use case.

This component wraps a single `alicloud_kvstore_instance` resource.
Unlike RDS, Redis does not bundle accounts -- the instance-level `password`
field handles authentication for 80% of use cases. Redis 6.0+ ACL accounts
(`alicloud_kvstore_account`) are an advanced feature that can be managed
separately if needed.

Provider resources:
  Terraform: alicloud_kvstore_instance
  Pulumi:    kvstore.Instance

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudRedisInstance
metadata:
  name: alicloudredisinstance-demo
spec:
  region: cn-hangzhou
  instanceClass: redis.master.small.default
  password: DemoP@ssw0rd123
  vswitchId:
    value: vsw-demo123
  engineVersion: "7.0"
  dbInstanceName: demo-redis
  securityIps:
    - "10.0.0.0/8"
  config:
    maxmemory-policy: allkeys-lru
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vswitchId` | `string \| valueFrom` | yes |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.instanceClass` | `string` | yes |  |  |
| `spec.password` | `string` (sensitive) | yes |  |  |
| `spec.engineVersion` | `string` |  | `7.0` |  |
| `spec.instanceType` | `string` |  | `Redis` |  |
| `spec.dbInstanceName` | `string` |  |  |  |
| `spec.zoneId` | `string` |  |  |  |
| `spec.secondaryZoneId` | `string` |  |  |  |
| `spec.paymentType` | `string` |  | `PostPaid` |  |
| `spec.securityIps` | `[]string` |  |  |  |
| `spec.securityGroupId` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.shardCount` | `int32` |  |  |  |
| `spec.readOnlyCount` | `int32` |  |  |  |
| `spec.sslEnable` | `string` |  |  |  |
| `spec.tdeStatus` | `string` |  |  |  |
| `spec.encryptionKey` | `string` |  |  |  |
| `spec.vpcAuthMode` | `string` |  | `Open` |  |
| `spec.config` | `map<string, string>` |  |  |  |
| `spec.instanceReleaseProtection` | `bool` |  |  |  |
| `spec.maintainStartTime` | `string` |  |  |  |
| `spec.maintainEndTime` | `string` |  |  |  |
| `spec.backupPeriod` | `[]string` |  |  |  |
| `spec.backupTime` | `string` |  |  |  |
| `spec.privateConnectionPrefix` | `string` |  |  |  |
| `spec.autoRenew` | `bool` |  |  |  |
| `spec.autoRenewPeriod` | `int32` |  |  |  |
| `spec.period` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the Redis instance will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vswitchId

`string | valueFrom` · required

VSwitch ID where the Redis instance is placed. The instance inherits the
VPC and availability zone from the VSwitch.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.instanceClass

`string` · required

Instance specification class that determines memory capacity and
performance tier. Examples: "redis.master.small.default",
"redis.master.large.default", "redis.sharding.mid.default".
See Alibaba Cloud documentation for the full list of instance classes.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.password

`string` · required · sensitive

Instance authentication password. 8-32 characters, must contain at least
three of: uppercase letters, lowercase letters, digits, special characters.

- rule: {"required":true,"string":{"minLen":"8","maxLen":"32"}}

### spec.engineVersion

`string` · optional (explicit presence)

Redis engine version. Alibaba Cloud KVStore supports multiple versions.
Default: "7.0"

- default: `7.0`
- rule: engine_version must be one of: 2.8, 4.0, 5.0, 6.0, 7.0

### spec.instanceType

`string` · optional (explicit presence)

KVStore engine type. "Redis" is the default and recommended engine.
"Memcache" is supported but rare; the component is named RedisInstance
because Redis represents 95%+ of KVStore usage.
Default: "Redis"

- default: `Redis`
- rule: instance_type must be one of: Redis, Memcache

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

Secondary availability zone for multi-zone high availability deployments.
When set, the standby node runs in this zone for cross-AZ failover.
Must differ from zone_id.

### spec.paymentType

`string` · optional (explicit presence)

Billing method. "PostPaid" for pay-as-you-go, "PrePaid" for subscription.
Default: "PostPaid"

- default: `PostPaid`
- rule: payment_type must be one of: PostPaid, PrePaid

### spec.securityIps

`[]string`

IP addresses or CIDR blocks allowed to access the Redis instance.
Default: ["127.0.0.1"] (no access). Set to ["0.0.0.0/0"] to allow all
(not recommended for production).

### spec.securityGroupId

`string`

Security group ID to associate with the Redis instance. Multiple
security group IDs can be provided as a comma-separated string.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the instance is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the Redis instance.

### spec.shardCount

`int32` · optional (explicit presence)

Number of data shards for cluster (sharding) instances. Setting this to
a value > 1 creates a cluster-mode Redis instance. Cluster mode is
recommended for large-scale caching workloads that exceed the memory
capacity or throughput of a single node.
Valid range depends on the instance class; common values: 1, 2, 4, 8, 16, 32.

- rule: {"int32":{"gte":1}}

### spec.readOnlyCount

`int32` · optional (explicit presence)

Number of read-only replicas in the primary zone (1-9). Read replicas
enable read scaling by distributing read traffic across multiple nodes.
Only available for master-slave and cluster instances.

- rule: {"int32":{"lte":9,"gte":1}}

### spec.sslEnable

`string` · optional (explicit presence)

SSL/TLS encryption for client connections.
"Enable" turns on SSL. "Disable" turns off SSL.
"Update" rotates the SSL certificate.

- rule: ssl_enable must be one of: Enable, Disable, Update

### spec.tdeStatus

`string` · optional (explicit presence)

Transparent Data Encryption (TDE) status. Once enabled, TDE cannot be
disabled. Encrypts data at rest at the storage level.

- rule: tde_status must be Enabled

### spec.encryptionKey

`string`

Custom KMS key ID for TDE encryption. When set with tde_status = Enabled,
uses a customer-managed key instead of the service-managed default.

### spec.vpcAuthMode

`string` · optional (explicit presence)

VPC authentication mode. "Open" requires password authentication for
VPC-internal connections. "Close" disables password authentication for
VPC-internal connections (connections from within the VPC don't need a
password). Default: "Open"

- default: `Open`
- rule: vpc_auth_mode must be one of: Open, Close

### spec.config

`map<string, string>`

Redis configuration parameters as key-value pairs. Allows tuning Redis
engine settings such as maxmemory-policy, timeout, hz, etc.
Refer to Alibaba Cloud documentation for valid parameter names.

### spec.instanceReleaseProtection

`bool` · optional (explicit presence)

Enable release protection to prevent accidental deletion via the
console or API.

### spec.maintainStartTime

`string`

Maintenance window start time in UTC format (e.g., "02:00Z").
The Redis service may perform minor version upgrades or patches
during this window.

### spec.maintainEndTime

`string`

Maintenance window end time in UTC format (e.g., "06:00Z").

### spec.backupPeriod

`[]string`

Backup days of the week. Each entry is a day name.
Examples: ["Monday", "Wednesday", "Friday"].

### spec.backupTime

`string`

Backup time window in UTC. Format: "HH:mmZ-HH:mmZ".
Example: "02:00Z-03:00Z".

### spec.privateConnectionPrefix

`string`

Custom prefix for the private connection string.
When set, the private connection domain uses this prefix instead of
the auto-generated one.

### spec.autoRenew

`bool` · optional (explicit presence)

Auto-renewal for PrePaid instances. Ignored for PostPaid instances.

### spec.autoRenewPeriod

`int32` · optional (explicit presence)

Auto-renewal period in months (1-12). Only applicable when auto_renew
is true and payment_type is "PrePaid".

- rule: {"int32":{"lte":12,"gte":1}}

### spec.period

`string` · optional (explicit presence)

Subscription period in months. Only applicable when payment_type
is "PrePaid". Valid values: 1-9, 12, 24, 36.

- rule: period must be one of: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudRedisInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The Redis instance ID assigned by Alibaba Cloud (e.g., "r-xxxxx"). |
| `status.outputs.connection_domain` | `string` | The intranet (VPC-internal) connection domain for the Redis instance. Applications within the same VPC use this domain to connect. |
| `status.outputs.private_connection_port` | `string` | The private connection port for the Redis instance (default: "6379"). |
| `status.outputs.private_ip` | `string` | The private IP address assigned to the Redis instance within the VSwitch. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](../README.md)
