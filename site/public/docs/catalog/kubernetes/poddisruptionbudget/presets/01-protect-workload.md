---
title: "Protect Workload"
description: "This preset is the standard availability floor for one workload: at least one of its pods must survive any voluntary disruption. Node drains, cluster upgrades, and autoscaler consolidation all go..."
type: "preset"
rank: "01"
presetSlug: "01-protect-workload"
componentSlug: "poddisruptionbudget"
componentTitle: "PodDisruptionBudget"
provider: "kubernetes"
icon: "package"
order: 1
---

# Protect Workload

This preset is the standard availability floor for one workload: at least one of its pods must survive any voluntary disruption. Node drains, cluster upgrades, and autoscaler consolidation all go through the eviction API, and the budget makes that API refuse any step that would leave zero pods available — the drain waits and retries as replacements become ready.

## When to Use

- Protecting operator-managed pods or non-Planton workloads that no Planton workload kind's built-in budget covers
- Any multi-replica service whose availability during node maintenance matters
- As the first budget in a namespace being hardened for cluster upgrades

## Key Configuration Choices

- **`selector.match_labels.app`** — targets one workload's pods via the workload label contract: every Planton workload kind stamps `app: <workload-metadata-name>` on its pods as immutable selection identity
- **`min_available: "1"`** — the floor form: at least one selected pod stays available through any drain. With 2+ replicas, drains proceed one pod at a time; with a single replica, every drain touching that pod is blocked (a budget can only refuse evictions, not create availability)
- **Voluntary disruptions only** — node crashes and OOM kills never consult the budget; replica count is still the availability mechanism for involuntary failures
- **One budget per set of pods** — a pod covered by more than one budget makes evictions FAIL. For a Planton Deployment's or StatefulSet's OWN pods, use the workload's built-in `availability.pod_disruption_budget` block instead of this kind — never both

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the protected pods run in — a budget governs only its own namespace | Your namespace management |
| `<your-workload-name>` | The target workload's `metadata.name` (its pods carry `app: <name>`) | The workload's manifest |

## Related Presets

- **02-tier-percentage** — one budget spanning several workloads with a percentage ceiling
- **03-crashloop-tolerant** — the same floor plus protection against crash-looping pods wedging drains
