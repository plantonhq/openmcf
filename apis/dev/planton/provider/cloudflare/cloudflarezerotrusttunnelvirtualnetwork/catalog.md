# Zero Trust Tunnel Virtual Network on Cloudflare

Provisions a Cloudflare Tunnel virtual network: an isolated routing segment that lets the same private CIDR (for example `10.0.0.0/8`) be connected through more than one tunnel without collision. Routes (`CloudflareZeroTrustTunnelRoute`) attach a private network to a tunnel within one virtual network, and WARP clients select which virtual network to reach. A virtual network is account-scoped and outlives any individual tunnel. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Virtual Network** -- a named, account-scoped routing segment
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Cloudflare Tunnel edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Zero Trust enabled** -- the account must have Cloudflare Zero Trust (Cloudflare One) set up.

## Deploy

### Console

Open the deployment store, find **Zero Trust Tunnel Virtual Network on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account, a name, an optional comment, and whether this virtual network is the account default.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareZeroTrustTunnelVirtualNetwork
metadata:
  name: prod-overlay
  org: acme-corp
  env: prod
spec:
  accountId: a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4
  name: prod-overlay
  comment: prod data center segment
  isDefaultNetwork: false
```

```shell
planton apply -f cloudflare-zero-trust-tunnel-virtual-network.yaml
```

This creates a named routing segment for the prod data center. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a virtual network. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Account (`accountId`)** -- The account that owns the virtual network. Immutable -- changing it replaces the virtual network.

**Name (`name`)** -- A human-readable name shown in the Zero Trust dashboard and used to disambiguate overlapping routes.

**Account Default (`isDefaultNetwork`)** -- When on, routes and WARP clients that do not name a virtual network use this one. Exactly one virtual network can be the default at a time -- promoting this one demotes the previous default.

## Outputs and Dependencies

### What This Component Consumes

Nothing -- a virtual network is a self-contained, account-scoped leaf.

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `virtual_network_id` | The Cloudflare-assigned UUID | Referenced by `CloudflareZeroTrustTunnelRoute` |
| `virtual_network_name` | The virtual network name | Auditing, grouping |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Isolated segment** -- a named virtual network keeps a site's overlapping CIDRs separate from other sites.

**Default network** -- mark one virtual network as the account default for the common single-overlay case.

## Works With

- [**Zero Trust Tunnel Route on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel-route) -- routes advertise CIDRs within this virtual network
- [**Zero Trust Tunnel on Cloudflare**](/cloud-catalog/cloudflare-zero-trust-tunnel) -- the tunnel a route binds a network to
