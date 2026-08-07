---
title: "Load Balancer"
description: "Load Balancer deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareloadbalancer"
---

# Load Balancer on Cloudflare

Deploys a Cloudflare Load Balancer with an origin pool, health monitor, and configurable traffic steering. Integrates with Planton's Provider Connections for Cloudflare credential management and supports ValueFromRef wiring to DNS zones for cross-resource dependency resolution in InfraPipelines.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Load Balancer Monitor** -- an HTTP health check that probes each origin at the configured `healthProbePath`, expecting `2xx` responses, with 2 retries and a 5-second timeout
- **Load Balancer Pool** -- a pool named `{metadata.name}-pool` containing all declared origins with their respective weights, linked to the health monitor
- **Load Balancer** -- the load balancer resource bound to the specified hostname and DNS zone, using the pool as both the default and fallback, with the configured proxy, steering policy, and session affinity settings
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Load Balancing edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **An existing Cloudflare DNS zone** -- the hostname must resolve within a zone you control. Provide the zone ID directly or reference a CloudflareDnsZone resource via ValueFromRef.
- **Load Balancing add-on** -- Cloudflare Load Balancing is a paid feature and must be enabled on your account before deploying.
- **Reachable origins** -- each origin address (IP or hostname) must be accessible from the internet over HTTP(S) for health checks and traffic routing.

## Deploy

### Console

Open the deployment store, find **Load Balancer on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Active-Passive Failover** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareLoadBalancer
metadata:
  name: api-lb
  org: acme-corp
  env: prod
spec:
  hostname: api.example.com
  zoneId:
    value: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  origins:
    - name: primary
      address: 203.0.113.10
      weight: 1
  proxied: true
  healthProbePath: /healthz
```

```shell
planton apply -f cloudflare-load-balancer.yaml
```

This creates a proxied load balancer for `api.example.com` with a single origin health-checked at `/healthz`. No session affinity or geo-steering is configured -- traffic uses static failover by default. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying alongside the DNS zone, use ValueFromRef to wire the load balancer to a zone deployed in the same InfraPipeline:

```yaml
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: prod-zone
      fieldPath: status.outputs.zone_id
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the load balancer with the resolved zone ID.

## Key Configuration

These are the most important decisions when configuring a load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Steering policy** -- Controls how traffic is distributed across origins. `off` (default) uses static failover where origins are tried in declaration order. `geo` routes clients to the geographically nearest healthy origin. `random` distributes traffic by weight for A/B testing or canary deployments.

**Session affinity** -- Set `sessionAffinity` to `cookie` to pin returning clients to the same origin using a Cloudflare-managed cookie. Leave as `none` (default) for stateless workloads where any origin can serve any request.

**Health probe path** -- The `healthProbePath` (default `/`) defines the HTTP GET endpoint used by the monitor. Choose a lightweight endpoint that accurately reflects origin health. The monitor expects `2xx` responses and retries twice with a 5-second timeout.

**Origin weights** -- Each origin's `weight` controls its relative traffic share when using `random` steering. For failover (`off`), weights are ignored and origins are tried in order.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareDnsZone** | `zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | The unique identifier of the created Cloudflare Load Balancer | Monitoring dashboards, management automation |
| `load_balancer_dns_record_name` | The hostname DNS record associated with the load balancer | DNS verification, application endpoint configuration |
| `load_balancer_cname_target` | The canonical CNAME target that the hostname resolves to (Cloudflare endpoint) | External DNS configuration, CNAME record setup |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Active-passive failover** -- Two origins with static steering. All traffic goes to the first healthy origin; if it fails health checks, traffic fails over to the second. Use for primary/backup server architectures and DR failover. Start from the **Active-Passive Failover** preset.

**Geographic routing** -- Multiple origins with `geo` steering. Cloudflare routes clients to the geographically nearest healthy origin for lowest latency. Use for multi-region deployments serving US, EU, and APAC traffic. Start from the **Geographic Routing** preset.

**Weighted A/B testing** -- Two or more origins with `random` steering and different weights. Traffic is distributed proportionally (e.g., 70/30 split) for A/B tests, canary deployments, or gradual rollouts. Start from the **Weighted A/B Testing** preset.

## Works With

- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- provides the zone ID that determines which domain the load balancer hostname belongs to