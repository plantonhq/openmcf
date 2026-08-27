# CloudflareHealthcheck guide

The judgment this guide protects you from: this is not a load-balancer monitor, it is a paid zone feature, and the config block must match the type. Reach for the wrong sibling -- or send `http_config` on a TCP check -- and the failure is at apply or as refresh drift.

## Declare check_regions, and give it a real address

Two walls measured live (2026-08-27): (1) always declare `check_regions` explicitly on Terraform-managed checks -- when omitted, Cloudflare picks and echoes a default region (measured: `WNAM`) into a provider attribute modeled optional-not-computed, so every later plan proposes a no-op update forever; (2) the `address` must be a real, publicly routable IP or hostname -- documentation-range addresses (TEST-NET) are rejected at create with "origin address is invalid" (code 1002). The origin does not have to answer the probe; an unhealthy check is still a fully manageable check.

## This is not CloudflareLoadBalancerMonitor

Standalone health checks watch an origin and record status. They need no load balancer. `CloudflareLoadBalancerMonitor` is the account-scoped monitor a pool consumes to drive failover. They do not share IDs, APIs, or import formats. If you are about to paste a monitor id into a pool, you are on the wrong kind.

## Paid zone feature: Pro+

Cloudflare enforces the plan gate at create. A free zone is rejected at the API (403/402), not here. Pro includes a small allotment; Business and Enterprise more. Extra checks beyond the allotment are a zone-plan / add-on decision -- real money, billed on the zone, not through this component.

## type is HTTP, HTTPS, or TCP -- and the config must match

The provider accepts any string for `type` and rejects bad values only at the API. This spec tightens the wall to the documented set: `HTTP`, `HTTPS`, `TCP`.

`http_config` is only valid on HTTP/HTTPS. `tcp_config` is only valid on TCP. Sending the wrong block is rejected at validation. The module also never sends the unused block, because both are Computed upstream and the wrong one reads back as drift.

`http_config.headers` is a map of name → `{values: [...]}` because proto3 maps cannot hold repeated values. The provider argument is `header` (a plain map of string lists). Set a `Host` header for name-based virtual hosts. The User-Agent header cannot be overridden.

For HTTPS on port 443, set `http_config.port` explicitly -- Cloudflare's HTTP default port is 80.

## API-token auth defect class

Some upstream provider tests blank `CLOUDFLARE_API_TOKEN` and fall through to the global API key. That class of auth defect is real: a token-only environment can fail in ways a key+email environment does not, and the error is often an opaque 403 rather than "wrong credential type." Prefer a scoped API token with Zone → Health Checks → Edit, and do not mix token and key in the same process. If a create fails with 403 on a Pro+ zone, check the credential before you check the allotment.

## Destroy is a real delete

Destroy removes the check and its history. `suspended: true` pauses probing and keeps the ID -- use that when you want the object to stay.

## Pairs well with

- [CloudflareLoadBalancerMonitor](../cloudflareloadbalancermonitor/README.md) -- the other health-check family; read the boundary above before you pick.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`; the zone plan is what gates create.
