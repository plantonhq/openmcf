# Hetzner Cloud Network

Deploys a private network with subnets and optional static routes on Hetzner Cloud. A network provides isolated IPv4 connectivity between servers, load balancers, and other Hetzner Cloud resources using RFC 1918 private address space. Subnets carve the network into smaller blocks assigned to specific network zones, enabling multi-region private connectivity. At least one subnet is required because Hetzner Cloud resources attach to subnets, not directly to the network.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Network** -- an `hcloud_network` resource with the specified private IP range
- **Subnets** -- one or more `hcloud_network_subnet` resources, each assigned to a network zone with a CIDR block carved from the network's IP range
- **Routes** (optional) -- `hcloud_network_route` resources defining static routing for VPN gateways, NAT instances, or inter-network traffic

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **IP range planning** -- choose a private CIDR block (10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16) and plan subnet allocations within it.
- **Network zone selection** -- decide which zones (eu-central, us-east, us-west, ap-southeast) your subnets will span.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Network**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including IP range, subnets, and routes.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudNetwork
metadata:
  name: main-network
  org: acme-corp
  env: prod
spec:
  ipRange: "10.0.0.0/16"
  subnets:
    - type: cloud
      networkZone: eu-central
      ipRange: "10.0.1.0/24"
    - type: cloud
      networkZone: eu-central
      ipRange: "10.0.2.0/24"
```

```shell
planton apply -f hetznercloud-network.yaml
```

This creates a /16 network with two /24 subnets in the eu-central zone. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a server environment, servers and load balancers reference this network via ValueFromRef:

```yaml
# In the HetznerCloudServer manifest:
spec:
  networkId:
    valueFrom:
      kind: HetznerCloudNetwork
      name: main-network
      fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, creates the network and subnets first, then provisions servers attached to the network.

## Key Configuration

These are the most important decisions when configuring a network. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**IP range** -- The `ipRange` field sets the top-level CIDR block. Must be from RFC 1918 private space (10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16). All subnet IP ranges must fall within this block. Changing this value forces replacement of the network.

**Subnets** -- The `subnets` field defines one or more subdivisions of the network. Each subnet specifies a `type` (cloud for standard servers, server for dedicated Robot servers, vswitch for hybrid connectivity), a `networkZone`, and an `ipRange` that must not overlap with other subnets. At least one subnet is required.

**Routes** -- The `routes` field adds optional static routes. Each route specifies a `destination` CIDR and a `gateway` IP within one of the network's subnets. Useful for VPN gateways, NAT instances, or routing between networks.

**Delete protection** -- The `deleteProtection` field prevents accidental deletion via the Hetzner Cloud API.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `network_id` | Hetzner Cloud numeric ID of the network | HetznerCloudServer `networkId`, HetznerCloudLoadBalancer `networkId` |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single-zone network** -- A /16 network with one /24 subnet in eu-central. Suitable for single-region deployments where all servers share one private network.

**Multi-zone network** -- A /16 network with subnets in multiple zones for cross-region private connectivity between servers.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers attach to this network via `networkId`
- [**Hetzner Cloud Load Balancer**](/cloud-catalog/hetznercloud-load-balancer) -- load balancers connect to this network for private backend communication
