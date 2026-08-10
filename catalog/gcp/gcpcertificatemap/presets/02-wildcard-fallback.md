# Wildcard Fallback

One wildcard certificate serving every subdomain through a single
PRIMARY entry — the simplest certificate map that can work.

## What it configures

- A single `matcher: PRIMARY` entry bound to a wildcard certificate:
  every handshake (any SNI, or none) is served by it.
- `PREVENT` teardown — live TLS routing is protected.

## Adjust before deploying

- **certificates** — reference your wildcard GcpCertManagerCert's
  `certificate_id` output (`*.example.com` plus the apex as SANs).

## After deploying

Set a GcpTargetHttpsProxy's `certificate_map` to this map's `map_uri`
output. Grow per-domain entries onto the map later (dedicated EV cert
for www, separate api cert) — the PRIMARY stays as the safety net.

## When to choose something else

Different certificates per domain from day one? Start from the
**Hostname Routing** preset.
