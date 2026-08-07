# GcpAlloydbInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `gcp.planton.dev/v1alpha1`

GcpAlloydbInstanceSpec defines an AlloyDB instance (`google_alloydb_instance`)
attached to an existing cluster.

Read pool instances scale read traffic independently of the bundled primary
in GcpAlloydbCluster. PRIMARY and SECONDARY types exist for advanced
topologies; most presets target READ_POOL.

## Example

```yaml
# Exercises a READ_POOL instance offline: cluster by literal path, two read
# nodes, and TLS-only client connections.
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpAlloydbInstance
metadata:
  name: hack-orders-read-pool
spec:
  cluster:
    value: projects/my-project/locations/us-central1/clusters/hack-orders
  instanceId: orders-read-pool
  instanceType: READ_POOL
  cpuCount: 2
  readPoolConfig:
    nodeCount: 2
  availabilityType: REGIONAL
  requireConnectors: true
  sslMode: ENCRYPTED_ONLY
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.projectId` | `string \| valueFrom` |  |  | GcpProject (`status.outputs.project_id`) |
| `spec.cluster` | `string \| valueFrom` | yes |  | GcpAlloydbCluster (`status.outputs.cluster_id`) |
| `spec.instanceId` | `string` | yes |  |  |
| `spec.instanceType` | `string` |  | `READ_POOL` |  |
| `spec.cpuCount` | `int32` |  |  |  |
| `spec.machineType` | `string` |  |  |  |
| `spec.readPoolConfig` | `GcpAlloydbInstanceReadPoolConfig` |  |  |  |
| `spec.readPoolConfig.nodeCount` | `int32` |  |  |  |
| `spec.availabilityType` | `string` |  |  |  |
| `spec.databaseFlags` | `map<string, string>` |  |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.queryInsightsConfig` | `GcpAlloydbInstanceQueryInsightsConfig` |  |  |  |
| `spec.queryInsightsConfig.queryPlansPerMinute` | `int32` |  |  |  |
| `spec.queryInsightsConfig.queryStringLength` | `int32` |  |  |  |
| `spec.queryInsightsConfig.recordApplicationTags` | `bool` |  |  |  |
| `spec.queryInsightsConfig.recordClientAddress` | `bool` |  |  |  |
| `spec.requireConnectors` | `bool` |  |  |  |
| `spec.sslMode` | `string` |  |  |  |
| `spec.activationPolicy` | `string` |  |  |  |
| `spec.enablePublicIp` | `bool` |  |  |  |
| `spec.enableOutboundPublicIp` | `bool` |  |  |  |
| `spec.authorizedExternalNetworks` | `[]GcpAlloydbInstanceAuthorizedExternalNetwork` |  |  |  |
| `spec.authorizedExternalNetworks[].cidrRange` | `string` | yes |  |  |
| `spec.pscInstanceConfig` | `GcpAlloydbInstancePscInstanceConfig` |  |  |  |
| `spec.pscInstanceConfig.allowedConsumerProjects` | `[]string` |  |  |  |
| `spec.pscInstanceConfig.pscAutoConnections` | `[]GcpAlloydbInstancePscAutoConnection` |  |  |  |
| `spec.pscInstanceConfig.pscAutoConnections[].consumerNetwork` | `string` |  |  |  |
| `spec.pscInstanceConfig.pscAutoConnections[].consumerProject` | `string` |  |  |  |
| `spec.pscInstanceConfig.pscInterfaceConfigs` | `[]GcpAlloydbInstancePscInterfaceConfig` |  |  |  |
| `spec.pscInstanceConfig.pscInterfaceConfigs[].networkAttachmentResource` | `string` |  |  |  |

## Field Details

### spec.projectId

`string | valueFrom`

The GCP project that owns the AlloyDB cluster.

- references: GcpProject (`status.outputs.project_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpProject, name: <that resource's name>, fieldPath: status.outputs.project_id}} -- a bare string does not parse

### spec.cluster

`string | valueFrom` · required

The AlloyDB cluster this instance belongs to. Accepts the full cluster
resource path or a reference to a GcpAlloydbCluster resource. Immutable.

