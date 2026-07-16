---
title: "Dedicated Serving Compute"
description: "Deploys an index with a pinned machine type and explicit replica bounds — predictable performance and cost for production serving."
type: "preset"
rank: "02"
presetSlug: "02-dedicated"
componentSlug: "vertex-ai-deployed-index"
componentTitle: "Vertex AI Deployed Index"
provider: "gcp"
icon: "package"
order: 2
---

# Dedicated Serving Compute

Deploys an index with a pinned machine type and explicit replica bounds —
predictable performance and cost for production serving.

## What this preset creates

A DeployedIndex on the referenced endpoint, served by 2–10
`e2-standard-16` replicas. The machine type is pinned, so query latency
and per-replica cost are known quantities instead of Vertex-managed
variables.

## When to use

- Production serving with latency SLOs that need a known machine profile
- Cost planning that requires a fixed per-replica price
- Corpora whose shard size dictates a specific machine class

## Remix ideas

- Match `machineType` to the index's shard size: MEDIUM shards (the
  default) take `e2-standard-16` and up; LARGE shards need
  `e2-highmem-16` or `n2d-standard-32`; SMALL shards allow
  `e2-standard-2`.
- Size `minReplicaCount` for your P99 floor — dedicated replicas are
  always on and always billed. The replica bounds are the only in-place
  knobs; `machineType` itself is immutable and changing it redeploys.
- Decide the machine type together with the index's `shardSize` at
  index-creation time; the two are chosen together, not independently.
