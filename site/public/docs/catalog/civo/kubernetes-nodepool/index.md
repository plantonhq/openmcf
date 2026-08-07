---
title: "Kubernetes NodePool"
description: "Kubernetes NodePool deployment documentation"
icon: "package"
order: 100
componentName: "civokubernetesnodepool"
---

# Kubernetes NodePool on Civo

Adds a dedicated node pool to an existing Civo Kubernetes cluster, giving you independent control over instance sizing, node count, and auto-scaling for specific workloads. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Node Pool** -- a `civo_kubernetes_node_pool` resource attached to the referenced Civo Kubernetes cluster with the configured instance size and node count
- **Civo Labels** -- metadata labels derived from the resource identity, applied to the node pool for organizational tracking

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **An existing Civo Kubernetes cluster** -- provide the cluster name directly or reference a CivoKubernetesCluster Cloud Resource via ValueFromRef.
- **A valid instance size slug** (e.g., `g4s.kube.medium`) -- check available sizes via the Civo CLI (`civo sizes ls`) or Civo dashboard.

## Deploy

### Console

Open the deployment store, find **Kubernetes NodePool on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Autoscaling** preset in the [Presets](#presets) tab for a production-ready auto-scaling pool.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoKubernetesNodePool
metadata:
  name: workers
  org: acme-corp
  env: prod
spec:
  nodePoolName: workers
  cluster:
    value: "app-cluster"
  size: g4s.kube.medium
  nodeCount: 3
```

```shell
planton apply -f civo-node-pool.yaml
```

This creates a fixed 3-node pool with medium-sized instances in the referenced cluster. No auto-scaling is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the node pool to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  cluster:
    valueFrom:
      kind: CivoKubernetesCluster
      name: app-cluster
      fieldPath: metadata.name
```

The InfraPipeline resolves the dependency graph, deploys the Kubernetes cluster first, then provisions the node pool on it.

## Key Configuration

These are the most important decisions when configuring a Civo Kubernetes node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Node sizing** -- The `size` field sets the instance type for each node (e.g., `g4s.kube.small` for development, `g4s.kube.medium` for production). Larger sizes provide more CPU and memory per node but at higher cost. Match the size to your workload's resource requests.

**Node count** -- The `nodeCount` field sets how many nodes to provision initially. Use at least 3 for production workloads to allow pod scheduling across nodes during rolling updates and node failures.

**Auto-scaling** -- Set `autoScale` to `true` and configure `minNodes` and `maxNodes` to let the cluster autoscaler manage node count based on pending pod demand. The pool starts at `nodeCount` and scales within the configured bounds. Disable for workloads with steady resource requirements.

**Tags** -- The `tags` field applies organizational tags to the node pool within Civo, useful for cost tracking and filtering across multiple pools in the same cluster.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoKubernetesCluster** | `cluster` | `metadata.name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `node_pool_id` | Unique identifier of the created node pool | Civo API operations, monitoring dashboards |
| `node_ids` | IDs of the individual nodes in the pool | Node-level monitoring, targeted operations |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Autoscaling pool** -- 2 to 5 medium-sized nodes with auto-scaling enabled. The cluster autoscaler adjusts node count based on pending pod demand, scaling down during idle periods for cost optimization. Start from the **Autoscaling** preset.

**Fixed-size pool** -- Static 3-node pool with no auto-scaling. Predictable capacity and cost for workloads with steady resource requirements, such as stateful services that should not be disrupted by scale-down events. Start from the **Fixed-Size** preset.

## Works With

- [**Kubernetes Cluster on Civo**](/cloud-catalog/civo-kubernetes-cluster) -- the parent cluster to which this node pool is added