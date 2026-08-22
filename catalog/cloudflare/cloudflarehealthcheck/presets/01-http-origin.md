# HTTP origin probe

A standalone HTTP health check against `example.com/health`. No load balancer is involved -- this kind watches an origin. Health checks are a paid zone feature (Pro+); a free zone is rejected at the API.

## When to Use

- First standalone health check on a Pro+ zone
- An origin that is not behind a Cloudflare load balancer
- Prefer this over `CloudflareLoadBalancerMonitor` when no pool consumes the result

## Key Configuration Choices

- **type: HTTP** -- the default. Use HTTPS (and set `http_config.port` to 443) for TLS origins; use TCP with `tcp_config` for a handshake check.
- **http_config.path: /health** -- Cloudflare's default path is `/`. expected_codes of `200` matches the API default; add `2xx` if any success code is healthy.
- **No tcp_config** -- that block is only valid when type is TCP. Sending it here is rejected.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone the health check belongs to | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `address` | Origin hostname or IP | Your origin inventory |
| `http_config.path` | Path probed | Your origin's health endpoint |

## Related Presets

None yet -- a TCP variant belongs on a zone that already has the Health Checks entitlement.
