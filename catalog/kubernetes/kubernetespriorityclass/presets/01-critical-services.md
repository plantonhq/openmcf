# Critical Services

This preset creates the top rung of the user-definable importance ladder: value 1,000,000 with preemption enabled. Pods that reference this class (via the workload pod spec's `priority_class_name`) schedule ahead of everything below it — and when the cluster has no room, the scheduler EVICTS lower-priority pods to make room for them. This is the class for workloads that must run even if something else has to move.

## When to Use

- Revenue-path services and user-facing APIs whose downtime is an incident
- Core platform services (ingress controllers, cert managers) that everything else depends on
- Any workload for which "stays Pending during capacity pressure" is unacceptable

## Key Configuration Choices

- **`value: 1000000`** — high in the user range but well below the ceiling (user classes max out at 1,000,000,000; above that is reserved for Kubernetes' own `system-*` classes). The million-unit headroom below the ceiling leaves room for future higher tiers without renumbering — the value is IMMUTABLE, and changing it replaces the class
- **Preemption left at its default (`preempt_lower_priority`)** — this is the point of a critical tier: under pressure, lower-priority pods are evicted (gracefully, respecting termination grace periods and, best-effort, PodDisruptionBudgets) to make room
- **`description` set** — surfaced by `kubectl describe priorityclass`; it is the only in-cluster documentation of when to use the class. Keep it honest, or everything becomes "critical"
- **NOT the global default** — criticality must be opted into per workload, never inherited by unmarked pods

> Reserve this class ruthlessly. If most workloads are critical, none are — and every admission to this tier licenses eviction of everything below it.

## Placeholders to Replace

None — this preset is deployable as-is. Adjust `value` only if your ladder uses different spacing.

## Related Presets

- **02-standard-default** — the tier unmarked pods should land in; deploy alongside this one
- **03-preemptable-batch** — the bottom rung this class's pods may evict
