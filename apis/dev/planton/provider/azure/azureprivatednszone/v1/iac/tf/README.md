# AzurePrivateDnsZone Terraform Module

## Overview

This Terraform module provisions an Azure Private DNS zone using the
`azurerm` provider. It creates a single `azurerm_private_dns_zone` -- a
global record container -- with optional SOA customization and governance
tags.

Tags update in place; the zone's name is its ARM identity, so renaming
replaces the zone and every record in it. The SOA record is written at
creation and cannot be customized afterwards.

The zone is deliberately just the zone: which networks can resolve it is
declared through `AzurePrivateDnsZoneVirtualNetworkLink` resources
referencing this zone's `zone_id` output -- one link per network. A zone
with no links answers nobody.

## Resources Created

- `azurerm_private_dns_zone.main` -- the private DNS zone

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Private DNS zone specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `resource_group` | yes | Resource group name |
| `name` | yes | The zone's DNS domain name (renaming replaces the zone and its records) |
| `soa_record` | no | SOA customization: `email` (required inside the block) + `expire_time` / `minimum_ttl` / `refresh_time` / `retry_time` / `ttl` / record `tags` |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `zone_id` | Full ARM ID of the zone |
| `zone_name` | The zone's DNS name as deployed |
| `resource_group_name` | The zone's resource group |

## Usage

```hcl
module "private_dns_zone" {
  source = "./iac/tf"

  metadata = {
    name = "postgres-privatelink-zone"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    resource_group = "network-rg"
    name           = "privatelink.postgres.database.azure.com"
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/privateDnsZones/write`
on the resource group -- held via Private DNS Zone Contributor,
Contributor, or Owner.
