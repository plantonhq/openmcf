# Behavior Tuned

This preset keeps the fast scale-UP defaults and makes scale-DOWN deliberately slow: wait 10 minutes of consistently lower recommendations, then remove at most 10% of the fleet per minute. The Kubernetes defaults scale down to the recommendation's 300-second high-water mark and can then remove ALL surplus pods in one step — for spiky traffic that means cliff-dropping capacity right before the next spike, paying cold-start latency on every rebound.

## When to Use

- Spiky or periodic traffic where the default 5-minute scale-down window is shorter than the gap between spikes
- Workloads with expensive warm-up (JVM services, cache-heavy apps) where removed capacity is slow to rebuild
- Queue-driven fleets that should bleed down gradually after the backlog clears rather than dropping to the floor at once

## Key Configuration Choices

- **`scale_down.stabilization_window_seconds: 600`** — the flap damper: the controller scales down only to the SAFEST (highest) recommendation of the past 10 minutes. Traffic must stay low for the full window before capacity leaves
- **`percent 10 per 60s` policy** — once the window agrees, at most 10% of current replicas are removed per minute — a gradual bleed instead of a cliff. With several policies listed, `select_policy` arbitrates (`max_change` is the default; `disabled` turns a direction off entirely — the incident lever for freezing scale-down while leaving scale-up live)
- **`scale_up` untouched** — the Kubernetes default (the higher of doubling or +4 pods per 15s, no stabilization) stays: under-provisioning during a surge is usually worse than over-provisioning after one
- **`resource` cpu at 60%** — the workhorse signal; the behavior block composes identically with any metric family, including the queue-driven external form

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | The scale target's own namespace — an HPA cannot scale across namespaces | Your namespace management |
| `<your-workload-name>` | The target Deployment's name | The workload's manifest |

The window (600s) and bleed rate (10%/min) are working defaults — size the window to the gap between your traffic spikes.

## Related Presets

- **01-cpu-autoscale** — the same signal with default behavior
- **03-queue-driven** — the external-metric fleet this scale-down discipline pairs naturally with
