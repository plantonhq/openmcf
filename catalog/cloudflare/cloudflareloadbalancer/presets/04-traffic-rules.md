# Traffic Rules

One load balancer, two pools, and two rules: `/api` requests steer to a
dedicated API pool with session affinity switched off, and `/maintenance`
answers directly at the edge with a 503 page -- no origin involved. Everything
else uses the web pool with cookie affinity.

> **Plan requirement**: rule count is a subscription-tier limit. The
> entry-level (Basic) Load Balancing subscription allows exactly ONE custom
> rule per load balancer, so this two-rule preset is rejected there with
> `rule count 2 exceeds limit 1` (400, code 1002). Use a plan with a higher
> rule limit, or keep just one of the two rules on Basic.

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
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | Account that owns the monitor and pools |
| `<cloudflare-zone-id>` | Zone containing the hostname |
| `<app-subdomain>.replaceme.example.com` | Load balancer hostname |
| `192.0.2.1`, `192.0.2.2` | Web and API origin addresses |

## Related Presets

- **01-active-passive-failover** — static failover without rules
- **02-geographic-routing** — route by geography across regional pools
- **03-weighted-ab-testing** — split traffic across pools by weight
