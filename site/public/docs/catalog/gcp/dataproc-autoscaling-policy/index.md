---
title: "Dataproc Autoscaling Policy"
description: "Dataproc Autoscaling Policy deployment documentation"
icon: "package"
order: 100
componentName: "gcpdataprocautoscalingpolicy"
---

# GCP Dataproc Autoscaling Policy

Deploys a reusable Dataproc autoscaling policy: the scaling contract that governs how attached clusters add primary and secondary (spot) workers when YARN memory is pending and remove them when it is idle. One policy can govern many clusters — a platform team tunes scaling behavior in one place, and editing the policy re-tunes every attached cluster in place. The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Dataproc Autoscaling Policy** -- a regional policy resource with the configured worker bounds, capacity weights, and YARN scaling algorithm
- **Primary Worker Bounds** -- the min/max envelope and capacity weight for the on-demand worker group (the stable HDFS base)
- **Secondary Worker Bounds** -- created only when `secondaryWorkerConfig` carries non-zero bounds; the envelope and weight for the preemptible/spot group (the cost-optimized burst arm)
- **YARN Scaling Algorithm** -- the scale-up/scale-down factors (what fraction of suggested capacity change is applied per evaluation), the graceful decommission window running tasks get before scale-down removes their worker, deadband fractions, and the evaluation cooldown

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the policy will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Dataproc API** enabled in the target project.
- **Region co-location** -- a Dataproc cluster can only attach policies in its own region; create the policy in the region your clusters run in.

## Deploy

### Console

Open the deployment store, find **GCP Dataproc Autoscaling Policy**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Balanced Autoscaling** preset in the [Presets](#presets) tab for a moderate general-purpose policy.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpDataprocAutoscalingPolicy
metadata:
  name: balanced-autoscaling
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  policyId: balanced-autoscaling
  location: us-central1
  workerConfig:
    maxInstances: 10
    minInstances: 2
  basicAlgorithm:
    yarnConfig:
      gracefulDecommissionTimeout: "1800s"
      scaleUpFactor: 0.5
      scaleDownFactor: 0.5
```

```shell
planton apply -f autoscaling-policy.yaml
```

A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, deploy the policy before the clusters that attach it:

```yaml
# In a GcpDataprocCluster spec:
clusterConfig:
  autoscalingPolicyUri:
    valueFrom:
      kind: GcpDataprocAutoscalingPolicy
      name: balanced-autoscaling
      fieldPath: status.outputs.name
```

The InfraPipeline resolves the dependency graph, deploys the policy first, then provisions the cluster with the policy attached.

## Key Configuration

These are the most important decisions when configuring an autoscaling policy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Scale factors** -- `scaleUpFactor` and `scaleDownFactor` (both required, 0.0-1.0) set the policy's personality: the fraction of YARN's suggested capacity change applied per evaluation. 1.0 reacts at maximum speed; 0.05 creeps smoothly. An explicit `scaleDownFactor: 0.0` disables scale-down entirely — pair it with the cluster's idle-delete TTL.

**Worker bounds and weights** -- `workerConfig.maxInstances` (required) is the primary group's hard ceiling; `minInstances: 0` accepts the API's default floor of 2. The weight pair splits new capacity between on-demand primaries and spot secondaries — weight 1:3 sends ~75% of new nodes to the spot group.

**Graceful decommission** -- `gracefulDecommissionTimeout` (required, e.g. `"1800s"`) is the window running tasks get to finish before scale-down removes their worker. Short windows kill long-running tasks.

**Cooldown** -- `basicAlgorithm.cooldownPeriod` sets the evaluation cadence (blank keeps GCP's `"120s"` default). Longer cooldowns smooth scaling further but slow every reaction.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `name` | Fully qualified policy resource name (`projects/{p}/locations/{l}/autoscalingPolicies/{id}`) | The exact value a GcpDataprocCluster's `autoscalingPolicyUri` attaches |
| `policy_id` | The short policy ID | Display, logging |
| `location` | The policy's region | Region co-location checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Balanced autoscaling** -- Half-strength factors, an even primary/secondary split, and a 30-minute drain window. A sensible default for mixed interactive and batch clusters. Start from the **Balanced Autoscaling** preset.

**Aggressive batch** -- Full-strength factors on a spot-heavy weight mix for bursty batch pipelines where reaction speed wins. Start from the **Aggressive Batch** preset.

**Conservative production** -- Small factors and long drain windows for SLA-bound shared clusters where oscillation hurts interactive users. Start from the **Conservative Production** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the policy is created
- [**GCP Dataproc Cluster**](/cloud-catalog/gcp-dataproc-cluster) -- attaches this policy via `autoscalingPolicyUri`; one policy can govern many clusters
