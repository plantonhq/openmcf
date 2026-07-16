---
title: "Catch-All HTTPS Route"
description: "This preset creates the endpoint's default route: every path, both protocols, HTTP redirected to HTTPS at the edge, no caching. The standard production entry rule for a dynamic application."
type: "preset"
rank: "01"
presetSlug: "01-catch-all-https"
componentSlug: "front-door-route"
componentTitle: "Front Door Route"
provider: "azure"
icon: "package"
order: 1
---

# Catch-All HTTPS Route

This preset creates the endpoint's default route: every path, both
protocols, HTTP redirected to HTTPS at the edge, no caching. The
standard production entry rule for a dynamic application.

## When to Use

- The first (often only) route on an endpoint serving a dynamic web app
  or API
- As the fallback beside more specific routes ("/api/*", "/static/*")
  -- Front Door picks the most specific pattern, so the catch-all only
  sees what nothing else matched

## Key Configuration Choices

- **Both protocols with the redirect left on** (Azure's default) -- the
  redirect needs HTTP to arrive and HTTPS to land on; the spec rejects
  redirect-with-one-protocol before Azure would
- **No `cache` block** -- deliberately absent: dynamic responses are
  fetched from the origin every time; add the block only for cacheable
  content (see the static-assets preset)
- **`originIds` for composed deploys** -- the ordering seam when route,
  origins, and group ship in one manifest set; omit when the origins
  already exist

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<endpoint-resource-name>` | The AzureFrontDoorEndpoint's Planton resource name | Your Front Door composition |
| `<origin-group-resource-name>` | The AzureFrontDoorOriginGroup's Planton resource name | Your Front Door composition |
| `<origin-resource-name>` | The AzureFrontDoorOrigin's Planton resource name | Your Front Door composition |
| `routeName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