- references: GcpAlloydbCluster (`status.outputs.cluster_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: GcpAlloydbCluster, name: <that resource's name>, fieldPath: status.outputs.cluster_id}} -- a bare string does not parse

### spec.instanceId

`string` · required

The instance ID within the cluster. Immutable.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-z][a-z0-9-]{0,61}[a-z0-9]$"}}

### spec.instanceType

`string` · optional (explicit presence)

Instance role: PRIMARY, READ_POOL, or SECONDARY. Immutable.
Presets default to READ_POOL for read scaling.

- default: `READ_POOL`
- rule: instance_type must be PRIMARY, READ_POOL, or SECONDARY

### spec.cpuCount

`int32`

Number of CPUs. GCP selects the machine family. Mutually exclusive with
machine_type.

### spec.machineType

`string`

Explicit machine type (e.g. "n2-highmem-4"). Mutually exclusive with
cpu_count.

### spec.readPoolConfig

`GcpAlloydbInstanceReadPoolConfig`

Read pool sizing. Required when instance_type is READ_POOL.

### spec.readPoolConfig.nodeCount

`int32`

Read capacity — number of nodes in the read pool instance.

- rule: {"int32":{"gte":1}}

### spec.availabilityType

`string`

ZONAL or REGIONAL placement. Read pools of size 1 can only be ZONAL;
pools with 2+ nodes can be REGIONAL.

- rule: availability_type must be ZONAL or REGIONAL

### spec.databaseFlags

`map<string, string>`

PostgreSQL database flags as key-value pairs.

### spec.displayName

`string`

Human-readable display name.

### spec.queryInsightsConfig

`GcpAlloydbInstanceQueryInsightsConfig`

Query insights configuration.

### spec.queryInsightsConfig.queryPlansPerMinute

`int32`

Number of query execution plans captured per minute. Range: 0-20.

- rule: {"int32":{"lte":20,"gte":0}}

### spec.queryInsightsConfig.queryStringLength

`int32`

Maximum length of the query string stored in insights. Range: 256-4500.
0 means unset — GCP applies its default (1024).

- rule: query_string_length must be between 256 and 4500 (or 0 for the GCP default)

### spec.queryInsightsConfig.recordApplicationTags

`bool`

Whether to record application tags for queries.

### spec.queryInsightsConfig.recordClientAddress

`bool`

Whether to record the client IP address for each query.

### spec.requireConnectors

`bool`

When true, only AlloyDB Auth Proxy / Language Connectors may connect.

### spec.sslMode

`string`

SSL mode: ENCRYPTED_ONLY or ALLOW_UNENCRYPTED_AND_ENCRYPTED.

- rule: ssl_mode must be ENCRYPTED_ONLY or ALLOW_UNENCRYPTED_AND_ENCRYPTED

### spec.activationPolicy

`string`

Instance activation: ALWAYS keeps the instance running (the default
posture); NEVER stops it. Flipping ALWAYS→NEVER→ALWAYS is the
stop/start lever — a stopped instance keeps its configuration and
storage but serves nothing and stops billing for compute. Mind the
ordering restrictions (stop read pools before the primary).
NOTE: managed connection pooling is deliberately not modeled — the
released google provider does not expose it for AlloyDB instances.

- rule: activation_policy must be ALWAYS or NEVER

### spec.enablePublicIp

`bool`

Enable a public IP on the instance.

### spec.enableOutboundPublicIp

`bool`

Enable outbound public IP for the instance.

### spec.authorizedExternalNetworks

`[]GcpAlloydbInstanceAuthorizedExternalNetwork`

CIDR ranges allowed to reach the public IP. Requires enable_public_ip.

### spec.authorizedExternalNetworks[].cidrRange

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.pscInstanceConfig

`GcpAlloydbInstancePscInstanceConfig`

Private Service Connect configuration.

### spec.pscInstanceConfig.allowedConsumerProjects

`[]string`

Consumer project numbers allowed to create PSC endpoints.

### spec.pscInstanceConfig.pscAutoConnections

`[]GcpAlloydbInstancePscAutoConnection`

PSC service automation connections.

### spec.pscInstanceConfig.pscAutoConnections[].consumerNetwork

`string`

Consumer network, e.g. "projects/vpc-host/global/networks/default".

### spec.pscInstanceConfig.pscAutoConnections[].consumerProject

`string`

Consumer project ID (not project number).

### spec.pscInstanceConfig.pscInterfaceConfigs

`[]GcpAlloydbInstancePscInterfaceConfig`

PSC interfaces for outbound connectivity (0 or 1 supported by AlloyDB).

### spec.pscInstanceConfig.pscInterfaceConfigs[].networkAttachmentResource

`string`

Network attachment resource in the consumer project.

## Validation Rules

- `machine_config_mutual_exclusion`: only one of cpu_count or machine_type may be set
- `read_pool_requires_node_count`: READ_POOL instances require read_pool_config.node_count >= 1
- `read_pool_config_only_for_read_pool`: read_pool_config applies to READ_POOL instances only
- `authorized_networks_require_public_ip`: authorized_external_networks requires enable_public_ip

## Outputs

Reference an output from another manifest as `valueFrom: {kind: GcpAlloydbInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_name` | `string` | Fully qualified instance resource name. Format: projects/{project}/locations/{location}/clusters/{cluster}/instances/{instance} |
| `status.outputs.ip_address` | `string` | Private IP address of the instance (primary connection endpoint). |
| `status.outputs.state` | `string` | Current state of the instance (e.g. READY, CREATING). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.projectId` | GcpProject | `status.outputs.project_id` |
| `spec.cluster` | GcpAlloydbCluster | `status.outputs.cluster_id` |

## See Also

- [Overview](../README.md)
