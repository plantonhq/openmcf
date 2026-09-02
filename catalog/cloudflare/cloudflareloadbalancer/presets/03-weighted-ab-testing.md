---
display_name: Weighted A/B Testing
---

# Weighted A/B Testing

A monitor, two pools (control and variant), and a load balancer with
`steeringPolicy: random` plus `randomSteering` weights. Cloudflare selects a pool
at random in proportion to the configured weights — ideal for canary releases and
A/B experiments.

## When to Use

- Canary rollouts: send a small share of traffic to a new version
- A/B experiments across two backends

## Key Configuration Choices

- **`steeringPolicy: random`** — weighted random pool selection.
- **`randomSteering.defaultWeight`** — base weight applied to pools not listed in
  `poolWeights`; set per-pool weights via `randomSteering.poolWeights` (keyed by
  pool ID) for finer control.
- **`defaultPools`** lists both pools; `fallbackPool` is the safe control pool.

## Placeholders to Replace

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | Account that owns the monitor and pools |
| `<cloudflare-zone-id>` | Zone containing the hostname |
| `<app-subdomain>.replaceme.example.com` | Load balancer hostname |
| `<control-origin-ip-or-hostname>` / `<variant-origin-ip-or-hostname>` | Origin addresses |

## Related Presets

- **01-active-passive-failover** — simple primary/backup failover
- **02-geographic-routing** — route by geography across regional pools
