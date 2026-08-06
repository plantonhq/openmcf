# Plain Origin (No CDN)

A backend bucket with edge caching off: the load balancer proxies every request to the bucket. The simplest way to put GCS content behind the same global VIP, hostname, and URL map as the rest of an application.

## When to Use

- Content that changes in place under a stable path (no fingerprinting) and must never be served stale
- Low-traffic download paths where CDN cache-fill cost outweighs the benefit
- A first step: ship it plain, flip `enableCdn` later — that flip is an in-place update

## Remix Notes

- Turning CDN on later is one field (`enableCdn: true`) plus an optional `cdnPolicy` — the URL map and everything above it are untouched.
- Add `customResponseHeaders` for security headers (e.g. `Strict-Transport-Security: max-age=31536000`) applied uniformly by the load balancer.
- For cross-region internal ALBs set `loadBalancingScheme: INTERNAL_MANAGED` (CDN must stay off there — the spec enforces it).
