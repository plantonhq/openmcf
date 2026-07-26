---
title: "Presets"
description: "Ready-to-deploy configuration presets for ClickHouse"
type: "preset-list"
componentSlug: "clickhouse"
componentTitle: "ClickHouse"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-minimal"
    rank: "01"
    title: "Dev minimal preset"
    excerpt: "The smallest declarable ClickHouse that actually serves: one host, a PVC, a pinned server version and one named user. No Keeper is deployed because a single-replica topology needs no coordination —..."
  - slug: "02-production-replicated"
    rank: "02"
    title: "Production replicated preset"
    excerpt: "The durability posture: one shard carried by three replicas in ReplicatedMergeTree lockstep, coordinated by a three-node managed Keeper (survives one Keeper loss), with replicas forced onto different..."
  - slug: "03-sharded-analytics"
    rank: "03"
    title: "Sharded analytics preset"
    excerpt: "The capacity posture: four shards, each a disjoint slice of the data carried by two replicas — eight ClickHouse hosts a Distributed table queries in parallel. This is the shape for datasets or write..."
---

# ClickHouse Presets

Ready-to-deploy configuration presets for ClickHouse. Each preset is a complete manifest you can copy, customize, and deploy.
