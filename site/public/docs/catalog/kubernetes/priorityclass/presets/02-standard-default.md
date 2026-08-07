---
title: "Standard Default"
description: "This preset creates the middle rung of the ladder AND makes it the cluster-wide default: every pod that names no priority class receives value 1000 instead of the bare-Kubernetes default of 0. This..."
type: "preset"
rank: "02"
presetSlug: "02-standard-default"
componentSlug: "priorityclass"
componentTitle: "PriorityClass"
provider: "kubernetes"
icon: "package"
order: 2
---

# Standard Default

This preset creates the middle rung of the ladder AND makes it the cluster-wide default: every pod that names no priority class receives value 1000 instead of the bare-Kubernetes default of 0. This lifts all ordinary workloads above the batch tier (which sits at a negative value) without anyone having to label anything.

## When to Use

- Every cluster that adopts a priority ladder — unmarked pods should land on a deliberate rung, not at 0
- Making the critical/standard/batch distinction meaningful: with the default at 1000, batch at -100 genuinely yields to ordinary workloads

## Key Configuration Choices

- **`global_default: true`** — pods that set no `priority_class_name` get this class at admission. **Exactly ONE class per cluster may sensibly carry this flag.** Kubernetes does NOT reject multiple global defaults — when several claim it, it silently uses the SMALLEST such value, which is rarely what anyone intended. Audit `kubectl get priorityclass` for an existing `GLOBAL-DEFAULT: true` entry before deploying this preset
- **Changing the default never re-prioritizes existing pods** — it applies only to pods created afterwards; expect the ladder to converge as workloads roll
- **`value: 1000`** — deliberately modest: well below the critical tier (1000000) so critical pods preempt standard ones under pressure, and well above 0 and the batch tier so ordinary workloads outrank opportunistic work. The 1000-unit spacing leaves room for intermediate tiers without renumbering (the value is immutable; renumbering is a replace)
- **Preemption left at its default (`preempt_lower_priority`)** — standard workloads may displace the batch tier when capacity is scarce, which is the intended economics

## Placeholders to Replace

None — this preset is deployable as-is, after confirming no other class in the cluster carries `global_default: true`.

## Related Presets

- **01-critical-services** — the tier above; opted into explicitly per workload
- **03-preemptable-batch** — the tier below; the negative value only means "yields to everything" once this default exists
