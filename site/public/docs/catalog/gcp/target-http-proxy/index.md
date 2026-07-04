---
title: "Target HTTP Proxy"
description: "Target HTTP Proxy deployment documentation"
icon: "package"
order: 100
componentName: "gcptargethttpproxy"
---

# GCP Target HTTP Proxy

Creates a global Compute Engine target HTTP proxy — the plaintext-HTTP frontend adapter that binds a global forwarding rule (the VIP) to a URL map (the routing table). Its standard production role is serving the http→https redirect on port 80 while a target HTTPS proxy serves the application on 443.

## What Gets Created

A single `google_compute_target_http_proxy` in the chosen GCP project.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A URL map** — a `GcpUrlMap` self-link for the proxy to route through
- **IAM permissions** — any role carrying `compute.targetHttpProxies.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpTargetHttpProxy
metadata:
  name: web-http-frontend
spec:
  projectId:
    value: my-gcp-project-123
  urlMap:
    value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/urlMaps/http-redirect
```

```shell
planton apply -f http-proxy.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `urlMap` | `StringValueOrRef` | — | Required. The URL map to route through. Mutable in place (zero-downtime repoint) |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the proxy. Immutable |
| `proxyName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable |
| `description` | `string` | — | What this proxy fronts. Immutable |
| `httpKeepAliveTimeoutSec` | `int32` | GCP default (610) | Idle client keep-alive, 5-1200s; `EXTERNAL_MANAGED` only. Immutable |
| `proxyBind` | `bool` | `false` | Bind to Traffic Director mesh VIPs (`INTERNAL_SELF_MANAGED` only). Immutable |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a global forwarding rule references as its target |
| `proxy_name` | Name of the proxy in GCP |
| `proxy_id` | Server-assigned numeric ID |
| `fingerprint` | Fingerprint for optimistic concurrency control |

## Related Components

- [GcpUrlMap](/docs/catalog/gcp/url-map) — the routing table this proxy consults
- [GcpGlobalForwardingRule](/docs/catalog/gcp/global-forwarding-rule) — the VIP that binds to this proxy
- [GcpTargetHttpsProxy](/docs/catalog/gcp/target-https-proxy) — the TLS-terminating sibling
- [GcpProject](/docs/catalog/gcp/project) — provides the project that owns the proxy
