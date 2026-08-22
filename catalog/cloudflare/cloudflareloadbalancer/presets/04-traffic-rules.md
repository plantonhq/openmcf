# Traffic Rules

One load balancer, two pools, and two rules: `/api` requests steer to a
dedicated API pool with session affinity switched off, and `/maintenance`
answers directly at the edge with a 503 page -- no origin involved. Everything
else uses the web pool with cookie affinity.

## When to Use

- Per-path routing (API vs web) on a single hostname
- Disabling session affinity for stateless endpoints only
- Serving a maintenance or fallback response at the edge for specific paths

## Key Configuration Choices

- **Rule order** — rules run in `priority` order (lower first); when no rule
  sets one, list order decides. This preset relies on list order.
- **`overrides.sessionAffinity: none`** — an EXPLICIT `none` switches affinity
  off for matched traffic; leaving it unset would inherit the top-level
  `cookie` mode.
- **`fixedResponse`** — a fixed-response rule is always terminating; no
  origin is contacted for matched requests.
- **Condition expressions** — see Cloudflare's load-balancing rules
  expression language for what a `condition` can match.

## Placeholders to Replace

| Placeholder | Description |
|---|---|
| `<cloudflare-account-id>` | Account that owns the monitor and pools |
| `<cloudflare-zone-id>` | Zone containing the hostname |
| `<app-subdomain>.<your-domain.com>` | Load balancer hostname |
| `192.0.2.1`, `192.0.2.2` | Web and API origin addresses |

## Related Presets

- **01-active-passive-failover** — static failover without rules
- **02-geographic-routing** — route by geography across regional pools
- **03-weighted-ab-testing** — split traffic across pools by weight
