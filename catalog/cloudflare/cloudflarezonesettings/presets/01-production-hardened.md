# Production Hardened

The standard production security posture for a zone: every plain-http request redirects to HTTPS, connections to the origin are encrypted and certificate-validated (`ssl: strict`), TLS below 1.2 is refused, and browsers receive a one-year HSTS pin covering subdomains. HTTP/3 is enabled and the security level sits at `medium`. Every other zone setting stays unmanaged and untouched.

## When to Use

- Any zone serving production HTTPS traffic where the origin presents a valid certificate
- Codifying the security decisions usually made by hand in the dashboard's SSL/TLS and Security tabs
- A baseline to extend with performance settings or companions once the security posture is locked in

## Key Configuration Choices

- **ssl: strict** -- Cloudflare validates the origin certificate. If your origin serves a self-signed or expired certificate, fix that first (a CloudflareOriginCaCertificate works); dropping to `full` gives up validation.
- **min_tls_version: "1.2"** -- quoted, because unquoted 1.2 is a YAML number. Raise to "1.3" only after confirming your clients support it.
- **security_header** -- `max_age: 31536000` is a one-year browser pin; `include_subdomains: true` requires every subdomain to serve HTTPS for that year. `preload` stays false because preload-list removal is slow and manual -- enable it deliberately.
- **security_level: medium** -- the balanced challenge posture; move to `high` or `under_attack` during incidents (and remember to move back: destroying this resource does not revert anything).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone ID these settings manage | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |

## Related Presets

- **02-minimal** -- just HTTPS enforcement and the TLS floor, nothing else
- **03-performance** -- the speed-side counterpart: compression, caching, and image optimization
