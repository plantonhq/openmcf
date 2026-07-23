---
title: "Spot Diversified"
description: "This preset declares a cost-optimized spot pool with deliberate instance-family diversity: a wide `In` list of families plus `minValues` forcing the launched fleet to span at least four of them, so a..."
type: "preset"
rank: "02"
presetSlug: "02-spot-diversified"
componentSlug: "karpenter-node-pool"
componentTitle: "Karpenter Node Pool"
provider: "kubernetes"
icon: "package"
order: 2
---

# Spot Diversified

This preset declares a cost-optimized spot pool with deliberate
instance-family diversity: a wide `In` list of families plus `minValues`
forcing the launched fleet to span at least four of them, so a reclaim of
one spot pool cannot take out most of the fleet at once. Disruption
budgets bound how much of the pool churns concurrently and freeze
voluntary disruption during business hours.

## When to Use

- Interruption-tolerant workloads (stateless services behind load
  balancers, batch/queue consumers) chasing spot pricing
- As a companion to **01-general-purpose-on-demand**: pods that tolerate
  interruptions land here; the rest land on the on-demand pool

## Key Configuration Choices

- **Spot only** (`karpenter.sh/capacity-type In [spot]`) — the pool's
  entire posture; pair with the Karpenter installation's interruption
  queue so nodes drain ahead of reclaims
- **`instance-family In [...8 families...]` with `minValues: 4`** —
  `minValues` (ALPHA) is the diversity knob for spot pools: the fleet
  must span at least 4 distinct families, and the requirement must list
  at least that many values
- **Budgets** — an always-active `15%` ceiling on concurrent disruption,
  plus a `0`-node budget active weekdays 09:00 for 8 hours that blocks
  voluntary disruption during business hours; a budget's `schedule` and
  `duration` must be set together, and when several budgets are active
  the most restrictive value applies. Spot reclaims are involuntary and
  not governed by budgets
- **`consolidationPolicy: WhenEmptyOrUnderutilized`** (CRD default) with
  a 1-minute settle window
- **`limits.cpu: "500"`** — the spot pool's cost ceiling

## Placeholders to Replace

| Placeholder             | Description                            | Where to Find                                             |
| ----------------------- | -------------------------------------- | --------------------------------------------------------- |
| `<ec2-node-class-name>` | Name of the EC2NodeClass to build from | `metadata.name` of your `KubernetesKarpenterEc2NodeClass` |

## Related Presets

- **01-general-purpose-on-demand** — the on-demand baseline pool
- **03-gpu-dedicated** — tainted GPU pool for accelerated workloads
