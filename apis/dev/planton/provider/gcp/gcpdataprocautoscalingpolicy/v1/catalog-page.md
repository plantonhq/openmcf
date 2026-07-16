# GCP Dataproc Autoscaling Policy

Creates a Dataproc autoscaling policy — the reusable resource that governs how Dataproc clusters scale their worker groups on YARN memory pressure. One policy can govern many clusters: each attaches it by reference, so scaling behavior is tuned in one place and every attached cluster follows.

## What Gets Created

A regional `google_dataproc_autoscaling_policy` resource holding worker-group bounds, group weights, and the YARN-based scaling algorithm. Clusters attach it through their `clusterConfig.autoscalingPolicyUri` field — an in-place update on the cluster.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A GCP project** (the Dataproc API is enabled automatically)
- Clusters that will attach the policy must live in the **same region**

## Quick Start

Create a file `policy.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpDataprocAutoscalingPolicy
metadata:
  name: etl-autoscaling
spec:
  projectId:
    value: my-gcp-project
  policyId: etl-autoscaling
  location: us-central1
  workerConfig:
    maxInstances: 10
    minInstances: 2
  secondaryWorkerConfig:
    maxInstances: 20
    weight: 3
  basicAlgorithm:
    cooldownPeriod: "120s"
    yarnConfig:
      gracefulDecommissionTimeout: "3600s"
      scaleUpFactor: 0.5
      scaleDownFactor: 1.0
```

Deploy:

```shell
planton apply -f policy.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `policyId` | `string` | The GCP resource name of the policy. Immutable. | 3-50 chars; letters, digits, underscores, hyphens; starts/ends alphanumeric |
| `location` | `string` | Region the policy lives in (e.g., `us-central1`). Clusters attach same-region policies only. Immutable. | Required, `^[a-z]+-[a-z]+[0-9]+$` |
| `workerConfig.maxInstances` | `int` | Ceiling on primary workers. | Required, >= 1 |
| `basicAlgorithm.yarnConfig.gracefulDecommissionTimeout` | `string` | Drain window before a worker is forcefully removed on scale-down (e.g., `"3600s"`). Bounds: 0s to 1 day. | Required, `^[0-9]+s$` |
| `basicAlgorithm.yarnConfig.scaleUpFactor` | `double` | Fraction of pending YARN memory acted on per evaluation. `1.0` scales up as fast as possible. | Required, 0.0-1.0 |
| `basicAlgorithm.yarnConfig.scaleDownFactor` | `double` | Fraction of idle YARN memory removed per evaluation. | Required, 0.0-1.0 |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | GCP project. Can reference GcpProject via `valueFrom`. |
| `workerConfig.minInstances` | `int` | `2` (API default) | Primary floor. `0` accepts the API default; otherwise must be >= 2 (the Dataproc minimum for autoscaled clusters). |
| `workerConfig.weight` | `int` | `1` | Relative share of new capacity landing on primaries. |
| `secondaryWorkerConfig` | `object` | not scaled | Bounds for the secondary (spot) group. |
| `secondaryWorkerConfig.maxInstances` | `int` | `0` | Secondary ceiling. |
| `secondaryWorkerConfig.minInstances` | `int` | `0` | Secondary floor — `0` lets the group scale to zero when idle. |
| `secondaryWorkerConfig.weight` | `int` | `1` | Relative share of new capacity landing on secondaries (e.g., primary 1 + secondary 3 sends ~75% of new nodes to spot). |
| `basicAlgorithm.cooldownPeriod` | `string` | `"120s"` | Time between scaling evaluations. Bounds: 2 minutes to 1 day. |
| `basicAlgorithm.yarnConfig.scaleUpMinWorkerFraction` | `double` | `0.0` | Minimum fractional cluster change worth acting on when scaling up (e.g., `0.05` ignores sub-5% recommendations). |
| `basicAlgorithm.yarnConfig.scaleDownMinWorkerFraction` | `double` | `0.0` | Same filter for scale-down. |

## Example: Attaching from a Cluster

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpDataprocCluster
metadata:
  name: batch-spark
spec:
  region: us-central1
  clusterName: batch-spark
  clusterConfig:
    secondaryWorkerConfig:
      preemptibility: SPOT
    autoscalingPolicyUri:
      valueFrom:
        kind: GcpDataprocAutoscalingPolicy
        name: etl-autoscaling
        fieldPath: status.outputs.name
```

The reference resolves to the policy's full resource name. Attaching, swapping, or detaching the policy updates the cluster in place.

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `name` | `string` | Fully qualified policy resource name (`projects/{project}/locations/{location}/autoscalingPolicies/{policy_id}`) — the handle a cluster's `autoscalingPolicyUri` reference resolves to |
| `policy_id` | `string` | The policy ID (same as `spec.policyId`) |
| `location` | `string` | Region the policy lives in |

## Related Components

- [GcpDataprocCluster](/docs/catalog/gcp/gcpdataproccluster) — attaches this policy via `clusterConfig.autoscalingPolicyUri`
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project
