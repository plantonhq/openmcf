# Minimal

Manages exactly two settings -- redirect all plain-http traffic to HTTPS and refuse TLS below 1.2 -- and leaves the other 61 settings unmanaged. This is the smallest useful zone-settings resource and the safest starting point: nothing you have configured in the dashboard is touched.

## When to Use

- First adoption of managed zone settings on a zone that already has hand-tuned dashboard configuration
- Zones where the only hard requirement is HTTPS enforcement with a modern TLS floor
- A starting point to grow field by field, taking ownership of one setting at a time

## Key Configuration Choices

- **Unset fields stay unmanaged** -- the module only sends what you set. Add fields as you decide to own them; each one you add becomes yours until you explicitly set it to something else (zone settings have no delete at Cloudflare).
- **min_tls_version: "1.2"** -- quoted, because unquoted 1.2 is a YAML number, not the string the API expects.
- **No ssl mode set** -- the zone's current origin-encryption mode stays as-is. Add `ssl: strict` once the origin presents a valid certificate.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone ID these settings manage | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |

## Related Presets

- **01-production-hardened** -- the full production security posture with strict SSL and HSTS
- **03-performance** -- compression, caching, and image optimization settings
