# CDN-Cached API

A read-heavy API served through Cloud CDN with the origin in control: `USE_ORIGIN_HEADERS` caches exactly what the application marks cacheable, tracking parameters are stripped from the cache key, and stale content bridges origin blips.

## When to Use

- Public read-heavy APIs (catalogs, content feeds, search suggestions) where some routes are cache-safe and others are not
- Any origin that already emits correct `Cache-Control` headers

## Remix Notes

- `USE_ORIGIN_HEADERS` forbids the TTL fields (the spec enforces the pairing) — lifetimes belong to the origin's headers in this mode.
- If every response is safe to cache, `CACHE_ALL_STATIC` with explicit `defaultTtl`/`clientTtl` needs no origin header discipline.
- Vary responses by header (e.g. `Accept-Language`)? Add it to `cacheKeyPolicy.includeHttpHeaders` — every distinct value gets its own cache entry, so keep the list short.
- CDN only fronts external schemes; the spec rejects `enableCdn` on INTERNAL_* before deploy.
