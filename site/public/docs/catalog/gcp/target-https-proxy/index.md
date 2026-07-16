---
title: "Target HTTPS Proxy"
description: "Target HTTPS Proxy deployment documentation"
icon: "package"
order: 100
componentName: "gcptargethttpsproxy"
---

# GCP Target HTTPS Proxy

Creates a global Compute Engine target HTTPS proxy — the TLS-termination node that binds a global forwarding rule (the VIP) to a URL map (the routing table) and owns the client-facing handshake: certificates, SSL policy, QUIC (HTTP/3), and TLS 1.3 early data.

## What Gets Created

A single `google_compute_target_https_proxy` in the chosen GCP project.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A URL map** — a `GcpUrlMap` self-link for the proxy to route through
- **A certificate source** — typically one or more `GcpManagedSslCertificate` resources
- **IAM permissions** — any role carrying `compute.targetHttpsProxies.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpTargetHttpsProxy
metadata:
  name: web-https-frontend
spec:
  projectId:
    value: my-gcp-project-123
  urlMap:
    value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/urlMaps/web-routing
  sslCertificates:
    - value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/sslCertificates/web-cert
```

```shell
planton apply -f https-proxy.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `urlMap` | `StringValueOrRef` | — | Required. Routing table for decrypted requests. Mutable in place |
| `sslCertificates` | list of refs | `[]` | 1-15 compute SSL certificates; exactly one certificate mechanism allowed |
| `certificateManagerCertificates` | list of refs | `[]` | Certificate Manager certs (cross-region internal ALB only) |
| `certificateMap` | `string` | — | SNI-scale certificate map URI (external ALBs only) |
| `sslPolicy` | `StringValueOrRef` | GCP default policy | TLS version/cipher constraints. Mutable |
| `serverTlsPolicy` | `StringValueOrRef` | — | mTLS ServerTlsPolicy; Traffic Director's only TLS lever. Mutable, clearable |
| `quicOverride` | `string` | `NONE` | QUIC negotiation: `NONE` / `ENABLE` / `DISABLE`. Mutable |
| `tlsEarlyData` | `string` | GCP default (`DISABLED`) | TLS 1.3 0-RTT mode. Immutable |
| `httpKeepAliveTimeoutSec` | `int32` | GCP default (610) | Idle keep-alive, 5-1200s; `EXTERNAL_MANAGED` only. Immutable |
| `proxyBind` | `bool` | `false` | Traffic Director mesh binding. Immutable |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the proxy. Immutable |
| `proxyName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value a global forwarding rule references as its target |
| `proxy_name` | Name of the proxy in GCP |
| `proxy_id` | Server-assigned numeric ID |
| `fingerprint` | Fingerprint for optimistic concurrency control |

## Related Components

- [GcpManagedSslCertificate](/docs/catalog/gcp/managed-ssl-certificate) — Google-managed certificates for `sslCertificates`
- [GcpCertManagerCert](/docs/catalog/gcp/certificate-manager-certificate) — Certificate Manager certificates
- [GcpUrlMap](/docs/catalog/gcp/url-map) — the routing table this proxy consults
- [GcpGlobalForwardingRule](/docs/catalog/gcp/global-forwarding-rule) — the VIP that binds to this proxy
- [GcpTargetHttpProxy](/docs/catalog/gcp/target-http-proxy) — the port-80 redirect sibling
