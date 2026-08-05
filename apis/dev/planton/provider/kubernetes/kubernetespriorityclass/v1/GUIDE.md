# KubernetesPriorityClass Guide

The judgment this guide carries: a PriorityClass is not just a number —
above a threshold it lets the scheduler EVICT running pods to make room,
so the cluster's priority ladder is a deliberate, small design, not a
per-workload afterthought.

## Design the ladder once, cluster-wide

Priority is only meaningful as a comparison, so decide the whole ladder
deliberately and keep it small: e.g. "critical" (high value, preempting)
for revenue-path services, "standard" (the global default) for everything
unmarked, "batch" (negative, non-preempting) for work that must never
displace services. Workloads then reference a class by name through the
shared pod spec's `priorityClassName`. Scattering ad-hoc numbers across
teams defeats the point — the ordering only holds if one ladder governs
the cluster.

## Preemption is the sharp edge

Above the preemption threshold, when a high-priority pod cannot schedule,
the scheduler EVICTS lower-priority pods to fit it (the reference page has
the exact semantics). That is the feature for revenue-path work and the
foot-gun for everything else: mark a chatty batch job critical and it
starts evicting the services it was meant to stay clear of. Non-preempting
classes get the ordering benefit without the eviction power — the right
choice for most tiers above default. The built-in
`system-cluster-critical` / `system-node-critical` classes are
Kubernetes's own and sit above the user range.

## On the diagram

PriorityClasses are cluster-scoped policy, not composition nodes —
workloads carry a class NAME, not a drawn edge. The architecture-level
fact worth stating in a proposal is which tier a workload sits in, since
it decides who survives contention.

## Pairs well with

- Every workload kind — via the pod spec's `priorityClassName`.
- KubernetesResourceQuota — quotas can bound consumption per priority
  class, pairing importance with entitlement.
