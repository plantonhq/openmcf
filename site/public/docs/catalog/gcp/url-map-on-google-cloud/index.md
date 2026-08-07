---
title: "URL Map on Google Cloud"
description: "URL Map on Google Cloud deployment documentation"
icon: "package"
order: 100
componentName: "gcpurlmap"
---

# URL Map on Google Cloud

Deploys a global Compute Engine URL map — the L7 routing brain of a global external Application Load Balancer. The URL map matches each request's Host and path and decides what happens: send it to a backend service or backend bucket, split it across weighted backends (canary/blue-green), redirect, or rewrite. Target proxies bind this map's self link; the forwarding rule and its IP sit in front of the proxy. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to projects, backend services, and backend buckets.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Engine URL Map (global)** -- host rules, path matchers with path or route rules, default targets at every level, header policies, custom error pages, and routing self-tests

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the URL map will be created.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.
- **At least one backend** (GcpBackendService or GcpBackendBucket) — the default target is required, and it is the natural creation order: backends first, then the routing that references them.

## Deploy

### Console

Open the deployment store, find **URL Map on Google Cloud**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Host Path Fanout** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpUrlMap
metadata:
  name: prod-web-map
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  defaultService:
    value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/backendServices/web-api"
  hostRules:
    - hosts: [acme.com, www.acme.com]
      pathMatcher: site
  pathMatchers:
    - name: site
      defaultService:
        value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/backendServices/web-api"
      pathRules:
        - paths: ["/assets/*"]
          service:
            value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/backendBuckets/assets-backend"
```

```shell
planton apply -f url-map.yaml
```

This creates the classic fan-out: dynamic traffic to a backend service, /assets/* to a CDN-backed bucket.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the backends:

```yaml
spec:
  defaultService:
    valueFrom:
      kind: GcpBackendService
      name: web-api
      fieldPath: status.outputs.self_link
```

The InfraPipeline resolves the dependency graph — backends first, then this URL map — and a downstream GcpTargetHttpsProxy references this map's `self_link`.

## Key Configuration

These are the most important decisions when configuring a URL map. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Default target** -- Exactly one of: a backend service/bucket, a URL redirect (the http→https front half), or a weighted split across backend services. All of it is MUTABLE — repointing a live map is an in-place, zero-downtime update, which is the canary lever. On the backend arm, a host/path rewrite may accompany the service.

**Path matchers and rules** -- Named routing tables host rules point at. Each matcher uses path rules (longest prefix — the simple choice) OR route rules (priority-ordered with header/query matching), never both. Every rule target is the same three-arm choice, so a weighted canary can be scoped to exactly one path.

**Host rules** -- Fan hosts out to matchers by name ("*.acme.com" covers subdomains, not the apex). Every host routed here must also be covered by a certificate on the target HTTPS proxy.

**Routing tests** -- Self-tests GCP evaluates on every create/update; a failing test BLOCKS the change. One test per critical path is the cheapest routing-regression guard on GCP.

**Advanced traffic management** -- The provider's per-route timeout/retry/mirroring/CORS/fault-injection sub-policies are a deliberate, documented coverage boundary of this kind (they overlap the backend service's own resilience settings); weighted splits and rewrites — the routing-defining parts — are fully modeled.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpBackendService** | `defaultService`, path/route rule `service`, `weightedBackendServices[].backendService`, `tests[].service` | `status.outputs.self_link` |
| **GcpBackendBucket** | `defaultService`, path rule `service`, error policy `errorService` | `status.outputs.self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `self_link` | Self-link URI of the URL map | GcpTargetHttpProxy / GcpTargetHttpsProxy `url_map` |
| `url_map_name` | Name as it exists in GCP | Audit, fleet inventory |
| `map_id` | Server-assigned numeric ID | Diagnostics |
| `fingerprint` | Optimistic-concurrency token | Out-of-band gcloud updates |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Host + path fan-out** -- One domain, paths split across API and static backends. Start from the **Host Path Fanout** preset.

**Weighted canary** -- A 90/10 split walked up as confidence grows; weight 0 drains a backend. Start from the **Weighted Canary** preset.

**Apex redirect** -- The redirect-only map a port-80 HTTP proxy serves (http→https 301). Start from the **Apex Redirect** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the URL map is created
- [**GCP Backend Service**](/cloud-catalog/gcp-backend-service) -- the dynamic targets this map routes to
- [**GCP Backend Bucket**](/cloud-catalog/gcp-backend-bucket) -- the static targets and custom-error-page origins
- [**GCP Target HTTPS Proxy**](/cloud-catalog/gcp-target-https-proxy) -- consumes this map's `self_link` behind TLS
- [**GCP Target HTTP Proxy**](/cloud-catalog/gcp-target-http-proxy) -- consumes a redirect-only map for the http→https upgrade
