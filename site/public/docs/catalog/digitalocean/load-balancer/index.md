---
title: "Load Balancer"
description: "Load Balancer deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceanloadbalancer"
---

# Load Balancer on DigitalOcean

Deploys a DigitalOcean Load Balancer with configurable forwarding rules, health checks, VPC placement, backend targeting by Droplet IDs or tags, optional SSL termination, and sticky sessions. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC and Droplet dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Load Balancer** -- a regional load balancer in the specified VPC with the configured forwarding rules routing traffic from entry ports to backend Droplet target ports
- **Health Check** -- created only when `healthCheck` is provided; probes backend Droplets on the specified port, protocol, and path at the configured interval
- **Sticky Sessions** -- configured only when `enableStickySessions` is true; uses cookie-based session affinity to route repeated requests from the same client to the same backend
- **Backend Droplet Attachments** -- targets Droplets via explicit `dropletIds` or a `dropletTag` (mutually exclusive); tag-based targeting automatically includes/excludes Droplets as they are created/destroyed

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A VPC network** in the target region. The load balancer and its backend Droplets must be in the same VPC. Provide the VPC UUID directly or reference a DigitalOceanVpc Cloud Resource via ValueFromRef.
- **Backend Droplets** -- at least one Droplet or a Droplet tag with matching instances. Droplets must be in the same region and VPC as the load balancer.
- **A TLS certificate** (for HTTPS) -- required when a forwarding rule uses `https` as the entry protocol. Upload the certificate to DigitalOcean and reference it by name in `certificateName`.

## Deploy

### Console

Open the deployment store, find **Load Balancer on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTPS SSL Termination** preset in the [Presets](#presets) tab for a production HTTPS configuration with tag-based targeting.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  loadBalancerName: web-lb
  region: nyc3
  vpc:
    value: "abc12345-6789-def0-1234-567890abcdef"
  forwardingRules:
    - entryPort: 80
      entryProtocol: http
      targetPort: 8080
      targetProtocol: http
  dropletTag: web
```

```shell
planton apply -f do-load-balancer.yaml
```

This creates an HTTP load balancer forwarding port 80 to port 8080 on all Droplets tagged `web` in the VPC. No health check, SSL termination, or sticky sessions are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the load balancer to a VPC and Droplets deployed in the same InfraPipeline:

```yaml
spec:
  vpc:
    valueFrom:
      kind: DigitalOceanVpc
      name: app-network
      fieldPath: status.outputs.vpc_id
  dropletIds:
    - valueFrom:
        kind: DigitalOceanDroplet
        name: web-server-1
        fieldPath: status.outputs.droplet_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC and Droplets first, then provisions the load balancer with the resolved values.

## Key Configuration

These are the most important decisions when configuring a load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Forwarding rules** -- Each rule maps an `entryPort`/`entryProtocol` pair on the load balancer to a `targetPort`/`targetProtocol` on the backend. For HTTPS termination, set `entryProtocol: https` and provide `certificateName` with the name of a TLS certificate uploaded to DigitalOcean.

**Backend targeting** -- Use `dropletTag` for dynamic, tag-based membership where Droplets are added/removed automatically, or `dropletIds` for explicit control over which Droplets receive traffic. The two options are mutually exclusive.

**Health checks** -- Configure `healthCheck` with the port, protocol, optional path, and check interval. Unhealthy Droplets are removed from the rotation. Without a health check, the load balancer forwards to all attached Droplets regardless of health.

**Sticky sessions** -- Set `enableStickySessions: true` for cookie-based session affinity. Required for applications that store session state on individual backend servers rather than in a shared store.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** | `vpc` | `status.outputs.vpc_id` |
| **DigitalOceanDroplet** (optional) | `dropletIds` | `status.outputs.droplet_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | Unique load balancer identifier (UUID) | DigitalOcean API operations, firewall rules |
| `ip` | Public IP address assigned to the load balancer | DNS records, client-facing endpoint configuration |
| `dns_name` | DNS name for the load balancer | DNS CNAME records |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS with SSL termination** -- TLS terminates on port 443 at the load balancer, forwarding to backends over HTTP on port 80. Uses tag-based Droplet targeting and HTTP health checks. Start from the **HTTPS SSL Termination** preset.

**Simple HTTP load balancer** -- HTTP on port 80 forwarded to port 8080 on explicitly listed Droplets. No TLS. Suitable for development, staging, or internal services behind a CDN or reverse proxy. Start from the **HTTP Basic** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the VPC network for load balancer placement
- [**DigitalOcean Droplet**](/cloud-catalog/digital-ocean-droplet) -- provides backend compute instances for traffic routing