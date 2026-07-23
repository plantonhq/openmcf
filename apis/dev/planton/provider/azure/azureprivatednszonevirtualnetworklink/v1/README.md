# AzurePrivateDnsZoneVirtualNetworkLink

## Overview

`AzurePrivateDnsZoneVirtualNetworkLink` provisions the attachment that
makes an Azure Private DNS zone resolvable from a virtual network. A zone
without links answers nobody; each link adds one network to its audience,
optionally with automatic VM record registration and a resolution-fallback
policy.

## Why a First-Class Resource?

Links are real infrastructure with their own lifecycle:

- **Many per zone** -- a hub-and-spoke topology links one zone (say,
  `privatelink.postgres.database.azure.com`) to the hub and every spoke
  network
- **Independent lifecycle** -- networks join and leave the resolution
  audience without touching the zone, its records, or each other
- **One link resource per zone-network pair** -- each individually
  reviewable and removable

## Key Features

- **Parent-derived identity** -- the link takes the zone's ARM ID; the
  zone's name and resource group are derived from it, so they can never
  contradict the referenced zone
- **VM auto-registration** -- optionally auto-register A records for VMs
  in the linked network (custom internal zones; one registration-enabled
  link per network)
- **Resolution policy** -- opt into `NX_DOMAIN_REDIRECT` to retry names
  the private zone cannot answer against public DNS
- **Composable** -- both the zone and the network are referenced by ARM
  ID, defaulting to `AzurePrivateDnsZone` and `AzureVirtualNetwork`
  outputs in composed environments

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | The link's name under the parent zone (1-80 chars, unique per zone) |
| `private_dns_zone_id` | StringValueOrRef | Yes | ARM ID of the parent zone (defaults to an AzurePrivateDnsZone reference) |
| `virtual_network_id` | StringValueOrRef | Yes | ARM ID of the network the zone becomes resolvable from (defaults to an AzureVirtualNetwork reference) |
| `registration_enabled` | bool | No | Auto-register VM A records (default false; keep false for privatelink zones) |
| `resolution_policy` | enum | No | `DEFAULT` or `NX_DOMAIN_REDIRECT` (unset = Azure's per-zone-type default) |
| `tags` | map | No | User tags, merged over Planton-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `link_id` | Full ARM ID of the link (`{zone-id}/virtualNetworkLinks/{name}`) |
| `link_name` | The link's name as deployed |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: postgres-zone-hub-link
  org: mycompany
  env: production
spec:
  name: hub-vnet
  privateDnsZoneId:
    valueFrom:
      name: postgres-privatelink-zone
  virtualNetworkId:
    valueFrom:
      name: hub-network
```

The complete private-DNS story composes three kinds: the network
(`AzureVirtualNetwork`), the zone (`AzurePrivateDnsZone`), and this link
making the zone resolvable from the network -- one link per zone-network
pair.

## Lifecycle Notes

- `registration_enabled`, `resolution_policy`, and tags update **in
  place**; name, zone, and network are the link's ARM identity, so
  changing any of them **replaces the link** (a brief resolution gap for
  the affected network, nothing else)
- Azure allows only **one registration-enabled link per virtual network**
- The deploying credential needs
  `Microsoft.Network/privateDnsZones/virtualNetworkLinks/write` on the
  zone and `Microsoft.Network/virtualNetworks/join/action` on the network

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
