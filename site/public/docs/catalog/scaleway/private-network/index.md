---
title: "Private Network"
description: "Private Network deployment documentation"
icon: "package"
order: 100
componentName: "scalewayprivatenetwork"
---

# Scaleway Private Network

Deploys a regional Private Network inside a Scaleway VPC with optional explicit IPv4 and IPv6 subnet configuration. Private Networks are the primary attachment point for Kapsule clusters, RDB instances, Redis clusters, Load Balancers, and Instances that need secure, private connectivity. Built-in DHCP assigns IP addresses automatically, with optional CIDR control for predictable address planning. Supports ValueFromRef for VPC dependency wiring in InfraCharts.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Private Network** -- a regional `network.PrivateNetwork` in the specified VPC with built-in DHCP and the configured subnet settings
- **IPv4 Subnet** -- created only when `ipv4Subnet` is specified; a user-defined CIDR block for the network. When omitted, Scaleway's IPAM service auto-allocates a subnet from a default range
- **IPv6 Subnets** -- created only when `ipv6Subnets` entries are provided; dual-stack networking for workloads requiring IPv6 reachability
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway VPC** in the target region. Provide the VPC UUID directly or reference a ScalewayVpc Cloud Resource via ValueFromRef.
- **Non-overlapping CIDR** -- if specifying an explicit `ipv4Subnet`, ensure the range does not overlap with other Private Networks in the same VPC. Overlapping ranges prevent routing between networks.

## Deploy

### Console

Open the deployment store, find **Scaleway Private Network**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Auto-Subnet** preset in the [Presets](#presets) tab to let Scaleway IPAM assign the subnet automatically.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayPrivateNetwork
metadata:
  name: app-network
  org: acme-corp
  env: prod
spec:
  vpcId:
    value: "abc12345-6789-def0-1234-567890abcdef"
  region: fr-par
```

```shell
planton apply -f scaleway-private-network.yaml
```

This creates a Private Network in the Paris region with an IPAM-managed auto-assigned subnet. No explicit CIDR, IPv6, or route propagation is configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Private Network to a VPC deployed in the same InfraPipeline:

```yaml
spec:
  vpcId:
    valueFrom:
      kind: ScalewayVpc
      name: main-vpc
      fieldPath: status.outputs.vpc_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the Private Network with the resolved VPC ID.

## Key Configuration

These are the most important decisions when configuring a Private Network. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Subnet mode** -- Omit `ipv4Subnet` to let Scaleway IPAM auto-allocate a CIDR. Specify an explicit CIDR (e.g., `10.0.1.0/24`) when running multiple Private Networks in a routing-enabled VPC, integrating with VPN tunnels, or requiring documented address ranges. The assigned CIDR is always available in `status.outputs.ipv4_subnet_cidr`.

**Region** -- The `region` must match the parent VPC's region. Cannot be changed after creation.

**Default route propagation** -- Set `enableDefaultRoutePropagation` to true when the parent VPC has routing enabled and this Private Network's resources need to reach resources in other Private Networks within the same VPC. Leave disabled for isolated single-network environments.

**IPv6 subnets** -- Add entries to `ipv6Subnets` for dual-stack networking. Most production deployments use IPv4 only; IPv6 is primarily useful for workloads that must be reachable on IPv6 or for advanced networking scenarios.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayVpc** | `vpcId` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `private_network_id` | Unique identifier (UUID) of the created Private Network | ScalewayKapsuleCluster, ScalewayInstance, ScalewayLoadBalancer, ScalewayRdbInstance network attachment |
| `ipv4_subnet_cidr` | IPv4 CIDR of the Private Network (explicit or IPAM-assigned) | Network planning, firewall rule configuration, VPN tunnel CIDR documentation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Auto-subnet Private Network** -- A Private Network with IPAM-managed automatic subnet allocation. The fastest path to a functional network for attaching instances, Kapsule clusters, and databases without upfront address planning. Start from the **Auto-Subnet** preset.

**Explicit-subnet Private Network** -- A Private Network with a user-defined IPv4 CIDR block for environments requiring predictable addressing, non-overlapping ranges for inter-network routing, or VPN integration. Start from the **Explicit-Subnet** preset.

## Works With

- [**Scaleway VPC**](/cloud-catalog/scaleway-vpc) -- provides the VPC that contains this Private Network