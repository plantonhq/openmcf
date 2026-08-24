# Preemptable Batch

This preset creates the bottom step of the ladder: a negative-value, non-preempting class for work that should run only on spare capacity. Pods of this class yield to everything — including unmarked pods at priority 0 — and, because preemption is disabled, they never evict anything already running even when they queue ahead of other pending work. This is the honest expression of "opportunistic": take idle capacity, give it back under pressure, disturb nothing.

## When to Use

- Batch jobs, nightly reports, data pipelines, and ML training runs that tolerate eviction and restarts
- Backfill and reprocessing work that should soak up idle capacity without ever competing with services
- Any workload whose deadline is soft and whose eviction is cheap

## Key Configuration Choices

- **`value: -100`** — negative values are valid and idiomatic for always-preemptable tiers: this class sits below even priority-0 unmarked pods, so under pressure the scheduler evicts batch pods FIRST to make room for anything else
- **`preemption_policy: never`** — the class's own pods jump the queue ahead of lower-priority pending pods but NEVER evict running pods to schedule. Note the policy governs what this class's pods do to others; it does not protect them — they remain fully preemptable BY higher tiers, which the negative value guarantees
- **Deliberately NOT a global default** — batch status must be opted into per workload; inheriting it silently would demote ordinary services to evict-first
- **Workloads must tolerate eviction** — preemption is a graceful delete (termination grace periods apply); pair this class with Jobs that handle retries, checkpointing for long training runs, and idempotent processing

## Placeholders to Replace

None — this preset is deployable as-is. Adjust `value` only if your ladder uses different spacing.

## Related Presets

- **01-critical-services** — the top step whose pods may evict this tier
- **02-standard-default** — the default step; its value 1000 is what makes this tier's negative value meaningful
