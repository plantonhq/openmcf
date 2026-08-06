---
title: "URL Map"
description: "URL Map deployment documentation"
icon: "package"
order: 100
componentName: "gcpurlmap"
---

# GCP URL Map

Creates a global Compute Engine URL map — the routing brain that matches each request's host and path and decides what happens next (forward, split, rewrite, redirect, or custom error page). Target HTTP(S) proxies attach to its self-link; forwarding rules and addresses sit in front of the proxy.

## What Gets Created

A single `google_compute_url_map` in the chosen GCP project.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A backend target** for the default route — typically a `GcpBackendService` self-link
- **IAM permissions** — any role carrying `compute.urlMaps.*` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpUrlMap
metadata:
  name: web-routing
spec:
  projectId:
    value: my-gcp-project-123
  defaultService:
    value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/backendServices/web-backend
```

```shell
planton apply -f url-map.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `defaultService` / `defaultUrlRedirect` / `defaultRouteAction` | ref / object / object | — | Required, exactly one top-level default target |
| `hostRules` | list | `[]` | Map Host headers to named path matchers |
| `pathMatchers` | list | `[]` | Path-level routing (`pathRules` or `routeRules`, not both) |
| `headerAction` | object | — | URL-map-level header mutations |
| `tests` | list | `[]` | Routing self-tests evaluated at update time |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the map. Immutable. |
| `urlMapName` | `string` | `metadata.name` | Cloud-side name (RFC1035). Immutable. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `self_link` | Self-link URI — the value target proxies reference |
| `url_map_name` | Name of the URL map in GCP |
| `map_id` | Server-assigned numeric ID |
| `fingerprint` | Fingerprint for optimistic concurrency control |

## Related Components

- [GcpBackendService](/docs/catalog/gcp/backend-service) — backends this map routes to
- [GcpBackendBucket](/docs/catalog/gcp/backend-bucket) — static assets and custom error pages
- [GcpProject](/docs/catalog/gcp/project) — provides the project that owns the map
