# ScalewayRedisCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayRedisClusterSpec defines the specification for a Scaleway Managed
Redis cluster.

Scaleway Redis provides fully managed, in-memory data stores ideal for
caching, session management, real-time analytics, and message brokering.

This is a **standalone resource** wrapping a single `scaleway_redis_cluster`
Terraform resource. Unlike ScalewayRdbInstance (which bundles 5 sub-resources),
Redis clusters are self-contained: ACL rules and Private Network configuration
are inline properties of the cluster resource itself.

Redis clusters are **zonal** resources (e.g., "fr-par-1", "nl-ams-1"),
unlike RDB instances which are regional.

**Cluster sizing determines deployment mode:**
  - cluster_size = 1: Standalone (single node, no redundancy).
  - cluster_size = 2: High Availability (1 main + 1 standby, automatic failover).
  - cluster_size >= 3: Cluster mode (sharding across multiple nodes).

**Networking constraint:** ACL rules and Private Network are mutually
exclusive. Scaleway does not support ACL when the cluster is attached to
a Private Network. This constraint is enforced via CEL validation below
and by the Terraform/Pulumi providers at apply time.

**Composition pattern:** The cluster accepts a Private Network reference
via `StringValueOrRef`, enabling private connectivity from applications
and other resources on the same network. Downstream resources can
reference `status.outputs.cluster_id` for monitoring or dependent services.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zone` | `string` | yes |  |  |
| `spec.version` | `string` | yes |  |  |
| `spec.nodeType` | `string` | yes |  |  |
| `spec.clusterSize` | `uint32` |  | `1` |  |
| `spec.tlsEnabled` | `bool` |  |  |  |
| `spec.userName` | `string` | yes |  |  |
| `spec.password` | `string` (sensitive) | yes |  |  |
| `spec.aclRules` | `[]ScalewayRedisAclRule` |  |  |  |
| `spec.aclRules[].ip` | `string` | yes |  |  |
| `spec.aclRules[].description` | `string` |  |  |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.settings` | `map<string, string>` |  |  |  |

## Field Details

### spec.zone

`string` · required

The Scaleway availability zone where the cluster will be created.

Examples: "fr-par-1", "nl-ams-1", "pl-waw-1"

Redis is a ZONAL resource (unlike RDB which is regional). The zone
determines which data center hosts the cluster nodes.

Cannot be changed after creation.

- rule: {"required":true}

### spec.version

`string` · required

Redis engine version.

Format: semantic version (e.g., "7.2.5", "6.2.7").

Can be upgraded (triggers an online migration) but never downgraded.
Check Scaleway documentation for currently supported versions.

- rule: {"required":true,"string":{"pattern":"^[0-9]+\\.[0-9]+\\.[0-9]+$"}}

### spec.nodeType

`string` · required

Node type determining CPU, RAM, and performance characteristics.

Common types:
  - "RED1-MICRO"  -- Smallest, for development and testing.
  - "RED1-S"      -- Small production workloads.
  - "RED1-M"      -- Medium production workloads.
  - "RED1-L"      -- Large production workloads.
  - "RED1-XL"     -- High-traffic production.

Can be upgraded (triggers an online migration) but never downgraded.

- rule: {"required":true}

### spec.clusterSize

`uint32`

Number of nodes in the cluster. Determines the deployment mode:

  - 1: Standalone mode (single node, no redundancy).
  - 2: High Availability (1 main + 1 standby with automatic failover).
  - 3+: Cluster mode (data sharded across nodes, minimum 3).

IMPORTANT lifecycle constraints:
  - Standalone (1) to Cluster (3+): DESTROYS and recreates the cluster.
  - Cluster to smaller Cluster: DESTROYS and recreates the cluster.
  - Cluster to larger Cluster: online migration (safe scale-out).
  - HA adjustments (1 to 2, 2 to 1): migration-based.

Default: 1 (standalone).

- default: `1`

### spec.tlsEnabled

`bool`

Whether to enable TLS encryption for client connections.

When enabled, all client connections must use TLS. The cluster's
TLS certificate is exported in `status.outputs.certificate`.

IMPORTANT: Changing this value DESTROYS and RECREATES the cluster.
Plan TLS requirements before initial deployment.

Default: false.

### spec.userName

`string` · required

Username for the Redis cluster's initial user.

