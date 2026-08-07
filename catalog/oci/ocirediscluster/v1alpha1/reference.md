# OciRedisCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciRedisClusterSpec defines the specification for an OCI Cache (Redis)
cluster -- a fully managed, Redis-compatible in-memory caching service
that supports both sharded and non-sharded deployment modes.

Non-sharded clusters provide a single primary with optional replicas
for high availability. Sharded clusters distribute data across multiple
shards for horizontal scaling, each with its own primary and replicas.

This component manages the cluster resource itself. Config sets are
separate OCI resources with independent lifecycles and are referenced
by OCID when custom Redis configuration is needed.

Excluded from v1:
  - defined_tags, system_tags -- managed by platform via freeform_tags
  - security_attributes -- specialized Oracle Zero-Trust Packet Routing
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.nodeCount` | `int32` |  |  |  |
| `spec.nodeMemoryInGbs` | `float` |  |  |  |
| `spec.softwareVersion` | `string` | yes |  |  |
| `spec.clusterMode` | `enum` |  |  |  |
| `spec.shardCount` | `int32` |  |  |  |
| `spec.nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.configSetId` | `string \| valueFrom` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the Redis cluster will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.subnetId

`string | valueFrom` · required

OCID of the subnet where the Redis cluster will be placed.
Changing this forces recreation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.nodeCount

`int32`

Number of nodes in the cluster. For non-sharded clusters this is the
total node count (1 primary + N-1 replicas). For sharded clusters
this is the number of nodes per shard. Updatable.

- rule: {"int32":{"gte":1}}

### spec.nodeMemoryInGbs

`float`

Memory allocated to each node in gigabytes. Determines the total
cache capacity. Updatable. Common values: 2, 4, 8, 16, 32.

- rule: {"float":{"gt":0}}

### spec.softwareVersion

`string` · required

OCI Cache engine version (e.g. "V7.0.5", "V7.1.1").
Available versions depend on the region. Updatable.

- rule: {"string":{"minLen":"1"}}

### spec.clusterMode

`enum`

Cluster topology mode. Changing this forces recreation.
When unset, OCI defaults to NONSHARDED.

Allowed values (use exactly as shown):

- `cluster_mode_unspecified`
- `nonsharded` -- Single shard with one primary and optional replicas.
- `sharded` -- Multiple shards, each with its own primary and replicas.

### spec.shardCount

`int32`

Number of shards in the cluster. Only applicable when cluster_mode
is sharded. Each shard has node_count nodes. Updatable.

### spec.nsgIds

`[]string | valueFrom`

OCIDs of network security groups controlling access to the cluster.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.configSetId

`string | valueFrom`

OCID of an OCI Cache Config Set providing custom Redis configuration
parameters (e.g. maxmemory-policy, timeout). When omitted, the
default configuration is used.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `shard_count_required_when_sharded`: shard_count must be greater than zero when cluster_mode is sharded

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciRedisCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | OCID of the Redis cluster. |
| `status.outputs.primary_fqdn` | `string` | FQDN of the primary (read-write) endpoint. This is the main connection point for non-sharded clusters. |
| `status.outputs.primary_endpoint_ip_address` | `string` | Private IP address of the primary endpoint. |
| `status.outputs.replicas_fqdn` | `string` | FQDN of the replica (read-only) endpoint. |
| `status.outputs.discovery_fqdn` | `string` | FQDN of the discovery endpoint for sharded clusters. Clients use this to discover shard topology. |

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
