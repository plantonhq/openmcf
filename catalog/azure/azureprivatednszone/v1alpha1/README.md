# AzurePrivateDnsZone

## Overview

`AzurePrivateDnsZone` provisions an Azure Private DNS zone: name
resolution inside virtual networks without running a DNS server. Private
DNS zones power two scenarios -- automatic resolution for Azure Private
Endpoints (the `privatelink.*` zones) and custom internal DNS
(`corp.internal`).

## The Zone Is Just the Zone

The zone is a global record container with no reach of its own. Which
networks can resolve it is declared through
`AzurePrivateDnsZoneVirtualNetworkLink` resources referencing this zone's
`zone_id` output -- one link per network, added and removed without
touching the zone. A zone with no links answers nobody, so every
deployment pairs the zone with at least one link.

This separation is what makes hub-and-spoke DNS work: one
`privatelink.postgres.database.azure.com` zone, linked to the hub and
every spoke, each link an independently reviewable and removable node.

## Key Features

- **Private Link ready** -- create the Azure-defined `privatelink.*` zone
  names for automatic private endpoint resolution
- **Custom internal DNS** -- any valid domain for internal name resolution
- **SOA customization** -- pin the zone's contact email and
  negative-caching TTL when operational tooling requires it
- **Governance tags** -- user tags merged over the Planton-derived ones
- **Composable** -- the resource group is referenced by name, defaulting
  to an `AzureResourceGroup`'s output; the zone's outputs feed links,
  private endpoints, and VNet-integrated databases

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `resource_group` | StringValueOrRef | Yes | Resource group name (defaults to an AzureResourceGroup reference) |
| `name` | string | Yes | The zone's DNS domain name (privatelink zone name or custom domain) |
| `soa_record` | object | No | SOA customization (email + timers); written at creation, ForceNew |
| `tags` | map | No | User tags, merged over Planton-derived tags (user wins) |

Private DNS zones are **global** resources -- there is no `region` field.

## Outputs

| Output | Description |
|--------|-------------|
| `zone_id` | Full ARM ID of the zone -- the join key for links, private endpoints, and databases |
| `zone_name` | The zone's DNS name as deployed |
| `resource_group_name` | The zone's resource group (for record/link tooling that joins on name+RG) |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateDnsZone
metadata:
  name: postgres-privatelink-zone
  org: mycompany
  env: production
spec:
  resourceGroup:
    valueFrom:
      name: network-rg
  name: privatelink.postgres.database.azure.com
```

The complete private-DNS story composes three kinds: the network
(`AzureVirtualNetwork`), this zone, and an
`AzurePrivateDnsZoneVirtualNetworkLink` per network that should resolve
it.

## Common Private Link Zone Names

| Service | Zone name |
|---------|-----------|
| PostgreSQL Flexible Server | `privatelink.postgres.database.azure.com` |
| MySQL Flexible Server | `privatelink.mysql.database.azure.com` |
| Azure SQL Database | `privatelink.database.windows.net` |
| Cosmos DB (SQL) | `privatelink.documents.azure.com` |
| Azure Cache for Redis | `privatelink.redis.cache.windows.net` |
| Blob Storage | `privatelink.blob.core.windows.net` |
| Key Vault | `privatelink.vaultcore.azure.net` |

## Lifecycle Notes

- Tags update **in place**; the zone's name is its ARM identity, so
  renaming **replaces the zone and every record in it**
- The SOA record is written at creation and cannot be customized
  afterwards
- Deleting a zone requires its links to be gone first -- composed destroy
  ordering handles this through the dependency graph

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