This is the only user created with the cluster. Redis on Scaleway
supports a single authentication principal (unlike RDB which has
multi-user support).

- rule: {"required":true,"string":{"maxLen":"63"}}

### spec.password

`string` · required · sensitive

Password for the Redis cluster user.

Must meet minimum complexity requirements. For production, use a
strong, randomly generated password and manage it through your
organization's secrets workflow.

- rule: {"required":true,"string":{"minLen":"8"}}

### spec.aclRules

`[]ScalewayRedisAclRule`

Network ACL rules restricting which IPs can connect to the cluster.

Each rule allows a CIDR range to connect. The complete set of rules
defines who can reach the cluster over its public endpoint.

IMPORTANT: ACL rules CONFLICT with Private Network attachment.
Scaleway does not support both simultaneously. If `private_network_id`
is set, this field MUST be empty. This constraint is enforced by
CEL validation below and by the Scaleway API at apply time.

If neither ACL nor Private Network is configured, Scaleway creates
a public endpoint with no restrictions (all IPs allowed).

### spec.aclRules[].ip

`string` · required

CIDR range to allow.

Examples:
  - "10.0.0.0/24" -- Allow a /24 subnet
  - "1.2.3.4/32"  -- Allow a single IP
  - "0.0.0.0/0"   -- Allow all IPs (NOT recommended for production)

- rule: {"required":true}

### spec.aclRules[].description

`string`

Human-readable description for this rule.

Helps operators understand the purpose of each ACL entry.
Examples: "Office IP", "VPN egress", "CI/CD pipeline"

### spec.privateNetworkId

`string | valueFrom`

Private Network to attach the cluster to.

When set, the cluster receives a private endpoint reachable only
from resources on the same Private Network. No public endpoint is
created.

IMPORTANT: Private Network CONFLICTS with ACL rules. If this
field is set, `acl_rules` MUST be empty. This constraint is
enforced by CEL validation below and by the Scaleway API.

IPAM is used automatically to assign IP addresses. Static IP
assignment is not exposed in this version.

In Cluster mode (3+ nodes), the Private Network is set at creation
time and CANNOT be changed afterward without destroying and
recreating the cluster.

In infra charts, this is typically wired via valueFrom:

  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.settings

`map<string, string>`

Redis-specific configuration settings.

Key-value pairs passed to the Redis engine configuration. Applied
on both creation and updates.

Common settings:
  - "maxclients" = "1000"        -- Maximum concurrent client connections
  - "tcp-keepalive" = "120"      -- TCP keepalive interval in seconds
  - "maxmemory-policy" = "allkeys-lru" -- Eviction policy when memory is full
  - "timeout" = "300"            -- Client idle timeout in seconds

Available settings depend on the Redis version. Use the Scaleway
API or CLI to list available settings for a given version.

Optional. If empty, Scaleway uses Redis defaults optimized for
the node type.

## Validation Rules

- `acl_private_network_mutual_exclusivity`: acl_rules and private_network_id are mutually exclusive -- Scaleway does not support ACL when the cluster is on a Private Network

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayRedisCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_id` | `string` | The unique identifier of the created Redis cluster. Format: zonal ID (e.g., "fr-par-1/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"). This is the primary output referenced by downstream resources. |
| `status.outputs.public_network_port` | `uint32` | Public network port number. The TCP port for connecting to the cluster over the public network. Populated when the cluster is NOT attached to a Private Network. Zero when using Private Network mode. |
| `status.outputs.public_network_ips` | `[]string` | Public network IP addresses. The IPv4 addresses for connecting to the cluster over the public network. Populated when the cluster is NOT attached to a Private Network. Empty when using Private Network mode. |
| `status.outputs.private_network_port` | `uint32` | Private network port number. The TCP port for connecting to the cluster from the Private Network. Populated when the cluster IS attached to a Private Network. Zero when using public network mode. |
| `status.outputs.private_network_ips` | `[]string` | Private network IP addresses. The IPv4 addresses for connecting to the cluster from the Private Network. Populated when the cluster IS attached to a Private Network. Empty when using public network mode. |
| `status.outputs.certificate` | `string` | TLS certificate in PEM format for verifying the Redis server. Available only when `tls_enabled` is true in the spec. Clients should use this CA certificate to establish encrypted connections and verify the server's identity. Empty when TLS is disabled. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](../README.md)
