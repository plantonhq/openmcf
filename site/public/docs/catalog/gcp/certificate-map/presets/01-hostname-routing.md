---
title: "Hostname Routing"
description: "Per-domain certificates selected by SNI at the load balancer, with a wildcard PRIMARY fallback — the multi-domain TLS edge as a routing table."
type: "preset"
rank: "01"
presetSlug: "01-hostname-routing"
componentSlug: "certificate-map"
componentTitle: "Certificate Map"
provider: "gcp"
icon: "package"
order: 1
---

# Hostname Routing

Per-domain certificates selected by SNI at the load balancer, with a
wildcard PRIMARY fallback — the multi-domain TLS edge as a routing
table.

## What it configures

- Two hostname entries (www, api) each bound to its own certificate.
- A `PRIMARY` fallback entry serving the wildcard when no hostname
  matches (including SNI-less clients).
- `PREVENT` teardown — live TLS routing is protected.

## Adjust before deploying

- **certificates** — reference your GcpCertManagerCert resources'
  `certificate_id` outputs.
- **hostnames** — add an entry per domain; wildcards
  (`*.example.com`) match a suffix set.
- Keep the **PRIMARY** entry — removing it hard-fails unmatched
  handshakes.

## After deploying

Set a GcpTargetHttpsProxy's `certificate_map` to this map's `map_uri`
output. Rotate certificates by editing an entry's list in place —
attach the new one before detaching the old.

## When to choose something else

One wildcard covering everything? Start from the **Wildcard Fallback**
preset — a single PRIMARY entry, no per-domain routing to maintain.
