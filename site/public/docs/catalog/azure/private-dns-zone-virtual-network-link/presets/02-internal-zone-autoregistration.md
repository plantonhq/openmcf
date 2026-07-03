---
title: "Internal Zone with VM Auto-Registration"
description: "This preset links a custom internal zone (e.g. `corp.internal`) to a virtual network with VM auto-registration on: every virtual machine in the network gets an A record at boot and loses it at..."
type: "preset"
rank: "02"
presetSlug: "02-internal-zone-autoregistration"
componentSlug: "private-dns-zone-virtual-network-link"
componentTitle: "Private DNS Zone Virtual Network Link"
provider: "azure"
icon: "package"
order: 2
---

# Internal Zone with VM Auto-Registration

This preset links a custom internal zone (e.g. `corp.internal`) to a
virtual network with VM auto-registration on: every virtual machine in
the network gets an A record at boot and loses it at deallocation, making
machines discoverable by hostname with zero record management.

Azure allows only ONE registration-enabled link per virtual network --
choose the zone that owns the network's machine names. Additional zones
link to the same network with registration off.

## When to Use

- Custom internal DNS zones where VMs should be discoverable by hostname
- Lift-and-shift estates that relied on on-premises dynamic DNS

## Key Configuration Choices

- **`registrationEnabled: true`** -- the one registration-enabled link
  this network gets; a second one fails at deploy time
- **Never for privatelink zones** -- their records belong to private
  endpoints; registration would add unrelated VM records to a
  service-resolution zone

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<network-name>` | A short name identifying the linked network | Your naming convention |
| `<private-dns-zone-arm-id>` | The internal zone's full ARM ID | The zone's `status.outputs.zone_id` |
| `<virtual-network-arm-id>` | The network's full ARM ID | The network's `status.outputs.virtual_network_id` |
