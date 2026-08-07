---
title: "Kapsule Pool"
description: "Kapsule Pool deployment documentation"
icon: "package"
order: 100
componentName: "scalewaykapsulepool"
---

# Scaleway Kapsule Pool

Deploys an additional node pool in an existing Scaleway Kapsule cluster with configurable instance types, autoscaling, autohealing, Kubernetes labels, taints, and upgrade policies. Node pools provide dedicated compute capacity for workload isolation -- use separate pools for different instance types, GPU workloads, or teams with distinct scheduling requirements. Supports ValueFromRef for cluster dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kapsule Node Pool** -- a `kubernetes.Pool` attached to the referenced cluster with the configured instance type, node count, and optional autoscaling bounds
- **Kubernetes Labels** -- applied to all nodes in the pool via Scaleway's Cloud Controller Manager tag convention, enabling pod scheduling with `nodeSelector` and node affinity rules
- **Kubernetes Taints** -- applied to all nodes in the pool via CCM tags, preventing pods without matching tolerations from being scheduled on these nodes
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) plus CCM-formatted label and taint tags applied automatically

## Before You Deploy

### Scaleway Account

- **An existing Kapsule cluster** in the target region. Provide the cluster ID directly or reference a ScalewayKapsuleCluster Cloud Resource via ValueFromRef. The pool's region must match the cluster's region.
- **Sufficient instance quota** in the target zone for the requested instance type and node count. Check Scaleway quotas in the console if provisioning fails with capacity errors.

## Deploy

### Console

Open the deployment store, find **Scaleway Kapsule Pool**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **General-Purpose** preset in the [Presets](#presets) tab for a fixed-size pool with GP1-XS instances.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayKapsulePool
metadata:
  name: workers
  org: acme-corp
  env: prod
spec:
  region: fr-par
  clusterId:
    value: "fr-par/abc12345-6789-def0-1234-567890abcdef"
  nodeType: GP1-XS
  size: 3
  autohealing: true
  publicIpDisabled: true
```

```shell
planton apply -f scaleway-kapsule-pool.yaml
```

This creates a fixed-size 3-node pool with GP1-XS instances, autohealing enabled, and no public IPs. No autoscaling, labels, or taints are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the pool to a cluster deployed in the same InfraPipeline:

```yaml
spec:
  clusterId:
    valueFrom:
      kind: ScalewayKapsuleCluster
      name: app-cluster
      fieldPath: status.outputs.cluster_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC, Private Network, and cluster first, then provisions the node pool with the resolved cluster ID.

## Key Configuration

These are the most important decisions when configuring a node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance type** -- The `nodeType` field determines CPU, RAM, and local storage per node. Use `DEV1-M` for development, `GP1-XS` for general workloads, or `PRO2-S`/`PRO2-M` for production with guaranteed resources. Cannot be changed after creation -- to change instance types, create a new pool and migrate workloads.

**Autoscaling** -- Set `autoScale` to true with `minSize` and `maxSize` bounds for elastic workloads. The cluster-level `autoscalerConfig` on the parent ScalewayKapsuleCluster controls autoscaler behavior (scale-down delays, utilization thresholds). When autoscaling is enabled, the `size` field sets only the initial node count.

**Labels and taints** -- Use `kubernetesLabels` to tag nodes for scheduling (e.g., `{"pool": "gpu", "team": "ml"}`). Use `taints` to restrict scheduling to pods with matching tolerations (e.g., `NoSchedule` taints for dedicated GPU pools). Labels and taints are synced to Kubernetes nodes via Scaleway's Cloud Controller Manager.

**Public IP** -- Set `publicIpDisabled` to true for production pools so nodes are reachable only via the Private Network. Requires a Public Gateway or NAT for outbound internet access to container registries and APIs.

**Upgrade policy** -- Configure `upgradePolicy` with `maxSurge` and `maxUnavailable` to control rolling update behavior during Kubernetes version upgrades. The default (0 surge, 1 unavailable) replaces one node at a time. Set `maxSurge: 1` for zero-downtime upgrades at the cost of temporarily running an extra node.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayKapsuleCluster** | `clusterId` | `status.outputs.cluster_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `pool_id` | Unique identifier of the created node pool | Monitoring, management API calls, distinguishing pools |
| `pool_version` | Kubernetes version running on pool nodes | Version tracking, upgrade verification |
| `current_size` | Actual number of nodes in the pool | Capacity monitoring, autoscaler verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**General-purpose pool** -- A fixed 3-node pool with GP1-XS instances, autohealing, and private-only nodes. The standard additional pool for extending a cluster with general workloads. Start from the **General-Purpose** preset.

**Autoscaling worker pool** -- PRO2-M instances scaling between 1 and 8 nodes with a surge-based upgrade policy for zero-downtime rolling updates. The standard pattern for elastic production workloads, batch jobs, and queue consumers. Start from the **Autoscaling Workers** preset.

## Works With

- [**Scaleway Kapsule Cluster**](/cloud-catalog/scaleway-kapsule-cluster) -- provides the cluster that this node pool belongs to