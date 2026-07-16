---
title: "Private Link Zone Attachment"
description: "This preset links a Private Link zone (e.g. `privatelink.postgres.database.azure.com`) to a virtual network so workloads inside it resolve the PaaS service's FQDN to its private endpoint IP instead..."
type: "preset"
rank: "01"
presetSlug: "01-privatelink-zone-link"
componentSlug: "private-dns-zone-virtual-network-link"
componentTitle: "Private DNS Zone Virtual Network Link"
provider: "azure"
icon: "package"
order: 1
---

# Private Link Zone Attachment

This preset links a Private Link zone (e.g.
`privatelink.postgres.database.azure.com`) to a virtual network so
workloads inside it resolve the PaaS service's FQDN to its private
endpoint IP instead of the public one -- the standard last step of every
private-endpoint deployment.

VM auto-registration stays off (the default): privatelink zone records
are created and removed by private endpoints, never by VM registration.

## When to Use

- After creating a private endpoint for PostgreSQL, MySQL, SQL, Cosmos,
  Redis, Storage, or Key Vault -- the zone must be resolvable from every
  network whose workloads call the service
- Extending an existing privatelink zone's audience to a new network

## Key Configuration Choices

- **Name the link after the network** -- a zone accumulates one link per
  network; the name is how they are told apart
- **Registration stays off** -- only custom internal zones
  (`corp.internal`) auto-register VM records
- **Use `valueFrom`** against an `AzurePrivateDnsZone` and an
  `AzureVirtualNetwork` in composed environments instead of literal IDs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `app-vnet-link` | A short name identifying the linked network | Your naming convention |
| `<private-dns-zone-arm-id>` | The zone's full ARM ID | The zone's `status.outputs.zone_id` |
| `<virtual-network-arm-id>` | The network's full ARM ID | The network's `status.outputs.virtual_network_id` |
