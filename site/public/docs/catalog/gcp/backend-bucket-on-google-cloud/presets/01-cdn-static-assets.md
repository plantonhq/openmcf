---
title: "CDN-Cached Static Assets"
description: "The workhorse backend bucket: fingerprinted static assets (JS/CSS/images with content hashes in their paths) served through Cloud CDN with sensible TTLs, negative caching for missing-asset storms,..."
type: "preset"
rank: "01"
presetSlug: "01-cdn-static-assets"
componentSlug: "backend-bucket-on-google-cloud"
componentTitle: "Backend Bucket on Google Cloud"
provider: "gcp"
icon: "package"
order: 1
---

# CDN-Cached Static Assets

The workhorse backend bucket: fingerprinted static assets (JS/CSS/images with content hashes in their paths) served through Cloud CDN with sensible TTLs, negative caching for missing-asset storms, request coalescing, and automatic compression.

## When to Use

- The `/assets/*` or `/static/*` path of a web application behind a global external HTTP(S) load balancer
- Any immutable, versioned content where cache invalidation is handled by changing the path (fingerprinting)

## Remix Notes

- Fingerprinted assets can take much longer TTLs (`defaultTtl: 86400` or more) — the path changes on every release, so stale content is impossible.
- If the origin sets correct `Cache-Control` headers itself, switch `cacheMode` to `USE_ORIGIN_HEADERS` and remove the TTL fields (the spec enforces that pairing).
- The `X-Cache-Status` header is free observability; remove it if response-header hygiene matters more.
- Objects must be publicly readable (`allUsers` → `roles/storage.objectViewer` on the bucket) — the load balancer does not authenticate to GCS.
