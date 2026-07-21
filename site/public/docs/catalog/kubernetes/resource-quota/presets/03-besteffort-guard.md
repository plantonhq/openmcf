---
title: "BestEffort Guard"
description: "This preset caps only the namespace's BestEffort pods — pods with no requests or limits at all — at 10. It is a scoped quota: the `best_effort` scope makes the quota track exclusively the naive pods,..."
type: "preset"
rank: "03"
presetSlug: "03-besteffort-guard"
componentSlug: "resource-quota"
componentTitle: "Resource Quota"
provider: "kubernetes"
icon: "package"
order: 3
---

# BestEffort Guard

This preset caps only the namespace's BestEffort pods — pods with no requests or limits at all — at 10. It is a scoped quota: the `best_effort` scope makes the quota track exclusively the naive pods, leaving every well-behaved (Burstable or Guaranteed) workload completely untouched. BestEffort pods are the cluster's least accountable tenants: the scheduler reserves nothing for them and they can consume whatever a node has, so containing their count is a targeted guard against unbounded `kubectl run`-style sprawl.

## When to Use

- Shared namespaces where declared workloads should be unconstrained but ad-hoc naive pods must be contained
- As a lightweight alternative to compute caps: no LimitRange needed, no risk of rejecting declared workloads
- Development namespaces where experimentation is allowed but must stay bounded

## Key Configuration Choices

- **`scopes: [best_effort]`** — the quota counts a pod only if it is BestEffort QoS (no requests or limits on any container). Everything else in the namespace is invisible to this quota
- **`hard` caps only `pods`** — a `best_effort`-scoped quota may cap ONLY pod counts; BestEffort pods have no requests or limits to meter, so compute entries are invalid here (the schema rejects the combination, mirroring the API rule)
- **No `limit_defaults`** — deliberately: adding container defaults to the namespace would convert naive pods to Burstable QoS, and this quota would then track nothing. Use this preset when you want naive pods allowed-but-bounded rather than defaulted
- **Composable** — this quota coexists with an unscoped compute or count quota in the same namespace; each tracks its own slice

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The namespace whose BestEffort pods to cap | Your namespace management |

## Related Presets

- **01-team-namespace-governed** — the opposite philosophy: default naive pods into declared values instead of capping them
- **02-object-count-caps** — unscoped counts for all pods and other objects
