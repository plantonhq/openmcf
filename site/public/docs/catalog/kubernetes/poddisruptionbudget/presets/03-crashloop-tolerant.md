---
title: "Crashloop Tolerant"
description: "This preset is the availability floor that cannot wedge a node drain. Under the default eviction policy, a budget counts only READY pods as available — so a crash-looping application (running, never..."
type: "preset"
rank: "03"
presetSlug: "03-crashloop-tolerant"
componentSlug: "poddisruptionbudget"
componentTitle: "PodDisruptionBudget"
provider: "kubernetes"
icon: "package"
order: 3
---

# Crashloop Tolerant

This preset is the availability floor that cannot wedge a node drain. Under the default eviction policy, a budget counts only READY pods as available — so a crash-looping application (running, never ready) holds its budget below the floor forever, and every drain touching its pods blocks indefinitely. `always_allow` lets running-but-not-ready pods be evicted regardless of the budget, keeping cluster maintenance moving while healthy pods stay protected.

## When to Use

- Budgets over workloads that can crash-loop — anything with flaky startup, external dependencies at boot, or a history of bad deploys
- Clusters where node drains are automated (managed upgrades, autoscaler consolidation) and a stuck drain pages a human
- As the default shape for operator-managed workloads whose readiness you do not control

## Key Configuration Choices

- **`unhealthy_pod_eviction_policy: always_allow`** — not-yet-ready pods may ALWAYS be evicted; only healthy pods count against (and are protected by) the floor. The trade-off: a pod that was about to become ready may be restarted during a drain
- **`min_available: "1"`** — the floor still protects healthy pods: a drain cannot take the last available pod
- **`selector.match_labels.app`** — targets one workload's pods via the workload label contract (`app: <workload-metadata-name>`)
- **Engine boundary — this preset deploys via the Pulumi provisioner only.** The Terraform kubernetes provider cannot express `unhealthyPodEvictionPolicy`, and the Terraform module rejects `always_allow` with a plan-time precondition by design — failing loudly beats silently deploying the default (`if_healthy_budget`), which is the opposite of this preset's intent

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace the protected pods run in — a budget governs only its own namespace | Your namespace management |
| `<your-workload-name>` | The target workload's `metadata.name` (its pods carry `app: <name>`) | The workload's manifest |

## Related Presets

- **01-protect-workload** — the same floor with the default (conservative) unhealthy-pod policy
- **02-tier-percentage** — cross-workload coverage with a percentage ceiling
