# GCP Target HTTP Proxies: The Thin Adapter Between VIP and Routing Table

## Why This Resource Exists At All

Google's global load balancer is deliberately decomposed: the forwarding rule owns the IP and port, the URL map owns routing, the backend service owns traffic policy. The target HTTP proxy is the small adapter that connects the first two — a forwarding rule cannot reference a URL map directly, it references a *target*, and the proxy is the target type that speaks plaintext HTTP.

That thinness is the point. Because the proxy is its own resource with a mutable `urlMap` reference, a live frontend can be repointed at a brand-new routing table (GCP's dedicated `setUrlMap` call swaps it in place) without touching the VIP, DNS, or TLS. The proxy is the pivot for zero-downtime routing migrations.

## The Pair Pattern (Where Nearly Every HTTP Proxy Belongs)

In production, plaintext HTTP exists to be escaped from. The near-universal deployment shape is a PAIR of proxies sharing one frontend story:

- A **target HTTP proxy** (this kind) on port 80, whose URL map is redirect-only: `defaultUrlRedirect { httpsRedirect: true }` returning 301s. No backends anywhere in that map.
- A **target HTTPS proxy** on port 443 serving the real application.
- Two **global forwarding rules** — one per port — sharing a single reserved static IP (allowed exactly because their port ranges do not overlap).

A target HTTP proxy serving real application traffic is a legitimate but niche shape (internal test environments, upstream TLS termination) — the presets model both.

## What Is Deliberately Immutable

Only `url_map` updates in place. Name, description, keep-alive, and `proxy_bind` are ForceNew on the provider: changing them destroys and recreates the proxy, and any forwarding rule referencing the old self-link breaks until it repoints. The spec documents this on every field so an operator knows which edits are safe on a live frontend.

## http_keep_alive_timeout_sec: One Field, One Load Balancer Family

The keep-alive dial is only honored by the envoy-based `EXTERNAL_MANAGED` global external ALB (default 610s there); the classic `EXTERNAL` ALB ignores it. The value must exceed your clients' own keep-alive to avoid the load balancer closing first. The spec validates the 5-1200 range and treats 0 as "let GCP decide" — the same zero-means-unset convention the catalog uses everywhere.

## The 90/10 Coverage Decision

| Provider field | Modeled | Notes |
|---|---|---|
| `name` | ✅ `proxyName` | Defaults to `metadata.name`; RFC1035 validated |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; empty → provider default |
| `description` | ✅ | |
| `url_map` | ✅ `urlMap` | Required ref → GcpUrlMap self_link; the only mutable field |
| `http_keep_alive_timeout_sec` | ✅ | 5-1200 CEL; 0 = GCP default |
| `proxy_bind` | ✅ `proxyBind` | Traffic Director binding |
| `proxy_id` / `fingerprint` (computed) | outputs only | |
| `deletion_policy` | ❌ | Absent from the released 6.x line; Terraform-provider lifecycle lever, Planton owns destroy semantics |
| `timeouts` | ❌ | Operation plumbing |

## Composition

The proxy is the second-to-last node walking outward from backends:

1. **GcpBackendService / GcpBackendBucket** — where traffic lands.
2. **GcpUrlMap** — decides which backend gets each request.
3. **GcpTargetHttpProxy** (this component) — the frontend adapter forwarding rules bind to.
4. **GcpGlobalForwardingRule** — the VIP; references this proxy's `self_link` as its `target`.

The kind registry declares `GcpUrlMap` as a prerequisite, so the E2E harness (and any composed chart) installs the routing table before the proxy.

## Operational Notes

- **Zero-downtime routing swap**: apply a new URL map, then update this proxy's `urlMap` — GCP repoints in place.
- **Recreate blast radius**: renaming the proxy briefly orphans its forwarding rule; do it by creating the replacement proxy first and repointing the rule's `target` (also in-place).
- The module enables `compute.googleapis.com` before creating the proxy, so a fresh project works on the first deploy.
