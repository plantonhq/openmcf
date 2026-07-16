---
title: "Cached Static Assets Route"
description: "This preset creates a route dedicated to static content: matched paths are cached at every Front Door edge location and text-based types are compressed on the way out. The origin sees each asset..."
type: "preset"
rank: "02"
presetSlug: "02-static-assets-cached"
componentSlug: "front-door-route"
componentTitle: "Front Door Route"
provider: "azure"
icon: "package"
order: 2
---

# Cached Static Assets Route

This preset creates a route dedicated to static content: matched paths
are cached at every Front Door edge location and text-based types are
compressed on the way out. The origin sees each asset roughly once per
edge, not once per user.

## When to Use

- Static assets (JS/CSS bundles, images, fonts) served beside a dynamic
  app on the same endpoint
- Storage static websites or SPA bundles where cache hit ratio is the
  whole value of the CDN

## Key Configuration Choices

- **`IGNORE_QUERY_STRING`** -- every query variant shares one cache
  entry; right for assets that version by PATH (`app.3f9c.js`). Use
  `USE_QUERY_STRING` only when query strings genuinely change content
- **Compression list covers text types only** -- binary media (images,
  video, archives) is already compressed; Azure also only compresses
  responses between 1 KiB and 8 MiB
- **Patterns are the split** -- this route takes `/static/*` and
  `/assets/*`; the uncached catch-all handles everything else. Front
  Door picks the most specific match, so ordering is automatic

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<endpoint-resource-name>` | The AzureFrontDoorEndpoint's Planton resource name | Your Front Door composition |
| `<origin-group-resource-name>` | The AzureFrontDoorOriginGroup's Planton resource name (often a storage-backed group) | Your Front Door composition |
| `routeName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
