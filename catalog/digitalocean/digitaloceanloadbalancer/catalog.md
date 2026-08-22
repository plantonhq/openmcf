# Load Balancer on DigitalOcean

Deploys a DigitalOcean Load Balancer with regional or global routing, configurable forwarding rules, health checks, VPC placement, backend targeting by Droplet IDs or tags, optional SSL termination, sticky sessions, an LB-level firewall, and connection-tuning knobs. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC, Droplet, and certificate dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Load Balancer** -- a regional or global balancer with the configured forwarding rules or global-routing settings
- **Health Check** -- created only when `healthCheck` is provided; probes backends on the specified port, protocol, and path
- **Sticky Sessions** -- cookie-based affinity when `stickySessions.type` is `cookies`
- **Backend Droplet Attachments** -- targets Droplets via explicit `dropletIds` or a `dropletTag` (mutually exclusive)
- **Firewall** -- source allow/deny rules when `firewall` is set

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A VPC network** in the target region is recommended (and required for `subnetUuid`). Omit `vpc` to use the region's default VPC. GLOBAL balancers take no VPC.
- **Backend Droplets** -- at least one Droplet or a Droplet tag, in the same region (and preferably the same VPC). A memberless balancer is valid but serves no traffic.
- **A TLS certificate** (for HTTPS) -- required when a forwarding rule uses `https` without `tlsPassthrough`. Reference a `DigitalOceanCertificate` by name.

## Deploy

### Console

Open the deployment store, find **Load Balancer on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **HTTPS SSL Termination** preset in the [Presets](#presets) tab for a production HTTPS configuration with tag-based targeting.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanLoadBalancer
metadata:
  name: web-lb
  org: acme-corp
  env: prod
spec:
  loadBalancerName: web-lb
  region: nyc3
  vpc:
    valueFrom:
      kind: DigitalOceanVpc
      name: app-network
      fieldPath: status.outputs.vpc_id
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

This creates an HTTP load balancer forwarding port 80 to port 8080 on all Droplets tagged `web`. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the load balancer to a VPC, Droplets, and a certificate deployed in the same InfraPipeline:

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
  forwardingRules:
    - entryPort: 443
      entryProtocol: https
      targetPort: 80
      targetProtocol: http
      certificateName:
        valueFrom:
          kind: DigitalOceanCertificate
          name: app-cert
          fieldPath: status.outputs.certificate_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC, Droplets, and certificate first, then provisions the load balancer with the resolved values.

## Key Configuration

These are the most important decisions when configuring a load balancer. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type and region** -- `REGIONAL` (the default) and `REGIONAL_NETWORK` require a `region`. `GLOBAL` forbids a region and uses `glbSettings`, `domains`, and `targetLoadBalancerIds` instead of `forwardingRules`.

**Sizing** -- `size` (`lb-small` / `lb-medium` / `lb-large`) or `sizeUnit` (1–200). They are mutually exclusive; unset means `lb-small`.

**Forwarding rules** -- Each rule maps an `entryPort`/`entryProtocol` pair to a `targetPort`/`targetProtocol`. For HTTPS termination, set `entryProtocol: https` and provide `certificateName`. For passthrough, set `tlsPassthrough: true` and omit the certificate.

**Backend targeting** -- Use `dropletTag` for dynamic membership or `dropletIds` for an explicit list. The two options are mutually exclusive.

**Health checks** -- Configure `healthCheck` with port, protocol, path (required for http/https, forbidden for tcp), and optional thresholds. Without a health check, DigitalOcean applies a TCP check against the first forwarding rule's target port.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** (optional) | `vpc` | `status.outputs.vpc_id` |
| **DigitalOceanDroplet** (optional) | `dropletIds` | `status.outputs.droplet_id` |
| **DigitalOceanCertificate** (optional) | `forwardingRules[].certificateName`, `domains[].certificateName` | `status.outputs.certificate_id` (the certificate NAME) |
| **DigitalOceanLoadBalancer** (optional) | `targetLoadBalancerIds` | `status.outputs.load_balancer_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | Unique load balancer identifier (UUID) | DigitalOcean API operations, GLOBAL balancer targets |
| `ip` | Public IPv4 address assigned to the load balancer | DNS A records, client-facing endpoint configuration |
| `urn` | Uniform resource name (`do:loadbalancer:<id>`) | DigitalOcean project resources |
| `ipv6` | IPv6 address when `networkStack` is `DUALSTACK` | DNS AAAA records |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTPS with SSL termination** -- TLS terminates on port 443 at the load balancer, forwarding to backends over HTTP on port 80. Uses tag-based Droplet targeting and HTTP health checks. Start from the **HTTPS SSL Termination** preset.

**Simple HTTP load balancer** -- HTTP on port 80 forwarded to port 8080 on explicitly listed Droplets. No TLS. Suitable for development, staging, or internal services behind a CDN. Start from the **HTTP Basic** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the VPC network for load balancer placement
- [**DigitalOcean Droplet**](/cloud-catalog/digital-ocean-droplet) -- provides backend compute instances for traffic routing
- [**DigitalOcean Certificate**](/cloud-catalog/digital-ocean-certificate) -- provides the TLS certificate for HTTPS termination
