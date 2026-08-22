# DigitalOcean Kubernetes Node Pool

An additional worker pool for an existing DigitalOcean Kubernetes (DOKS) cluster, described once in a Planton manifest: node sizing and count, autoscaling bounds, Kubernetes labels and taints, DigitalOcean tags, and AMD GPU partitioning. The cluster's own default pool belongs to the `DigitalOceanKubernetesCluster` kind; this kind grows a cluster with separately shaped pools.

## What this component models

The spec maps one-to-one onto DigitalOcean's standalone node pool:

| Spec field | What it controls |
|---|---|
| `nodePoolName` | The pool's name, unique within the cluster |
| `cluster` | The owning DOKS cluster — a literal UUID or a reference to a `DigitalOceanKubernetesCluster`; create-only |
| `size` | Droplet size slug for every node (`s-2vcpu-4gb`, GPU slugs, ...); changing it replaces the pool |
| `nodeCount` | Node count; with autoscaling on, only the initial count |
| `autoScale` / `minNodes` / `maxNodes` | DigitalOcean's cluster-autoscaler manages the count between the bounds |
| `labels` | Kubernetes node labels for scheduling (nodeSelector, affinity) |
| `taints` | Kubernetes taints keeping untolerating pods off the nodes |
| `tags` | DigitalOcean tags on the pool's Droplets — billing attribution and grouping |
| `gpuPartitionMode` | AMD GPU partitioning for GPU sizes; changing it replaces the pool |

## Quick start

A fixed one-node pool on an existing cluster:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanKubernetesNodePool
metadata:
  name: app-pool
spec:
  nodePoolName: app-pool
  cluster:
    valueFrom:
      kind: DigitalOceanKubernetesCluster
      name: my-doks-cluster
      fieldPath: status.outputs.cluster_id
  size: s-2vcpu-4gb
  nodeCount: 2
```

```shell
planton apply -f app-pool.yaml
```

## Outputs

Both provisioners export the identical output set:

| Output | Description |
|---|---|
| `node_pool_id` | The pool's UUID (import id for `digitalocean_kubernetes_node_pool`) |
| `cluster_id` | The owning cluster's UUID — the pool's API address needs both ids |
| `node_ids` | The DOKS node object UUIDs of the pool's current members |
| `droplet_ids` | The integer ids of the Droplets backing the nodes — wire Droplet-scoped resources (e.g. firewalls) to the pool's machines |

## Behavior worth knowing

- **`cluster` and `size` replace the pool when changed** (the nodes are recreated); everything else updates in place.
- **Autoscaling bounds are validated early**: `autoScale: true` requires `minNodes >= 1` and `maxNodes >= minNodes` — the API would reject them late, the spec rejects them at validation.
- **With autoscaling on, the live node count drifts by design.** `nodeCount` is only the initial count; the provider suppresses the diff while the count sits between the bounds.
- **Taints must spell their effect exactly as Kubernetes does** (`NoSchedule`, `PreferNoSchedule`, `NoExecute`), and a taint's `value` may be empty — Kubernetes allows valueless taints.
- **A cluster's default pool cannot be managed here** — it is part of the cluster resource itself, and DigitalOcean refuses to import a default pool as a standalone one.
- **Pulumi SDK v4.49.0 cannot express `gpuPartitionMode`.** The Pulumi module fails loudly if it is set; Terraform wires it. See the [GUIDE](GUIDE.md).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
