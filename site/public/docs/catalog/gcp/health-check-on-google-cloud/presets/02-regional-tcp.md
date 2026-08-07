---
title: "Regional TCP Probe"
description: "A regional health check proving TCP connectability on a fixed port — the shape internal load balancers and regional managed instance groups consume. Regional backend services can only reference..."
type: "preset"
rank: "02"
presetSlug: "02-regional-tcp"
componentSlug: "health-check-on-google-cloud"
componentTitle: "Health Check on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Regional TCP Probe

A regional health check proving TCP connectability on a fixed port — the shape internal load balancers and regional managed instance groups consume. Regional backend services can only reference health checks in their own region, so the `region` here must match theirs.

## When to Use

- Internal (regional) load balancers fronting databases, caches, or custom TCP services
- Regional MIG auto-healing where no HTTP endpoint exists
- Any service where "accepting connections" is the best available liveness signal

## Remix Notes

- A TCP probe is blind to application-level sickness — a deadlocked process with an open listener still passes. If the service has any HTTP endpoint, prefer the HTTP preset.
- Add `request`/`response` to verify an application banner (e.g. a protocol greeting) instead of bare connectivity.
- Switch to the `ssl` block if the service only accepts TLS connections.
- Instance-group backends need a firewall rule admitting `35.191.0.0/16` and `130.211.0.0/22` on the probed port.
