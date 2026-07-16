# AzureVirtualNetwork

## Overview

`AzureVirtualNetwork` provisions an Azure Virtual Network (VNet): the
isolated, private IP address space every network-attached Azure workload
lives inside. It models the full network surface -- multi-CIDR (and
dual-stack) address spaces, Azure Network Manager IPAM delegation, custom
DNS servers, BGP community advertisement, DDoS Protection Plan attachment,
VM-to-VM encryption, connection-tracking flow timeouts, edge zones, and
governance tags.

## The Network Is Just the Network

The virtual network deliberately contains nothing but the address space and
network-wide policy. Everything that partitions or extends it is its own
composable resource referencing this network's outputs:

- **`AzureSubnet`** partitions the address space into workload segments
- **`AzureNatGateway`** provides managed outbound connectivity for a subnet
- **`AzurePrivateDnsZoneVirtualNetworkLink`** makes a private DNS zone
  resolvable from this network

This separation means a hub-and-spoke topology adds subnets, DNS links, and
gateways over time without ever touching -- or risking -- the network
resource itself.

## Key Features

- **Multi-CIDR address space** -- networks routinely grow a second range
  when the first fills; dual-stack networks carry IPv4 and IPv6 side by side
- **IPAM delegation** -- alternatively, delegate address allocation to an
  Azure Network Manager IPAM pool (`ip_address_pools`); the provisioned
  ranges surface in the outputs
- **Custom DNS servers** -- for on-premises DNS integration or self-hosted
  resolvers (empty means Azure's default resolver, which private DNS zone
  resolution requires)
- **DDoS Protection Plan attachment** -- attach a shared plan by ARM ID,
  with ARM's attach/activate distinction preserved
- **VM-to-VM encryption** -- enable virtual network encryption with the
  enforcement mode for non-encrypting VM sizes
- **BGP community** -- advertise a community with the network's routes over
  ExpressRoute for on-premises route filtering
- **Composable** -- the resource group is referenced by name, defaulting to
  an `AzureResourceGroup`'s output in composed environments

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | Azure region for the network (a regional resource) |
| `resource_group` | StringValueOrRef | Yes | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | Network name, unique within the resource group (2-64 chars) |
| `address_spaces` | list(string) | One of | Self-managed CIDR blocks (exactly one of this or `ip_address_pools`) |
| `ip_address_pools` | list | One of | Network Manager IPAM pool allocations (max 2: one per IP version) |
| `dns_servers` | list(string) | No | Custom DNS server IPs (empty = Azure's default resolver) |
| `bgp_community` | string | No | BGP community in `asn:community` notation (ASN is always 12076 today) |
| `ddos_protection_plan` | object | No | Attach an existing DDoS Protection Plan (`id` + `enable`) |
| `encryption` | enum | No | VM-to-VM encryption enforcement: `ALLOW_UNENCRYPTED` or `DROP_UNENCRYPTED` |
| `flow_timeout_in_minutes` | int | No | Connection-tracking timeout, 4-30 (unset = Azure's 4-minute default) |
| `private_endpoint_vnet_policies` | enum | No | Network-wide private endpoint policy: `BASIC` (unset = Disabled) |
| `edge_zone` | string | No | Deploy into an Azure Edge Zone instead of the main region |
| `tags` | map | No | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `virtual_network_id` | Full ARM ID of the network -- the join key for subnets, peerings, and DNS links |
| `virtual_network_name` | The network's name as deployed |
| `guid` | ARM's stable network GUID (used by BGP community advertisement and diagnostics) |
| `address_spaces` | The ACTUAL address ranges (IPAM-provisioned when pools delegate allocation) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureVirtualNetwork
metadata:
  name: prod-network
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: prod-vnet
  addressSpaces:
    - "10.0.0.0/16"
  tags:
    cost-center: platform
```

A complete network foundation composes the network with subnets and
attachments -- see the presets for hub-and-spoke, dual-stack, and
DDoS-protected starting points.

## Lifecycle Notes

- Address space, DNS servers, BGP community, DDoS attachment, encryption,
  flow timeout, and tags all update **in place**
- Name, region, resource group, and edge zone are the network's ARM
  identity -- changing any of them **replaces the network and everything
  inside it**
- Address-space blocks can be added and removed live, but a block that
  subnets are carved from cannot shrink below them
