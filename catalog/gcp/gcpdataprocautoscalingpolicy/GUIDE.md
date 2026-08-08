# GcpDataprocAutoscalingPolicy Guide

The judgment this guide protects: the policy is fleet-wide tuning in a
single resource — one policy governs many clusters, so a change here
re-tunes every attached cluster in place. Treat edits like config
rollouts, not like editing one cluster.

## Weights steer money, factors steer speed

The worker/secondary `weight` pair decides WHERE new capacity lands —
weight 1 primary against weight 3 secondary sends ~75% of growth to spot
nodes, which is the cost posture most batch estates want. The YARN
`scaleUpFactor` / `scaleDownFactor` decide HOW HARD the autoscaler acts
on memory pressure: 1.0 is maximally aggressive, small factors smooth
scaling at the cost of reaction time. Tune weights for the bill, factors
for the workload's tolerance of node churn.

## Scale-down is where jobs die

`gracefulDecommissionTimeout` in the YARN config is the window running
tasks get before a node is removed — too short kills long tasks on every
scale-in, too long makes the autoscaler useless. Pair gentle scale-down
(`scaleDownFactor` well below 1.0) with the CLUSTER's `idleDeleteTtl` as
the cost backstop: the policy keeps capacity responsive while jobs run,
the TTL reclaims everything when they stop. Primary workers carry HDFS —
the API's floor of 2 exists because scaling them to zero would take the
filesystem with it.

## Destroy stance

The API refuses to delete a policy while any cluster references it —
teardown order is clusters first, policy last. `deletionPolicy: ABANDON`
releases the policy from IaC management (the referencing clusters keep
using it); `PREVENT` guards the shared policy an entire region's
clusters depend on — which is exactly the resource whose accidental
destroy hurts the most.

## On the diagram

A region-scoped leaf under `GcpProject` that many `GcpDataprocCluster`
nodes attach to via `autoscalingPolicyUri`. One policy with many
attached clusters is the intended shape — per-cluster copies of the same
bounds are tuning drift waiting to happen.

## Pairs well with

- `GcpDataprocCluster` — attaches the policy; its
  `gracefulDecommissionTimeout` and TTLs complete the scaling story.
- `GcpProject` — the scope; remember policies are per-region.
