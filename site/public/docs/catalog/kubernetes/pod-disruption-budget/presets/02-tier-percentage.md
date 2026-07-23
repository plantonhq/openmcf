---
title: "Tier Percentage"
description: "This preset spans several workloads with one budget: every pod labelled `tier: web` or `tier: api` is covered, and at most a quarter of them may be down at once during voluntary disruptions. It is..."
type: "preset"
rank: "02"
presetSlug: "02-tier-percentage"
componentSlug: "pod-disruption-budget"
componentTitle: "Pod Disruption Budget"
provider: "kubernetes"
icon: "package"
order: 2
---

# Tier Percentage

This preset spans several workloads with one budget: every pod labelled `tier: web` or `tier: api` is covered, and at most a quarter of them may be down at once during voluntary disruptions. It is the shape for tier-level protection — one budget guarding a whole class of services instead of one budget per workload.

## When to Use

- Guarding a tier of related services (all web frontends, all API services) with one policy instead of many
- Workloads that scale up and down — a percentage ceiling tracks the replica count where an absolute floor goes stale
- Cluster-upgrade readiness for namespaces with many small services sharing a labelling convention

## Key Configuration Choices

- **`match_expressions` with `In`** — set-based selection covers what exact-match cannot: any pod whose `tier` label is one of the listed values. This requires the workloads to share the labelling convention (`tier` is a user-defined label, not one Planton stamps automatically)
- **`max_unavailable: "25%"`** — the ceiling form: at most 25% of the selected pods (rounded UP, giving evictions more room) may be down at once. Prefer the ceiling for anything that scales; it adapts as replicas change
- **Cross-workload coverage is the point — and the caveat**: every pod matched here must NOT be covered by any other budget (including a Planton workload's built-in `availability.pod_disruption_budget` block). A pod under two budgets makes evictions fail and wedges drains
- **The percentage is computed per owning controller** — the disruption controller resolves desired replica counts through each pod's controller, so this composes with Deployments/StatefulSets/ReplicaSets

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the tier's pods run in — a budget governs only its own namespace | Your namespace management |

The `tier` values `web` and `api` are working examples — replace them with your own tier labelling convention.

## Related Presets

- **01-protect-workload** — the single-workload floor, selected by the automatic `app` label
- **03-crashloop-tolerant** — add drain protection against crash-looping pods
