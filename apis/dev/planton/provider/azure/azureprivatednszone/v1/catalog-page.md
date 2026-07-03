# Azure Private DNS Zone

Creates an Azure Private DNS zone -- name resolution inside virtual networks without running a DNS server. Powers both Private Link scenarios (resolving Azure PaaS private endpoints via the `privatelink.*` zones) and custom internal DNS (`corp.internal`).

## What Gets Created

When you deploy an AzurePrivateDnsZone resource, Planton provisions:

- **Private DNS Zone** — an `azurerm_private_dns_zone` in the specified resource group. Private DNS zones are global Azure resources with no region parameter.

The zone is deliberately just the zone. Which networks can resolve it is declared through `AzurePrivateDnsZoneVirtualNetworkLink` resources referencing this zone's `zone_id` output -- one link per network, added and removed without touching the zone. A zone with no links answers nobody, so pair the zone with at least one link.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the zone in (an `AzureResourceGroup` in composed environments)
- **Zone name planning** — for Private Link scenarios, the zone name must match the Azure-defined privatelink zone name for the target service (e.g., `privatelink.postgres.database.azure.com` for PostgreSQL Flexible Server)
- **Private DNS write rights**: `Microsoft.Network/privateDnsZones/write` (Private DNS Zone Contributor, Contributor, or Owner)

## Quick Start

Create a file `private-dns-zone.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZone
metadata:
  name: postgres-privatelink-zone
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZone.postgres-privatelink-zone
spec:
  resourceGroup:
    value: network-rg
  name: privatelink.postgres.database.azure.com
```

Deploy:

```shell
planton apply -f private-dns-zone.yaml
```

Then make it resolvable from a network with an `AzurePrivateDnsZoneVirtualNetworkLink` referencing this zone.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | DNS zone name. For Private Link, must match the Azure-defined privatelink zone name for the target service; for custom internal DNS, any valid domain. Renaming replaces the zone and every record in it. | Required, valid DNS domain |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `soaRecord` | `object` | Customize the zone's Start of Authority record: `email` (SOA host format, dots instead of `@`) plus `expireTime`, `minimumTtl` (negative-caching TTL -- the one timer with real operational impact), `refreshTime`, `retryTime`, `ttl`, and record `tags`. Written at creation; changing it replaces the zone. Nearly every deployment leaves this unset. |
| `tags` | `map(string)` | User tags, merged over Planton-derived tags (user wins on collision). |

## Examples

### Private Link Zone for PostgreSQL

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZone
metadata:
  name: postgres-privatelink-zone
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZone.postgres-privatelink-zone
spec:
  resourceGroup:
    valueFrom:
      name: network-rg
  name: privatelink.postgres.database.azure.com
```

### Custom Internal Zone with SOA Tuning

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZone
metadata:
  name: corp-internal-zone
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzurePrivateDnsZone.corp-internal-zone
spec:
  resourceGroup:
    valueFrom:
      name: network-rg
  name: corp.internal
  soaRecord:
    email: dnsadmin.mycompany.com
    minimumTtl: 30
  tags:
    cost-center: platform
```

### The Composed Trio: Network + Zone + Link

```yaml
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZone
metadata:
  name: redis-privatelink-zone
spec:
  resourceGroup:
    valueFrom:
      name: network-rg
  name: privatelink.redis.cache.windows.net
---
apiVersion: azure.planton.dev/v1
kind: AzurePrivateDnsZoneVirtualNetworkLink
metadata:
  name: redis-zone-hub-link
spec:
  name: hub-vnet
  privateDnsZoneId:
    valueFrom:
      name: redis-privatelink-zone
  virtualNetworkId:
    valueFrom:
      name: hub-network
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | `string` | Full ARM ID of the zone -- the join key for links, private endpoints, and databases |
| `zone_name` | `string` | The zone's DNS name as deployed |
| `resource_group_name` | `string` | The zone's resource group |

## Related Components

- [AzurePrivateDnsZoneVirtualNetworkLink](/docs/catalog/azure/private-dns-zone-virtual-network-link) — makes this zone resolvable from a virtual network
- [AzurePrivateEndpoint](/docs/catalog/azure/private-endpoint) — registers PaaS private-IP records in privatelink zones
- [AzureVirtualNetwork](/docs/catalog/azure/virtual-network) — the network the zone gets linked to
- [AzurePostgresqlFlexibleServer](/docs/catalog/azure/postgresql-flexible-server) — consumes the zone for VNet-integrated deployment
