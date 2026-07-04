---
title: "Backend Service"
description: "Backend Service deployment documentation"
icon: "package"
order: 100
componentName: "gcpbackendservice"
---

# GCP Backend Service

Creates a global Compute Engine backend service — the hub of GCP's L7 load balancing family. It owns how traffic reaches a set of backends: the backend list (instance groups / NEGs), health checking, session affinity, Cloud CDN policy, Identity-Aware Proxy, Cloud Armor attachment, and request logging. URL maps route host/path patterns to backend services.

## What Gets Created

When you deploy a GcpBackendService resource, Planton provisions:

- **Backend Service** — a global `google_compute_backend_service` with its backends, health check, affinity, CDN policy, IAP, Cloud Armor references, logging, and traffic policy
- **Signed-URL Keys** (optional) — up to 3 named signing keys for serving private CDN content

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **A health check** — referenced via `healthCheck`; required unless every backend is an internet or serverless NEG
- **IAM permissions** — any role carrying `compute.backendServices.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpBackendService
metadata:
  name: web-backend
spec:
  projectId:
    value: my-gcp-project-123
  healthCheck:
    valueFrom:
      kind: GcpHealthCheck
      name: web-hc
      fieldPath: status.outputs.self_link
  backends:
    - group:
        value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/zones/us-central1-a/instanceGroups/web-ig
      balancingMode: UTILIZATION
      maxUtilization: 0.8
```

```shell
planton apply -f backend-service.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `projectId` | `StringValueOrRef` | provider default | Project that owns the backend service. Immutable. |
| `backendServiceName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |
| `protocol` | `string` | `HTTP` | LB→backend protocol: HTTP, HTTPS, HTTP2, H2C, TCP, SSL, UDP, GRPC. |
| `loadBalancingScheme` | `string` | `EXTERNAL` | EXTERNAL, EXTERNAL_MANAGED, INTERNAL_MANAGED, INTERNAL_SELF_MANAGED. |
| `healthCheck` | `StringValueOrRef` | none | The ONE health check (GCP caps at one). Can reference a GcpHealthCheck. |
| `backends` | `list(object)` | `[]` | Instance groups / NEGs with balancing mode and capacity dials. |
| `sessionAffinity` | `string` | `NONE` | Client stickiness incl. cookie modes and STRONG_COOKIE_AFFINITY. |
| `enableCdn` + `cdnPolicy` | `bool` + object | off | Edge caching with the full cache-key policy (host/protocol/query/cookies/headers). |
| `securityPolicy` / `edgeSecurityPolicy` | `StringValueOrRef` | none | Cloud Armor backend/EDGE policies, attached by reference. |
| `iap` | object | off | Identity-Aware Proxy; `oauth2ClientSecret` is secret material. |
| `logConfig` | object | off | Request logging with sampling. |
| `circuitBreakers` / `outlierDetection` / `consistentHash` / `maxStreamDuration` | objects | none | Traffic Director / service-mesh traffic policy. |
| `securitySettings` / `tlsSettings` | objects | none | Backend mTLS, SAN pinning, SNI, AWS SigV4 origin signing (`accessKey` is secret). |
| `signedUrlKeys` | `list(object)` | `[]` | Up to 3 signing keys for CDN signed URLs; `keyValue` is secret material. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value URL maps reference as a default service or path-rule target |
| `backend_service_name` | Name of the backend service in GCP |
| `generated_id` | Server-assigned numeric ID |
| `fingerprint` | Optimistic-concurrency fingerprint |

## Related Components

- [GcpHealthCheck](/docs/catalog/gcp/health-check) — the probe deciding which backends receive traffic
- [GcpBackendBucket](/docs/catalog/gcp/backend-bucket) — the static-content counterpart for GCS origins
- [GcpCloudArmorPolicy](/docs/catalog/gcp/cloud-armor-policy) — WAF/rate-limiting policies attached by reference
- [GcpProject](/docs/catalog/gcp/project) — provides the project that owns the backend service
