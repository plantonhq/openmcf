# Zero Trust Tunnel Route on Cloudflare

Provisions a Cloudflare Tunnel route: it advertises a private IP range (CIDR) as reachable through a specific tunnel, within a virtual network. WARP clients and other tunnels can then reach that range. A route has an independent lifecycle from the tunnel -- you add or remove reachable networks without touching the tunnel -- and a tunnel commonly carries many routes. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Tunnel Route** -- a CIDR advertised through a tunnel within a virtual network
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Cloudflare Tunnel edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A tunnel** -- the route advertises a network through an existing `CloudflareZeroTrustTunnel`.

## Deploy

### Console

Open the deployment store, find **Zero Trust Tunnel Route on Cloudflare**, and click **Deploy**. The creation wizard captures the route (account + CIDR + comment) and its targets (the tunnel that serves it and, optionally, the virtual network it belongs to), with a live connection diagram.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareZeroTrustTunnelRoute
metadata:
  name: prod-subnet
  org: acme-corp
  env: prod
spec:
  accountId: a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4
  network: 10.0.0.0/24
  tunnelId:
    valueFrom:
      kind: CloudflareZeroTrustTunnel
      name: prod-connector
      fieldPath: status.outputs.tunnel_id
```

```shell
planton apply -f cloudflare-zero-trust-tunnel-route.yaml
```

This makes the `10.0.0.0/24` subnet reachable to WARP clients through the `prod-connector` tunnel. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a route. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Account (`accountId`)** -- The account that owns the route. Immutable -- changing it replaces the route.

**Network (`network`)** -- The private IPv4 or IPv6 range advertised by this route, in CIDR notation. A single host uses `/32` (IPv4) or `/128` (IPv6).

**Tunnel (`tunnelId`)** -- The tunnel that serves this network. Reference a `CloudflareZeroTrustTunnel` to keep the dependency in the graph.

**Virtual Network (`virtualNetworkId`)** -- Optional. Reference a `CloudflareZeroTrustTunnelVirtualNetwork`, or omit to use the account's default. Use distinct virtual networks to advertise overlapping CIDRs through different tunnels without collision.

## Outputs and Dependencies

### What This Component Consumes

The route references a **CloudflareZeroTrustTunnel** (via `tunnelId`) and, optionally, a **CloudflareZeroTrustTunnelVirtualNetwork** (via `virtualNetworkId`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `route_id` | The Cloudflare-assigned UUID of the route | Verification, dashboards |
| `network` | The advertised CIDR | Auditing, grouping |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Private subnet** -- advertise a `/24` through a tunnel so WARP clients reach internal hosts.

**Isolated overlap** -- place two overlapping CIDRs in distinct virtual networks so they route through different tunnels without collision.

## Works With

- [**Zero Trust Tunnel on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel) -- the tunnel that serves this route's network
- [**Zero Trust Tunnel Virtual Network on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel-virtual-network) -- the routing segment this route belongs to
