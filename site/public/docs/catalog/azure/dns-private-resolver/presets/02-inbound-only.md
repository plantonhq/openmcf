---
title: "Inbound Only (Pinned IP)"
description: "This preset deploys the one-way shape: on-premises (or another cloud) resolves Azure's private names by forwarding to the inbound endpoint, whose address is STATICALLY pinned so forwarder..."
type: "preset"
rank: "02"
presetSlug: "02-inbound-only"
componentSlug: "dns-private-resolver"
componentTitle: "DNS Private Resolver"
provider: "azure"
icon: "package"
order: 2
---

# Inbound Only (Pinned IP)

This preset deploys the one-way shape: on-premises (or another cloud) resolves Azure's private names by forwarding to the inbound endpoint, whose address is STATICALLY pinned so forwarder configurations never need to change. No outbound endpoint -- Azure-side resolution of on-premises names is not wired.

## When to Use

- The datacenter needs to resolve Azure private endpoints and private zones, but Azure workloads never query on-premises domains
- Environments where the forwarder fleet's configuration is fanned out and expensive to change (the pinned IP survives endpoint replacement)

## Key Configuration Choices

- **STATIC allocation requires the address** -- pick a free one from the delegated subnet's range (the first four addresses of any subnet are Azure-reserved)
- **Replace `10.0.4.4`** with an address that actually belongs to your inbound subnet
- **Adding an outbound endpoint later is additive** -- endpoints are keyed by name, so extending the list never touches the existing endpoint

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the resolver is created in | The resource-group component's name |
| `<your-hub-vnet>` | The AzureVirtualNetwork the resolver anchors to | The network component's name |
| `<your-dns-inbound-subnet>` | The delegated AzureSubnet for the inbound endpoint | The subnet component's name |

After deploy, `inbound_endpoint_ip` echoes the pinned address -- point every datacenter conditional forwarder for Azure domains at it.
