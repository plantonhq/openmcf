# Target HTTP Proxy on Google Cloud

Deploys a global Compute Engine target HTTP proxy — the plaintext-HTTP frontend adapter of a global external Application Load Balancer. The proxy binds a global forwarding rule (the VIP) to a URL map (the routing brain): the rule delivers client connections, the proxy consults the map for every request. It is deliberately thin — TLS lives on the target HTTPS proxy, routing on the URL map, traffic policy on the backend service. The standard production pattern is a PAIR sharing one static IP: this proxy serves a redirect-only URL map (http→https 301) while the HTTPS proxy serves the application. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to projects and URL maps.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Engine Target HTTP Proxy (global)** -- bound to the configured URL map, with optional keep-alive tuning and Traffic Director bind

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the proxy will be created.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.
- **A URL map** (GcpUrlMap) — the proxy's one required dependency; for the redirect pattern, a redirect-only map.

## Deploy

### Console

Open the deployment store, find **Target HTTP Proxy on Google Cloud**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTPS Redirect Frontend** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpTargetHttpProxy
metadata:
  name: web-http-proxy
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  urlMap:
    value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/urlMaps/http-redirect-map"
```

```shell
planton apply -f target-http-proxy.yaml
```

This creates the redirect half: a port-80 forwarding rule pointing here upgrades every request to HTTPS.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the URL map:

```yaml
spec:
  urlMap:
    valueFrom:
      kind: GcpUrlMap
      name: http-redirect-map
      fieldPath: status.outputs.self_link
```

The InfraPipeline resolves the dependency graph — the URL map first, then this proxy — and a downstream GcpGlobalForwardingRule references this proxy's `self_link` as its target.

## Key Configuration

These are the most important decisions when configuring a target HTTP proxy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**URL map** -- The one REQUIRED field, and the ONLY mutable one: GCP swaps it in place (a dedicated setUrlMap call), so repointing a live frontend at a new routing table causes zero downtime — the blue/green lever for whole routing schemes.

**Keep-alive timeout** -- 5-1200 seconds; only honored by the envoy-based EXTERNAL_MANAGED scheme (GCP default 610s). Raise it above your clients' own keep-alive so the LB never closes first. Immutable.

**Traffic Director bind** -- `proxyBind` attaches the proxy to the mesh's private IPs instead of Google's edge; only meaningful behind an INTERNAL_SELF_MANAGED forwarding rule. Immutable.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpUrlMap** | `urlMap` | `status.outputs.self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `self_link` | Self-link URI of the proxy | GcpGlobalForwardingRule `target` |
| `proxy_name` | Name as it exists in GCP | Audit, fleet inventory |
| `proxy_id` | Server-assigned numeric ID | Diagnostics |
| `fingerprint` | Optimistic-concurrency token | Out-of-band gcloud updates |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS redirect frontend** -- This proxy + a redirect-only URL map: the port-80 half of every production frontend. Start from the **HTTPS Redirect Frontend** preset.

**Plain HTTP frontend** -- Serving an application over plain HTTP (internal tools, pre-TLS testing). Start from the **Plain HTTP Frontend** preset.

**Traffic Director mesh** -- The proxy bound to mesh-private IPs for INTERNAL_SELF_MANAGED frontends. Start from the **Traffic Director Mesh** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the proxy is created
- [**GCP URL Map**](/cloud-catalog/gcp-url-map) -- the routing table this proxy consults
- [**GCP Global Forwarding Rule**](/cloud-catalog/gcp-global-forwarding-rule) -- consumes this proxy's `self_link` as its target
- [**GCP Target HTTPS Proxy**](/cloud-catalog/gcp-target-https-proxy) -- the TLS sibling serving the application half
