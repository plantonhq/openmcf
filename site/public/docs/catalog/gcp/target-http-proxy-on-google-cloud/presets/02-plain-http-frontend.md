---
title: "Plain HTTP Frontend"
description: "A target HTTP proxy that serves the application itself over plain HTTP — the right shape for internal test environments, health-probe endpoints, or services that terminate TLS elsewhere."
type: "preset"
rank: "02"
presetSlug: "02-plain-http-frontend"
componentSlug: "target-http-proxy-on-google-cloud"
componentTitle: "Target HTTP Proxy on Google Cloud"
provider: "gcp"
icon: "package"
order: 2
---

# Plain HTTP Frontend

A target HTTP proxy that serves the application itself over plain HTTP — the right shape for internal test environments, health-probe endpoints, or services that terminate TLS elsewhere.

## When to Use

- Pre-production environments where TLS provisioning is not worth the setup cost
- Load testing the backend path without TLS handshake overhead
- Serving traffic behind an upstream TLS terminator you do not control

## Remix Notes

- `httpKeepAliveTimeoutSec` is only honored by the envoy-based `EXTERNAL_MANAGED` load balancer; drop it if the forwarding rule uses the classic `EXTERNAL` scheme.
- Do NOT ship this to production for user-facing traffic — promote to the HTTPS-redirect pattern (preset 01) plus a `GcpTargetHttpsProxy` once a certificate exists.
- Reference the `GcpUrlMap` via `valueFrom` instead of a literal self-link.
