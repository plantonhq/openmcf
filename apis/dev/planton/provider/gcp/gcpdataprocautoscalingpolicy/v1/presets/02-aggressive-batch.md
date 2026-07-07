# Aggressive Batch

A maximum-speed scaling policy for retryable batch workloads: full
scale factors in both directions, a spot-heavy 4:1 weight split, a
scale-to-zero secondary group, and a short drain window.

## When to use

- Nightly or scheduled ETL where jobs arrive in bursts
- Retryable Spark batch jobs that tolerate preemption and fast
  scale-down
- Cost-sensitive pipelines that should hold zero burst capacity between
  runs

## What to customize

- `projectId` / `location` — the policy must live in the same region as
  every cluster that attaches it.
- `secondaryWorkerConfig.maxInstances` — the burst ceiling; size to the
  largest job's parallelism and your spot quota.
- `gracefulDecommissionTimeout` — lengthen if tasks regularly run
  longer than 5 minutes and retries are expensive.

## Key configuration

- **Factors of 1.0** — the autoscaler adds and removes everything the
  YARN memory metrics suggest, every evaluation; fastest possible
  reaction at the cost of some churn
- **Weight 4 on secondaries** — ~80% of new capacity lands on the spot
  group; the on-demand base stays small and stable
- **Scale-to-zero secondaries** — `minInstances: 0` means the burst arm
  costs nothing between job runs
- **Small primary group (2-4)** — just enough on-demand capacity to
  keep HDFS healthy
- **5-minute decommission** — quick scale-down; Spark task retry
  absorbs any interrupted work

## How clusters attach it

Pair this policy with a cluster whose secondary workers are SPOT (see
the `GcpDataprocCluster` preset `03-cost-optimized-batch`, which
references a policy exactly like this one by `valueFrom`).

## Related presets

- **01-balanced-autoscaling** — moderate factors for mixed workloads
- **03-conservative-production** — small factors and long drain windows for SLA-bound clusters
