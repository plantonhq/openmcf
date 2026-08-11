---
title: "Traffic-Managed Canary"
description: "A weighted canary hardened with the per-route traffic-management surface: bounded timeouts, deliberate retries, CORS answered at the load balancer, a response header attributing which arm served each..."
type: "preset"
rank: "04"
presetSlug: "04-traffic-managed-canary"
componentSlug: "url-map-on-google-cloud"
componentTitle: "URL Map on Google Cloud"
provider: "gcp"
icon: "package"
order: 4
---

# Traffic-Managed Canary

A weighted canary hardened with the per-route traffic-management surface: bounded timeouts, deliberate retries, CORS answered at the load balancer, a response header attributing which arm served each request, and a destroy stance that protects shared routing.

## When to Use

- Canary rollouts where dashboards must attribute errors to the canary arm, not the split
- APIs whose clients need retries and timeouts enforced at the edge, uniformly, without per-backend drift
- Browser-facing APIs that should answer CORS preflights before traffic reaches any backend

## Remix Notes

- The header attribution rides the canary arm only (`headerAction` on the weighted backend) — the stable arm stays untouched, so the header's presence IS the canary signal.
- `timeout` bounds the whole exchange including retries; `perTryTimeout` bounds one attempt. 15s total with 5s per try means both retries can actually run — keep that quotient honest when you change either.
- Add `requestMirrorPolicy` to shadow production traffic onto a candidate stack — mirror read-only routes only, and size the mirror backend for the full load.
- Add a route-scoped `cachePolicy` for CDN-backed routes; set `requestCoalescing: true` explicitly when you do (an unset boolean inside a written cache policy sends false).
- `deletionPolicy: PREVENT` fails any destroy at the map — swap to DELETE for stacks that tear down whole, e2e fixtures included.
