---
title: "Redis Cluster"
description: "Redis Cluster deployment documentation"
icon: "package"
order: 100
componentName: "scalewayrediscluster"
---

# Scaleway Redis Cluster

Deploys a managed Redis cluster on Scaleway with configurable deployment modes -- standalone, high-availability, or sharded cluster -- determined by node count. Supports TLS encryption, Private Network integration (mutually exclusive with ACL rules), and engine-level tuning. Supports ValueFromRef for Private Network dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Redis Cluster** -- a managed in-memory data store with the configured version, node type, cluster size (standalone, HA, or sharded), authentication credentials, optional TLS, and either ACL rules or Private Network attachment
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway Private Network** in the target region when using private connectivity. Redis clusters attached to a Private Network cannot use ACL rules -- these are mutually exclusive on Scaleway. Provide the Private Network UUID directly or reference a ScalewayPrivateNetwork Cloud Resource via ValueFromRef.
- **Redis is a zonal resource** -- specify an availability zone (e.g., `fr-par-1`, `nl-ams-1`), not a region. The zone determines which data center hosts the cluster nodes.

## Deploy

### Console

Open the deployment store, find **Scaleway Redis Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production HA** preset in the [Presets](#presets) tab for a 3-node TLS-enabled cluster with Private Network connectivity.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayRedisCluster
metadata:
  name: app-cache
  org: acme-corp
  env: prod
spec:
  zone: fr-par-1
  version: "7.2.5"
  nodeType: RED1-MICRO
  clusterSize: 1
  userName: redis-user
  password: changeme123
```

```shell
planton apply -f scaleway-redis-cluster.yaml
```

This creates a standalone Redis 7.2 cluster with a single RED1-MICRO node. No TLS, Private Network, or ACL rules are configured -- the public endpoint is accessible from all IPs.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the cluster to a Private Network deployed in the same InfraPipeline:

```yaml
spec:
  privateNetworkId:
    valueFrom:
      kind: ScalewayPrivateNetwork
      name: app-network
      fieldPath: status.outputs.private_network_id
```

The InfraPipeline resolves the dependency graph, deploys the Private Network first, then provisions the Redis cluster with the resolved Private Network ID.

## Key Configuration

These are the most important decisions when configuring a Redis cluster. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Cluster size** -- The `clusterSize` field determines the deployment mode: 1 for standalone, 2 for HA (main + standby with automatic failover), 3+ for cluster mode (sharding). Scaling from standalone to cluster mode destroys and recreates the cluster. Scaling up within cluster mode is an online migration.

**TLS** -- Set `tlsEnabled` to true for encrypted client connections. The cluster's TLS certificate is exported in `status.outputs.certificate`. Changing this value destroys and recreates the cluster -- plan TLS requirements before initial deployment.

**Networking** -- ACL rules and Private Network are mutually exclusive on Scaleway. Use `aclRules` for public endpoint access control or `privateNetworkId` for private-only connectivity, but not both. In cluster mode (3+ nodes), the Private Network cannot be changed after creation.

**Node type** -- Choose from `RED1-MICRO` (development) through `RED1-XL` (high-traffic production). Node type can be upgraded via online migration but never downgraded.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayPrivateNetwork** (optional) | `privateNetworkId` | `status.outputs.private_network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `cluster_id` | Unique identifier of the Redis cluster | Monitoring dashboards, management automation |
| `public_network_port` | Public endpoint TCP port | Application connection configuration (public mode only) |
| `public_network_ips` | Public endpoint IPv4 addresses | Application connection strings (public mode only) |
| `private_network_port` | Private Network endpoint TCP port | Application connection configuration (private mode only) |
| `private_network_ips` | Private Network endpoint IPv4 addresses | Application connection strings (private mode only) |
| `certificate` | TLS CA certificate in PEM format | Client TLS verification (only when `tlsEnabled` is true) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development standalone** -- A single-node RED1-MICRO cluster with public ACL access and no TLS. Minimal cost for development caching and session storage testing. Start from the **Dev Standalone** preset.

**Production HA** -- A 3-node RED1-M cluster with TLS encryption and Private Network connectivity. Automatic failover for production caching, session stores, and real-time data. Start from the **Production HA** preset.

## Works With

- [**Scaleway Private Network**](/cloud-catalog/scaleway-private-network) -- provides network isolation for private-only Redis endpoints