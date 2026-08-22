# Kubernetes Node Pool on DigitalOcean

Adds a worker pool to an existing DigitalOcean Kubernetes (DOKS) cluster with configurable node sizing, fixed or autoscaled node counts, Kubernetes labels and taints, DigitalOcean tags, and AMD GPU partitioning. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Kubernetes Node Pool** -- a worker pool attached to the referenced DOKS cluster, with the configured Droplet size and node count
- **Autoscaling** -- enabled by `autoScale`, with DigitalOcean's cluster-autoscaler managing the node count between `minNodes` and `maxNodes`
- **Kubernetes Node Labels and Taints** -- user labels plus the standard Planton identity labels on every node; taints applied for workload isolation
- **DigitalOcean Tags** -- user-supplied tags plus the standard Planton labels on the pool's Droplets

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A DOKS cluster** -- the pool attaches to an existing cluster; reference a `DigitalOceanKubernetesCluster` Cloud Resource via ValueFromRef or provide the cluster UUID directly (`doctl kubernetes cluster list`).
- **A valid Droplet size slug** (e.g., `"s-2vcpu-4gb"`) -- check DOKS-capable sizes via `doctl kubernetes options sizes`.

## Deploy

### Console

Open the deployment store, find **Kubernetes Node Pool on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Autoscaling Production** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: app-pool
  org: acme-corp
  env: prod
spec:
  nodePoolName: app-pool
  cluster:
    value: "fb7d9b81-fe06-4ee5-87f1-b9efc5af46fd"
  size: s-4vcpu-8gb
  nodeCount: 3
  autoScale: true
  minNodes: 2
  maxNodes: 6
```

```shell
planton apply -f node-pool.yaml
```

This attaches an autoscaling pool to the cluster. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the pool to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  cluster:
    valueFrom:
      kind: DigitalOceanKubernetesCluster
      name: platform-cluster
      fieldPath: status.outputs.cluster_id
```

The InfraPipeline resolves the dependency graph, deploys the cluster first, then attaches the pool.

## Key Configuration

These are the most important decisions when configuring a node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Sizing** -- The `size` field sets every node's CPU and memory. Changing it later replaces the whole pool (nodes are recreated), so plan capacity classes as separate pools rather than resizing one in place.

**Fixed vs. autoscaled** -- A fixed pool holds exactly `nodeCount` nodes. With `autoScale: true`, `nodeCount` is only the initial count and DigitalOcean's cluster-autoscaler moves it between `minNodes` and `maxNodes`; the live count drifting is normal and produces no configuration diff.

**Labels and taints** -- Labels make the pool targetable from Kubernetes (nodeSelector, affinity); taints keep untolerating pods off. Pair them for dedicated pools: a taint alone isolates, a label alone only attracts.

**GPU partitioning** -- `gpuPartitionMode` splits supported AMD GPU sizes into partitions. It is create-time-only in effect: changing it replaces the pool. Currently Terraform-only (Pulumi SDK gap; the Pulumi provisioner fails loudly if set).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanKubernetesCluster** (required) | `cluster` | `status.outputs.cluster_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `node_pool_id` | The pool's UUID | Import addressing, pool-scoped automation |
| `cluster_id` | The owning cluster's UUID | Anything addressing the pool through the cluster API |
| `node_ids` | DOKS node object UUIDs of the current members | Node-level automation against the DOKS API |
| `droplet_ids` | Integer ids of the Droplets backing the nodes | Firewall rules and other Droplet-scoped wiring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Autoscaling application pool** -- general-purpose nodes scaling 2–6 with a `workload: app` label. Start from the **Autoscaling Production** preset.

**Dedicated system pool** -- a fixed two-node pool with a `NoSchedule` taint reserving it for system components. Start from the **Fixed Size** preset.

## Works With

- [**DigitalOcean Kubernetes Cluster**](/cloud-catalog/digital-ocean-kubernetes-cluster) -- the cluster this pool attaches to
- [**DigitalOcean Firewall**](/cloud-catalog/digital-ocean-firewall) -- secures the pool's Droplets by id or tag
