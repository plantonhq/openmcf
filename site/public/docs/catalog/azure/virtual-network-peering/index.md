---
title: "Virtual Network Peering"
description: "Virtual Network Peering deployment documentation"
icon: "package"
order: 100
componentName: "azurevirtualnetworkpeering"
---

# Azure Virtual Network Peering

Creates one direction of an Azure virtual network peering -- private connectivity between two virtual networks over the Microsoft backbone. The building block for hub-and-spoke topologies, shared services networks, and cross-subscription connectivity.

## What Gets Created

When you deploy an AzureVirtualNetworkPeering resource, Planton provisions:

- **Virtual Network Peering** — an `azurerm_virtual_network_peering` written on the local network, pointing at the remote network

One resource is one direction. Traffic flows only after the reciprocal peering exists on the remote network -- typically stamped as a second resource with local and remote swapped.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Two virtual networks** to connect (each an `AzureVirtualNetwork` in composed environments)
- **Network write rights**: `Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write` (Network Contributor, Contributor, or Owner) on both networks

## Quick Start

Create a file `peering.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetworkPeering
metadata:
  name: hub-to-spoke1
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetworkPeering.hub-to-spoke1
spec:
  name: hub-to-spoke1
  virtualNetworkId:
    valueFrom:
      name: hub-vnet
  remoteVirtualNetworkId:
    valueFrom:
      name: spoke1-vnet
  allowForwardedTraffic: true
```

Deploy:

```shell
planton apply -f peering.yaml
```

Deploy the reciprocal direction (`spoke1-to-hub`) before expecting bidirectional traffic.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | Peering name within the local network. Name it after the far side so the network's peering list reads as a topology map. | Required, 1-80 chars, Azure naming rules |
| `virtualNetworkId` | `StringValueOrRef` | The local network's ARM ID -- the side this peering is written on. Defaults to referencing an `AzureVirtualNetwork`'s ID output. | Required |
| `remoteVirtualNetworkId` | `StringValueOrRef` | The remote network's ARM ID. Works across subscriptions and regions unchanged. Defaults to referencing an `AzureVirtualNetwork`'s ID output. | Required |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `allowVirtualNetworkAccess` | `bool` | Whether traffic from the local network can reach the remote network. Azure defaults to `true`. Set `false` to keep the peering established but block direct VM-to-VM traffic in this direction. |
| `allowForwardedTraffic` | `bool` | Whether traffic forwarded by the remote network (e.g. relayed by an NVA or VPN) is accepted. Azure defaults to `false`. Hub-to-spoke peerings in inspection topologies typically set `true`. |
| `allowGatewayTransit` | `bool` | Whether the local network's VPN/ExpressRoute gateway may be used by the remote network. Azure defaults to `false`. Set `true` on the hub side when spokes ride the hub gateway. |
| `useRemoteGateways` | `bool` | Whether the local network uses the remote network's gateway. Azure defaults to `false`. Set `true` on the spoke side; requires `allowGatewayTransit: true` on the reciprocal peering. Only one peering per network may set this; incompatible with global (cross-region) peering. |
| `peerCompleteVirtualNetworksEnabled` | `bool` | Whether the peering spans complete address spaces. Azure defaults to `true`. Set `false` for subnet-scoped peering with `localSubnetNames` / `remoteSubnetNames`. |
| `localSubnetNames` | `list(string)` | Local subnets included when subnet-scoped peering is enabled. |
| `remoteSubnetNames` | `list(string)` | Remote subnets included when subnet-scoped peering is enabled. |
| `onlyIpv6PeeringEnabled` | `bool` | Whether only IPv6 address space is peered (subnet-scoped dual-stack). Azure defaults to `false`. |

## Examples

### Hub-to-Spoke Pair

Two resources -- one per direction:

```yaml
# Written on the hub
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetworkPeering
metadata:
  name: hub-to-spoke1
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetworkPeering.hub-to-spoke1
spec:
  name: hub-to-spoke1
  virtualNetworkId:
    valueFrom:
      name: hub-vnet
  remoteVirtualNetworkId:
    valueFrom:
      name: spoke1-vnet
  allowForwardedTraffic: true
  allowGatewayTransit: true
---
# Written on the spoke
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetworkPeering
metadata:
  name: spoke1-to-hub
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetworkPeering.spoke1-to-hub
spec:
  name: spoke1-to-hub
  virtualNetworkId:
    valueFrom:
      name: spoke1-vnet
  remoteVirtualNetworkId:
    valueFrom:
      name: hub-vnet
  useRemoteGateways: true
```

### Subnet-Scoped Peering

Peer only named subnets instead of full address spaces:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetworkPeering
metadata:
  name: shared-to-app
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureVirtualNetworkPeering.shared-to-app
spec:
  name: shared-to-app
  virtualNetworkId:
    valueFrom:
      name: shared-services-vnet
  remoteVirtualNetworkId:
    valueFrom:
      name: app-vnet
  peerCompleteVirtualNetworksEnabled: false
  localSubnetNames:
    - shared-services
  remoteSubnetNames:
    - app
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `peering_id` | `string` | Full ARM ID of the peering |
| `peering_name` | `string` | The peering's name within its local network |
| `virtual_network_name` | `string` | The local network's name, derived from `virtual_network_id` |
| `resource_group_name` | `string` | The local network's resource group, derived from `virtual_network_id` |

## Related Components

- [AzureVirtualNetwork](/docs/catalog/azure/virtual-network) — the networks this peering connects (prerequisite for both local and remote sides)
- [AzureSubnet](/docs/catalog/azure/subnet) — subnet-scoped peering references subnets by name within those networks
- [AzureRouteTable](/docs/catalog/azure/route-table) — user-defined routing that steers traffic across peered networks
- [AzureNetworkSecurityGroup](/docs/catalog/azure/network-security-group) — traffic filtering, complementary to peering connectivity
