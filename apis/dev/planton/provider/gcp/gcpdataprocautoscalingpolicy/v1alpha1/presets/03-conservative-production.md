# Conservative Production

A smooth, stability-first scaling policy for SLA-bound clusters: small
scale factors, minimum-change fractions that filter out scaling noise, a
long cooldown, and a 2-hour graceful decommission window.

## When to use

- Long-running production clusters serving steady workloads
- Spark Structured Streaming or long-task jobs where a removed worker
  is expensive
- Environments where predictable capacity matters more than fast
  reaction

## What to customize

- `projectId` / `location` — the policy must live in the same region as
  every cluster that attaches it.
- `workerConfig.minInstances` — the guaranteed on-demand floor; size to
  the cluster's steady-state load.
- `gracefulDecommissionTimeout` — match your longest routine task (the
  API allows up to 1 day, `"86400s"`).
- The min-worker fractions — raise to `0.2` for even fewer scaling
  events, or drop to `0.0` to act on any recommendation.

## Key configuration

- **Small factors (0.2 up / 0.1 down)** — the autoscaler moves in small
  steps, and shrinks even more cautiously than it grows
- **Min-worker fractions of 0.1** — recommendations that would change
  the cluster by less than 10% are ignored entirely, eliminating
  scaling noise
- **5-minute cooldown** — longer settle time between evaluations
- **Primary-weighted (2:1)** — most new capacity is on-demand;
  the spot arm stays a modest supplement
- **2-hour graceful decommission** — long-running tasks finish before
  their worker disappears

## How clusters attach it

```yaml
# In a GcpDataprocCluster spec:
clusterConfig:
  autoscalingPolicyUri:
    valueFrom:
      kind: GcpDataprocAutoscalingPolicy
      name: conservative-production
      fieldPath: status.outputs.name
```

Tuning this policy later re-tunes every attached cluster in place — no
cluster recreation.

## Related presets

- **01-balanced-autoscaling** — moderate factors for mixed workloads
- **02-aggressive-batch** — maximum-speed scaling on a spot-heavy mix
