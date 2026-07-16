---
title: "HTTPS-Only API Route"
description: "This preset creates a strict API route: HTTPS is the only accepted protocol (plain HTTP fails instead of redirecting), the origin leg is pinned to HTTPS, and the public path prefix is rewritten on..."
type: "preset"
rank: "03"
presetSlug: "03-api-https-only"
componentSlug: "front-door-route"
componentTitle: "Front Door Route"
provider: "azure"
icon: "package"
order: 3
---

# HTTPS-Only API Route

This preset creates a strict API route: HTTPS is the only accepted
protocol (plain HTTP fails instead of redirecting), the origin leg is
pinned to HTTPS, and the public path prefix is rewritten on the origin
side.

## When to Use

- REST/GraphQL APIs where a silent 301 on HTTP is worse than an error
  (clients mishandle redirected POST bodies; failing loudly surfaces
  misconfigured base URLs)
- Backends whose path layout differs from the public URL scheme

## Key Configuration Choices

- **`supportedProtocols: [HTTPS]` + `httpsRedirectEnabled: false`** --
  the redirect requires both protocols, so single-protocol routes must
  disable it (the spec enforces the pairing); HTTP connections are
  refused at the edge
- **`forwardingProtocol: HTTPS_ONLY`** -- never downgrade on the origin
  leg, regardless of route or client configuration
- **`originPath: /v1`** -- Front Door PREPENDS the origin path to the
  request path (a call to `/api/users` fetches `/v1/api/users` from the
  backend); it namespaces one origin under a directory, it does not
  strip the public prefix. Drop it when public and origin paths match
- **No cache** -- API responses are dynamic by default; add the cache
  block per-route where responses are genuinely cacheable

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<endpoint-resource-name>` | The AzureFrontDoorEndpoint's Planton resource name | Your Front Door composition |
| `<origin-group-resource-name>` | The AzureFrontDoorOriginGroup's Planton resource name | Your Front Door composition |
| `routeName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
