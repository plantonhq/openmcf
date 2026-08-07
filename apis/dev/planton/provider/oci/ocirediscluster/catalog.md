# Redis Cluster on OCI

Deploys an Oracle Cloud Infrastructure Cache (Redis) cluster -- a fully managed, Redis-compatible in-memory caching service supporting both non-sharded and sharded deployment modes. Non-sharded clusters provide a single primary with optional replicas for high availability. Sharded clusters distribute data across multiple shards for horizontal scaling, each with its own primary and replicas. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Redis Cluster** -- a `redis.RedisCluster` in the specified compartment and subnet with configurable node count, memory per node, software version, cluster topology (sharded or non-sharded), and network security group associations
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the cluster

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the cluster in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- A private subnet for cluster placement. The cluster runs on private IPs within the subnet. Provide the subnet OCID directly or reference an OciSubnet Cloud Resource via ValueFromRef. Changing the subnet forces recreation.
- Optionally, one or more network security groups to control inbound/outbound traffic to the cluster endpoints.

## Deploy

### Console

Open the deployment store, find **Redis Cluster on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Non-Sharded Cluster** preset in the [Presets](#presets) tab to pre-populate a 3-node HA cluster with read replicas.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciRedisCluster
metadata:
  name: app-cache
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  subnetId:
    value: "ocid1.subnet.oc1..example"
  nodeCount: 3
  nodeMemoryInGbs: 8
  softwareVersion: V7.1.1
  clusterMode: nonsharded
```

```shell
planton apply -f redis-cluster.yaml
```

This creates a non-sharded Redis 7.1.1 cluster with 3 nodes (1 primary + 2 replicas) and 8 GB memory per node. No network security groups or custom config set are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a compartment, subnet, and security group deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: app-compartment
      fieldPath: status.outputs.compartmentId
  subnetId:
    valueFrom:
      kind: OciSubnet
      name: cache-subnet
      fieldPath: status.outputs.subnetId
  nsgIds:
    - valueFrom:
        kind: OciSecurityGroup
        name: cache-nsg
        fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, subnet, and security group first, then provisions the Redis cluster with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Redis cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster mode** -- Set `clusterMode` to `nonsharded` for a single-shard cluster with one primary and optional replicas (suitable for most caching workloads). Set to `sharded` for data distribution across multiple shards (each with its own primary and replicas) when a single node's memory is insufficient. Changing `clusterMode` forces recreation.

**Node count and sharding** -- For non-sharded clusters, `nodeCount` is the total number of nodes (1 primary + N-1 replicas). For sharded clusters, `nodeCount` is the number of nodes per shard, and `shardCount` sets the number of shards. Total nodes in a sharded cluster = `nodeCount` x `shardCount`. Both are updatable after creation.

**Memory per node** -- Set `nodeMemoryInGbs` to the memory allocated to each node. Common values: 2, 4, 8, 16, 32 GB. Total cache capacity = `nodeMemoryInGbs` x total nodes (non-sharded) or `nodeMemoryInGbs` x `nodeCount` x `shardCount` (sharded). Updatable.

**Software version** -- Set `softwareVersion` to the desired Redis version (e.g., `V7.0.5`, `V7.1.1`). Available versions depend on the region. Updatable -- OCI performs rolling upgrades.

**Config set** -- Set `configSetId` to reference a custom OCI Cache Config Set for Redis configuration parameters (e.g., `maxmemory-policy`, `timeout`). When omitted, the default configuration is used.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciSubnet** | `subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | OCID of the Redis cluster | Monitoring, IAM policy scoping, resource management |
| `primary_fqdn` | FQDN of the primary (read-write) endpoint | Application connection configuration for non-sharded clusters |
| `primary_endpoint_ip_address` | Private IP address of the primary endpoint | Firewall rules, direct IP-based connections |
| `replicas_fqdn` | FQDN of the replica (read-only) endpoint | Read-replica routing for non-sharded clusters |
| `discovery_fqdn` | FQDN of the discovery endpoint for sharded clusters | Client shard topology discovery for sharded clusters |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Non-sharded cluster** -- A 3-node Redis 7.1.1 cluster (1 primary + 2 replicas) with 8 GB per node and NSG-based access control. The standard configuration for application caching with high availability and read scaling. Start from the **Non-Sharded Cluster** preset.

**Sharded cluster** -- A 3-shard Redis 7.1.1 cluster with 3 nodes per shard (9 nodes total), 16 GB per node, and NSG-based access control. Designed for large-scale caching where data exceeds single-node memory capacity. Start from the **Sharded Cluster** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this Redis cluster
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the private subnet for cluster placement
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for cluster access control