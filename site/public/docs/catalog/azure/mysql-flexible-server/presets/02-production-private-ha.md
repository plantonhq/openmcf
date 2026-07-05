---
title: "Production Private Server with Zone-Redundant HA"
description: "This preset creates a production-grade MySQL Flexible Server: General Purpose compute, 256 GiB storage with auto-grow and elastic IOPS scaling, a zone-redundant standby with automatic failover,..."
type: "preset"
rank: "02"
presetSlug: "02-production-private-ha"
componentSlug: "mysql-flexible-server"
componentTitle: "MySQL Flexible Server"
provider: "azure"
icon: "package"
order: 2
---

# Production Private Server with Zone-Redundant HA

This preset creates a production-grade MySQL Flexible Server: General
Purpose compute, 256 GiB storage with auto-grow and elastic IOPS scaling,
a zone-redundant standby with automatic failover, 35-day geo-redundant
backups, VNet injection with no public endpoint, a pinned maintenance
window, and TLS-only connections.

## When to Use

- Production workloads that need zone-level fault tolerance and private
  connectivity
- Any environment where the database must never be reachable from the
  public internet

## Key Configuration Choices

- **VNet injection over a public endpoint** -- the server lives on a
  subnet delegated to `Microsoft.DBforMySQL/flexibleServers`;
  `publicNetworkAccess` stays unset so Azure derives DISABLED
- **`ZONE_REDUNDANT` HA with pinned zones** -- the standby synchronously
  replicates in another zone; the zone pair only changes via a planned
  failover
- **Elastic IOPS scaling (`ioScalingEnabled`)** -- Azure scales IOPS with
  demand instead of holding a provisioned value; mutually exclusive with
  a fixed `iops` dial
- **`require_secure_transport: ON`** -- rejects any client that does not
  negotiate TLS
- **Geo-redundant backups** -- enables cross-region GEO_RESTORE for
  disaster recovery (fixed at creation)

## Prerequisites

- An `AzureSubnet` delegated to `Microsoft.DBforMySQL/flexibleServers`
  with no other resources in it
- An `AzurePrivateDnsZone` (conventionally ending
  `.mysql.database.azure.com`) linked to the VNet

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the server in | The resource group's `status.outputs.resource_group_name` |
| `<globally-unique-server-name>` | 3-63 lowercase chars, globally unique | Your naming convention |
| `<admin-login>` | The MySQL admin login | Your convention |
| `<admin-password>` | 8-128 chars from 3+ character classes | A secret manager; never commit literals |
| `<database-subnet-resource-name>` | The delegated subnet's Planton resource name | Your network composition |
| `<mysql-dns-zone-resource-name>` | The private DNS zone's Planton resource name | Your network composition |
| `<application-database-name>` | The application's database | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

The `fqdn` output resolves to the server's private address through the
private DNS zone, so VNet-connected applications use the same connection
string shape as a public server:

```text
mysql://{administrator_login}:{password}@{status.outputs.fqdn}:3306/{database}?ssl-mode=REQUIRED
```
