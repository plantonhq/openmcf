---
title: "Kubernetes Node Pool"
description: "Kubernetes Node Pool deployment documentation"
icon: "package"
order: 100
componentName: "alicloudkubernetesnodepool"
---

# AliCloud Kubernetes Node Pool

Deploys a worker node pool for an ACK Managed Kubernetes cluster with configurable instance types, auto-scaling, spot pricing, managed lifecycle operations, Kubernetes labels and taints, and multi-AZ scheduling. Node pools have their own lifecycle and can be scaled independently of the cluster. The component integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to clusters, VSwitches, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **ACK Node Pool** -- an `alicloud_cs_kubernetes_node_pool` attached to the specified cluster, with the configured instance types, node count, system disk, data disks, and Kubernetes scheduling configuration
- **Auto Scaling Group** -- created automatically by ACK to manage the underlying ECS instances in the node pool; auto-scaling behavior is controlled by the `scalingConfig` settings
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **An ACK Managed Kubernetes cluster** -- the node pool attaches to an existing cluster. Provide the cluster ID directly or reference an AliCloudKubernetesCluster Cloud Resource via ValueFromRef.
- **1-5 VSwitches** in the same VPC as the cluster, preferably in distinct availability zones for high availability. Provide VSwitch IDs directly or reference AliCloudVswitch Cloud Resources via ValueFromRef.
- **A security group** (optional) -- if omitted, the cluster's default security group is used. Provide security group IDs directly or reference AliCloudSecurityGroup Cloud Resources via ValueFromRef.
- **An SSH key pair** (recommended) or password for node access. SSH key pairs are preferred for managed node pools.

## Deploy

### Console

Open the deployment store, find **AliCloud Kubernetes Node Pool**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **General Purpose Autoscaling** preset in the [Presets](#presets) tab to pre-populate a production-ready auto-scaling configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudKubernetesNodePool
metadata:
  name: general-workers
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  clusterId:
    value: "c-abc123"
  name: general-workers
  vswitchIds:
    - value: "vsw-abc001"
    - value: "vsw-abc002"
  instanceTypes:
    - ecs.g7.xlarge
  desiredSize: 3
```

```shell
planton apply -f node-pool.yaml
```

This creates a fixed-size node pool with 3 nodes using `ecs.g7.xlarge` instances, AliyunLinux3 OS, 120 GB cloud ESSD system disk, PostPaid billing, and the cluster's default security group. Auto-scaling, spot instances, managed lifecycle, labels, and taints are not configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the node pool to a cluster and VSwitches deployed in the same InfraPipeline:

```yaml
spec:
  clusterId:
    valueFrom:
      kind: AliCloudKubernetesCluster
      name: platform-cluster
      fieldPath: status.outputs.cluster_id
  vswitchIds:
    - valueFrom:
        kind: AliCloudVswitch
        name: worker-vswitch-a
        fieldPath: status.outputs.vswitch_id
    - valueFrom:
        kind: AliCloudVswitch
        name: worker-vswitch-b
        fieldPath: status.outputs.vswitch_id
  securityGroupIds:
    - valueFrom:
        kind: AliCloudSecurityGroup
        name: worker-sg
        fieldPath: status.outputs.security_group_id
```

The InfraPipeline resolves the dependency graph, deploys the cluster, VSwitches, and security group first, then provisions the node pool with the resolved values.

## Key Configuration

These are the most important decisions when configuring a node pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Auto-scaling vs fixed size** -- For variable workloads, configure `scalingConfig` with `enable: true`, `minSize`, and `maxSize`. The cluster auto-scaler adjusts node count based on pending pod resource requests. For predictable workloads, set `desiredSize` to a fixed count without `scalingConfig`. Use `multiAzPolicy` to control AZ distribution: `BALANCE` for even spread, `COST_OPTIMIZED` for cheapest AZ, or `PRIORITY` for sequential fill.

**Spot instances for cost savings** -- Set `spotStrategy` to `SpotAsPriceGo` (market price) or `SpotWithPriceLimit` (with per-instance-type price caps via `spotPriceLimits`). Specify multiple `instanceTypes` to spread across spot pools and improve availability. Spot instances can reduce costs up to 90% but may be reclaimed when capacity is needed.

**Managed lifecycle** -- Configure `management` with `enable: true` to let ACK automatically repair unhealthy nodes (`autoRepair`), upgrade kubelet versions (`autoUpgrade`), and patch vulnerabilities. Set `maxUnavailable` to control how many nodes can be offline during managed operations.

**Kubernetes scheduling** -- Use `labels` for pod affinity and node selectors (e.g., `{"workload-type": "compute"}`). Use `taints` with `NoSchedule`, `PreferNoSchedule`, or `NoExecute` effects to repel pods that lack matching tolerations. Set `cpuPolicy: static` for latency-sensitive workloads that benefit from CPU pinning.

**Disk configuration** -- The system disk defaults to 120 GB cloud ESSD. Configure `systemDisk` to change category, size, or performance level. Add `dataDisks` for additional storage (e.g., container image caching, local ephemeral storage). Each data disk supports independent category, size, and encryption settings.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudKubernetesCluster** | `clusterId` | `status.outputs.cluster_id` |
| **AliCloudVswitch** | `vswitchIds` | `status.outputs.vswitch_id` |
| **AliCloudSecurityGroup** (optional) | `securityGroupIds` | `status.outputs.security_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `node_pool_id` | ACK node pool ID | Monitoring dashboards, scaling activity queries |
| `scaling_group_id` | Auto Scaling group ID associated with the node pool | Scaling activity queries, node status monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**General-purpose with auto-scaling** -- An auto-scaling node pool with general-purpose instance types, managed lifecycle enabled, and `BALANCE` multi-AZ policy for even distribution across availability zones. Start from the **General Purpose Autoscaling** preset.

**Fixed-size for development** -- A small fixed-size node pool with minimal instance types, no auto-scaling, and PostPaid billing for cost-efficient development and testing. Start from the **Fixed Size Development** preset.

**Cost-optimized with spot** -- A spot instance node pool with multiple instance types, `SpotAsPriceGo` strategy, auto-scaling, and `COST_OPTIMIZED` multi-AZ policy. Suitable for fault-tolerant batch processing and CI/CD workloads. Start from the **Cost Optimized Spot** preset.

## Works With

- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- provides the ACK cluster that this node pool attaches to
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides the VSwitches for worker node network placement across availability zones
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- provides network access control for worker nodes