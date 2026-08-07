---
title: "Kubernetes NodePool"
description: "Kubernetes NodePool deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceankubernetesnodepool"
---

# Kubernetes NodePool on DigitalOcean

Adds an additional node pool to an existing DigitalOcean Kubernetes (DOKS) cluster with configurable Droplet sizing, fixed or auto-scaling node counts, Kubernetes labels and taints for workload scheduling, and DigitalOcean tags for billing attribution. Integrates with Planton's Provider Connections for DigitalOcean credential management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes Node Pool** -- a `digitalocean_kubernetes_node_pool` resource attached to the specified DOKS cluster with the configured Droplet size and node count
- **Auto-Scaling Policy** -- created only when `autoScale` is `true`; allows the cluster autoscaler to manage node count between `minNodes` and `maxNodes`
- **Kubernetes Labels** -- applied to every node in the pool, enabling pod scheduling via `nodeSelector` and node affinity rules
- **Kubernetes Taints** -- applied to every node in the pool, preventing pods without matching tolerations from being scheduled
- **DigitalOcean Tags** -- applied to the underlying Droplets for cost attribution and organizational grouping

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **An existing DOKS cluster** to attach the node pool to. Provide the cluster name directly or reference a DigitalOceanKubernetesCluster Cloud Resource via ValueFromRef.
- **A valid Droplet size slug** available in the cluster's region (e.g., `s-4vcpu-8gb`, `s-2vcpu-4gb`, `g-8vcpu-32gb`). Check available sizes via the DigitalOcean CLI or API.

## Deploy

### Console

Open the deployment store, find **Kubernetes NodePool on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Autoscaling Production Node Pool** preset in the [Presets](#presets) tab for a production-ready configuration with auto-scaling.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: app-pool
  org: acme-corp
  env: prod
spec:
  nodePoolName: app-pool
  cluster:
    value: "app-cluster"
  size: s-4vcpu-8gb
  nodeCount: 3
  autoScale: true
  minNodes: 2
  maxNodes: 6
  labels:
    workload: app
```

```shell
planton apply -f node-pool.yaml
```

This creates an auto-scaling node pool with 3 initial nodes (scaling between 2 and 6) using general-purpose 4-vCPU/8-GB Droplets. Nodes are labeled for workload-based scheduling.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the node pool to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  cluster:
    valueFrom:
      kind: DigitalOceanKubernetesCluster
      name: app-cluster
      fieldPath: metadata.name
```

The InfraPipeline resolves the dependency graph, deploys the Kubernetes cluster first, then provisions the node pool on it.

## Key Configuration

These are the most important decisions when configuring a Kubernetes node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Droplet sizing** -- The `size` field sets the instance type for each node (e.g., `s-2vcpu-4gb` for system workloads, `s-4vcpu-8gb` for general applications, `g-8vcpu-32gb` for memory-intensive workloads). All nodes in a pool share the same size.

**Auto-scaling vs. fixed size** -- Set `autoScale` to `true` with `minNodes` and `maxNodes` for workloads with variable traffic. The cluster autoscaler adds or removes nodes based on pending pod scheduling. Use a fixed `nodeCount` without auto-scaling for system pools with predictable resource needs.

**Labels and taints** -- Use `labels` to enable `nodeSelector` or `nodeAffinity` rules that target specific pools (e.g., `workload: app`). Use `taints` to prevent pods from scheduling on dedicated nodes unless they explicitly tolerate the taint (e.g., `dedicated=system:NoSchedule` for infrastructure pools).

**DigitalOcean tags** -- The `tags` field applies billing tags to the underlying Droplets. These are separate from Kubernetes labels and are visible in DigitalOcean's billing and management console for cost attribution.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanKubernetesCluster** | `cluster` | `metadata.name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `node_pool_id` | UUID of the created node pool | DigitalOcean API operations, monitoring dashboards |
| `node_ids` | IDs of the individual Droplet nodes in the pool | Node-level monitoring, maintenance operations |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Autoscaling production pool** -- General-purpose nodes (4 vCPU, 8 GB) with auto-scaling between 2 and 6 nodes. Labeled for application workload scheduling. Suited for web and API workloads with variable traffic. Start from the **Autoscaling Production Node Pool** preset.

**Fixed-size system pool** -- Smaller nodes (2 vCPU, 4 GB) with a fixed count of 2 and a `dedicated=system:NoSchedule` taint. Dedicated to cluster infrastructure like ingress controllers, cert-manager, and monitoring agents. Start from the **Fixed-Size System Node Pool** preset.

## Works With

- [**Kubernetes Cluster on DigitalOcean**](/cloud-catalog/digital-ocean-kubernetes-cluster) -- provides the parent DOKS cluster that this node pool is attached to