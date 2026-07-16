# GCP Target HTTP Proxy

Deploys a global Compute Engine target HTTP proxy (`google_compute_target_http_proxy`) — the plaintext-HTTP frontend adapter of a global external Application Load Balancer. It binds a global forwarding rule (the VIP) to a URL map (the routing brain); the standard production role is serving the http→https redirect on port 80 while the target HTTPS proxy serves the application on 443.

## What Gets Created

A single global target HTTP proxy in the chosen project. Global forwarding rules reference its `self_link`; the proxy references a URL map.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **A URL map** — a `GcpUrlMap` (or its self-link) for the proxy to route through
- **IAM permissions** — any role carrying `compute.targetHttpProxies.*` on the target project

## Quick Start

Create a file `http-proxy.yaml`:

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

Deploy:

```shell
planton apply -f http-proxy.yaml
```

This creates a proxy ready for a port-80 global forwarding rule to bind.

## Configuration Reference

| Field | Description |
|-------|-------------|
| `projectId` | Project owning the proxy (literal or `GcpProject` reference); empty uses the provider's default project |
| `proxyName` | Cloud-side name (RFC1035); defaults to `metadata.name`. Immutable |
| `description` | What this proxy fronts. Immutable |
| `urlMap` | The URL map to route through (required; `GcpUrlMap` reference or self-link). **Mutable in place** — repointing a live frontend causes no downtime |
| `httpKeepAliveTimeoutSec` | Idle client keep-alive, 5-1200s; only honored by `EXTERNAL_MANAGED` load balancers (GCP default 610). Immutable |
| `proxyBind` | Bind to Traffic Director mesh VIPs instead of Google's edge (`INTERNAL_SELF_MANAGED` only). Immutable |

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `self_link` | `string` | Self-link URI — the value a global forwarding rule references as its target |
| `proxy_name` | `string` | Name of the proxy in GCP |
| `proxy_id` | `string` | Server-assigned numeric ID |
| `fingerprint` | `string` | Fingerprint for optimistic concurrency control |

## Deployment Methods

### Pulumi (Go)

See [`iac/pulumi/README.md`](iac/pulumi/README.md).

### Terraform

See [`iac/tf/README.md`](iac/tf/README.md).

## Important Notes

- **Only `urlMap` is mutable** — every other field forces destroy-and-recreate, briefly breaking any forwarding rule referencing the old `self_link`.
- **The pair pattern**: point this proxy at a redirect-only URL map (`defaultUrlRedirect` with `httpsRedirect: true`) and let a `GcpTargetHttpsProxy` serve the real application — two forwarding rules share one static IP on ports 80 and 443.
- **TLS never lives here** — for HTTPS termination use `GcpTargetHttpsProxy`.

## Related Components

- [GcpUrlMap](/docs/catalog/gcp/gcpurlmap) — the routing table this proxy consults
- [GcpGlobalForwardingRule](/docs/catalog/gcp/gcpglobalforwardingrule) — the VIP that binds to this proxy
- [GcpTargetHttpsProxy](/docs/catalog/gcp/gcptargethttpsproxy) — the TLS-terminating sibling
- [GcpProject](/docs/catalog/gcp/gcpproject) — provides the GCP project that owns the proxy

## Additional Resources

- [Target proxies overview](https://cloud.google.com/load-balancing/docs/target-proxies)
- [Setting up HTTP-to-HTTPS redirect](https://cloud.google.com/load-balancing/docs/https/setting-up-http-https-redirect)

## Support

For issues, questions, or contributions, please refer to the Planton documentation or open an issue in the repository.
