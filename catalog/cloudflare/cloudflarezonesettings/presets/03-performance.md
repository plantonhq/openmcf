# Performance

The speed-side settings for a content-serving zone: Brotli compression, Early Hints (HTTP 103 from cached responses while the origin thinks), HTTP/3 for client connections, aggressive edge caching of static content across query-string variants, a four-hour browser cache TTL, and lossless image optimization via Polish. Security settings are left unmanaged -- pair with a security preset or add those fields here.

## When to Use

- Marketing sites, blogs, and documentation zones where static content dominates
- Zones on Pro plan or above (Polish is plan-gated; the apply fails on Free rather than billing anything)
- Tuning delivery on a zone whose security posture is already managed elsewhere

## Key Configuration Choices

- **cache_level: aggressive** -- caches static content with all query-string variants as separate entries. Use `basic` if query strings should bypass cache.
- **browser_cache_ttl: 14400** -- four hours; Cloudflare accepts a fixed value set (0 and specific durations from 30 seconds to one year), and 0 means respect the origin's cache headers.
- **polish: lossless** -- strips image metadata without touching pixel data. Switch to `lossy` for further savings on JPEG-heavy zones; pair with `webp: true` to serve WebP variants.
- **early_hints + http3** -- both are free-plan safe and benefit repeat visitors most.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone ID these settings manage | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |

## Related Presets

- **01-production-hardened** -- the security-side counterpart: HTTPS enforcement, strict SSL, HSTS
- **02-minimal** -- the smallest useful starting point
