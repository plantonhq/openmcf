# Geo-located pool (proximity steering)

A regional pool tagged with latitude/longitude and a least-connections origin
policy, for use with a load balancer's `proximity` or `geo` steering.

## When to use

- You run regional backends and want Cloudflare to send users to the nearest
  healthy pool (proximity steering) or route by geography (`region_pools`).

## Key choices

- `latitude` / `longitude`: the pool's data-center coordinates, used by proximity
  steering to compute distance.
- `checkRegions`: probe from the region nearest the origins to reflect real health.
- `originSteering.policy`: `least_connections` (or `least_outstanding_requests`)
  spreads load by live connection counts rather than at random.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
| `<origin-address>` | Origin IP or hostname in this region |

## Wiring into geo steering

Reference this pool from a load balancer's `regionPools` (by code) or rely on
`proximity` steering using the coordinates above.
