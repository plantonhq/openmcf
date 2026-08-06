# AzureVirtualNetworkPeering

## Overview

`AzureVirtualNetworkPeering` provisions one direction of an Azure virtual
network peering: private, low-latency connectivity between two virtual
networks over the Microsoft backbone, without gateways, public IPs, or
encryption overhead. Peered networks exchange traffic as if they were one
network while remaining separately owned and managed -- the building block
of hub-and-spoke and multi-VNet topologies.

## Why a First-Class Resource?

A peering is real infrastructure with its own lifecycle:

- **Many per network** -- a hub accumulates one peering per spoke; spokes
  come and go without touching the networks themselves
- **One direction per resource** -- ARM models each direction separately;
  connectivity only flows once both directions exist
- **Independent policy** -- access, forwarded-traffic, and gateway-transit
  flags differ by direction and evolve without replacing the networks
- **Cross-subscription and global** -- the remote network is referenced by
  ARM ID alone; no extra wiring for cross-subscription or cross-region pairs

## Key Features

- **Directional grain** -- one resource is one peering written on the local
  network; charts stamp the reciprocal direction as a sibling resource
- **Full connectivity surface** -- allow virtual network access, allow
  forwarded traffic, gateway transit, and use-remote-gateways
- **Subnet-scoped peering** -- peer complete address spaces (default) or
  named subnets only, with optional IPv6-only peering for dual-stack designs
- **FK-driven placement** -- the local network's resource group and name are
  derived from `virtual_network_id`; no duplicate identity fields to drift
- **Validated pairing** -- subnet name lists are rejected unless complete-
  network peering is disabled
- **No tags** -- ARM peerings are not taggable resources

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Peering name, unique within the local network (1-80 chars) |
| `virtual_network_id` | StringValueOrRef | Yes | Local network ARM ID (defaults to an AzureVirtualNetwork reference) |
| `remote_virtual_network_id` | StringValueOrRef | Yes | Remote network ARM ID (defaults to an AzureVirtualNetwork reference) |
| `allow_virtual_network_access` | bool | No | Whether local traffic can reach the remote network (default: true) |
| `allow_forwarded_traffic` | bool | No | Whether forwarded traffic from the remote network is accepted (default: false) |
| `allow_gateway_transit` | bool | No | Whether the local gateway may be used by the remote network (default: false) |
| `use_remote_gateways` | bool | No | Whether the local network uses the remote gateway (default: false) |
| `peer_complete_virtual_networks_enabled` | bool | No | Peer full address spaces vs named subnets only (default: true) |
| `local_subnet_names` | list(string) | No | Local subnets included when subnet-scoped peering is enabled |
| `remote_subnet_names` | list(string) | No | Remote subnets included when subnet-scoped peering is enabled |
| `only_ipv6_peering_enabled` | bool | No | Peer only IPv6 address space (subnet-scoped dual-stack; default: false) |

## Outputs

| Output | Description |
|--------|-------------|
| `peering_id` | Full ARM ID of the peering |
| `peering_name` | The peering's name within its local network |
| `virtual_network_name` | The local network's name, derived from `virtual_network_id` |
| `resource_group_name` | The local network's resource group, derived from `virtual_network_id` |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualNetworkPeering
metadata:
  name: hub-to-spoke1
  org: mycompany
  env: production
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

A working pair requires a second resource with local and remote swapped
(`spoke1-to-hub`). See the presets for the full hub-and-spoke pattern.

## Composition

- `virtual_network_id` → `AzureVirtualNetwork.status.outputs.virtual_network_id`
  (the local side the peering is written on)
- `remote_virtual_network_id` → `AzureVirtualNetwork.status.outputs.virtual_network_id`
  (the far side; may be a different stack, subscription, or region)
- Hub-and-spoke charts stamp **two** `AzureVirtualNetworkPeering` resources per
  spoke pair: hub→spoke and spoke→hub, with direction-appropriate gateway flags
- `peering_id` and `peering_name` identify this direction for diagnostics;
  `virtual_network_name` and `resource_group_name` let sibling resources compose
  without re-parsing the local network's ARM ID

## Lifecycle Notes

- Access, forwarding, gateway, and subnet-name fields update **in place**
- Name, both network references, complete-vs-subnet-scoped peering, and
  IPv6-only mode are the peering's ARM identity -- changing any of them
  **replaces the peering** (a brief connectivity gap for this direction only)
- The reciprocal direction can deploy concurrently; Azure retries until both
  sides exist
- `use_remote_gateways` cannot combine with global (cross-region) peering, and
  only one peering per network may set it

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
