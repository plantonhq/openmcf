# Balanced Autoscaling

A moderate, general-purpose autoscaling policy: half-strength scale
factors, an even primary/secondary split, and a 30-minute graceful
decommission window. A sensible default for mixed interactive and batch
clusters.

## When to use

- A first autoscaling policy for teams new to Dataproc autoscaling
- Clusters serving a mix of interactive queries and scheduled jobs
- Workloads where neither reaction speed nor cost dominates

## What to customize

- `projectId` / `location` — the policy must live in the same region as
  every cluster that attaches it.
- `workerConfig.maxInstances` / `secondaryWorkerConfig.maxInstances` —
  the hard ceilings; size them to your project quota.
- `scaleUpFactor` / `scaleDownFactor` — raise toward `1.0` for faster
  reaction, lower toward `0.1` for smoother scaling.

## Key configuration

- **Factors of 0.5** — the autoscaler acts on half of what YARN memory
  metrics suggest per evaluation, damping oscillation
- **Equal weights (1:1)** — new capacity splits evenly between
  on-demand primaries and spot secondaries
- **Primary floor of 2** — the Dataproc API's minimum for autoscaled
  clusters; secondaries scale to zero when idle
- **30-minute graceful decommission** — running tasks finish before
  their worker is removed
- **2-minute cooldown** — the default evaluation cadence

## How clusters attach it

```yaml
# In a GcpDataprocCluster spec:
clusterConfig:
  autoscalingPolicyUri:
    valueFrom:
      kind: GcpDataprocAutoscalingPolicy
      name: balanced-autoscaling
      fieldPath: status.outputs.name
```

One policy can govern many clusters; editing the policy re-tunes all of
them in place.

## Related presets

- **02-aggressive-batch** — maximum-speed scaling on a spot-heavy mix
- **03-conservative-production** — small factors and long drain windows for SLA-bound clusters
