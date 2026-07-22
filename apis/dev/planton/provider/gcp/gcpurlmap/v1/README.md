# GCP URL Map

Deploys a global Compute Engine URL map (`google_compute_url_map`) — the L7 routing brain of a global external Application Load Balancer. It matches each request's host and path and decides whether to forward to a backend, split traffic across weighted backends, rewrite or redirect the URL, inject headers, or serve custom error pages.

## What Gets Created

A single global URL map in the chosen project. Target HTTP(S) proxies reference its `self_link`; forwarding rules and addresses sit in front of the proxy.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **At least one backend target** — typically a `GcpBackendService` or `GcpBackendBucket` self-link for the default route
- **IAM permissions** — any role carrying `compute.urlMaps.*` on the target project

## Quick Start

Create a file `url-map.yaml`:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpUrlMap
metadata:
  name: web-routing
spec:
  projectId:
    value: my-gcp-project-123
  defaultService:
    value: https://www.googleapis.com/compute/v1/projects/my-gcp-project-123/global/backendServices/web-backend
```

Deploy:

```shell
planton apply -f url-map.yaml
```

This creates a URL map that sends all unmatched traffic to the named backend service.

## Configuration Reference

### Default target (exactly one required)

| Field | Description |
|-------|-------------|
| `defaultService` | Backend service or bucket self-link for unmatched traffic |
| `defaultUrlRedirect` | Redirect unmatched traffic (apex→www, http→https) |
| `defaultRouteAction` | Weighted split across backends (requires `weightedBackendServices`) |

### Routing table

| Field | Description |
|-------|-------------|
| `hostRules` | Map Host headers (with optional wildcards) to named path matchers |
| `pathMatchers` | Named path-level routing: `pathRules` (longest prefix) **or** `routeRules` (priority-ordered rich matching), plus an optional per-matcher default |
| `headerAction` | Headers added/removed at the URL-map level before per-route actions |
| `defaultCustomErrorResponsePolicy` | Custom error pages from a backend bucket (global external ALBs only) |
| `tests` | Routing self-tests GCP evaluates at create/update time |

### Route actions

Inside `defaultRouteAction`, `pathRules[].routeAction`, and `routeRules[].routeAction`:

| Sub-field | Description |
|-----------|-------------|
| `weightedBackendServices` | Relative weights splitting traffic across backend services |
| `urlRewrite` | `hostRewrite`, `pathPrefixRewrite`, or (route rules only) `pathTemplateRewrite` |

Timeout, retry, CORS, fault injection, and request mirroring inside `routeAction` are intentionally omitted — a documented coverage boundary; adding them later is purely additive.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `self_link` | `string` | Self-link URI — the value target proxies reference |
| `url_map_name` | `string` | Name of the URL map in GCP |
| `map_id` | `string` | Server-assigned numeric ID |
| `fingerprint` | `string` | Fingerprint for optimistic concurrency control |

## Deployment Methods

### Pulumi (Go)

See [`iac/pulumi/README.md`](iac/pulumi/README.md).

### Terraform

See [`iac/tf/README.md`](iac/tf/README.md).

## Important Notes

- **Immutability**: `urlMapName` and `projectId` are ForceNew — changing either destroys and recreates the map, briefly breaking every target proxy referencing the old `self_link`.
- **Path matcher exclusivity**: each path matcher uses either `pathRules` or `routeRules`, not both.
- **Evaluation order**: host rules → route rules (by priority) → path rules (longest prefix) → path matcher default → URL map default.
- **Regional URL maps** are a separate GCP resource reserved for a future kind.

## Related Components

- [GcpBackendService](/docs/catalog/gcp/gcpbackendservice) — the backends this map routes to
- [GcpBackendBucket](/docs/catalog/gcp/gcpbackendbucket) — static assets and custom error pages
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project that owns the map

## Additional Resources

- [URL map concepts](https://cloud.google.com/load-balancing/docs/url-map-concepts)
- [Routing rules and traffic management](https://cloud.google.com/load-balancing/docs/https/traffic-management)

## Support

For issues, questions, or contributions, please refer to the Planton documentation or open an issue in the repository.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
