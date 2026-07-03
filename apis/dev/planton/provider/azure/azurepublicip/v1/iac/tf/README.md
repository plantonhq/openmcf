# AzurePublicIp Terraform Module

## Overview

This Terraform module provisions an Azure Public IP Address using the
`azurerm` provider. It creates a single `azurerm_public_ip` with static
allocation (every current SKU requires it), covering the full azurerm
surface: SKU and tier, IP version, zones, prefix allocation, DNS label with
scope-based reuse, reverse FQDN, idle timeout, IP tags, DDoS protection
stance, and edge zones.

Reverse FQDN, DDoS settings, idle timeout, and tags update in place. Name,
SKU/tier, IP version, zones, prefix membership, IP tags, and edge zone are
fixed at creation -- changing any of them replaces the resource and with it
the actual address, so treat replacement as a coordinated migration (DNS,
allowlists).

Enum fields are mapped to provider strings in locals and left null when
unset, so an unspecified spec deploys Azure's defaults (Standard / Regional /
IPv4 / region-unique label / inherited DDoS stance) identically on both
engines.

## Resources Created

- `azurerm_public_ip.main` -- the static public IP address

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Public IP specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region (must match the region of the resource the address attaches to) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Address name, unique within the resource group; renaming replaces the resource and the actual address |
| `sku` | no | `STANDARD` or `STANDARD_V2` (enum name string); unset defers to Azure's default (Standard) |
| `sku_tier` | no | `REGIONAL` or `GLOBAL` (enum name string); GLOBAL requires the STANDARD SKU |
| `ip_version` | no | `IPV4` or `IPV6` (enum name string); unset defers to Azure's default (IPv4) |
| `zones` | no | Availability zones ("1", "2", "3"); multiple zones make the address zone-redundant |
| `public_ip_prefix_id` | no | ARM ID of the public IP prefix to allocate from |
| `domain_name_label` | no | Azure-managed DNS label (`{label}.{region}.cloudapp.azure.com`) |
| `domain_name_label_scope` | no | Label reuse policy (enum name string, e.g. `TENANT_REUSE`); unset keeps the classic region-unique behavior |
| `reverse_fqdn` | no | Reverse-DNS (PTR) name; the forward record must exist first |
| `idle_timeout_in_minutes` | no | TCP idle timeout in minutes, 4-30 (Azure defaults to 4) |
| `ip_tags` | no | Azure IP tags (routing metadata like `RoutingPreference`), not governance tags |
| `ddos_protection_mode` | no | `DISABLED` or `ENABLED` (enum name string); unset inherits from the network |
| `ddos_protection_plan_id` | no | ARM ID of the DDoS plan; only valid with the ENABLED mode |
| `edge_zone` | no | Azure Edge Zone; unset deploys into the standard region |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `public_ip_id` | Full ARM ID of the address |
| `ip_address` | The allocated address; static for the resource's lifetime |
| `fqdn` | The Azure-managed FQDN; populated only when a domain name label is set |
| `public_ip_name` | The address's name as deployed |

## Usage

```hcl
module "public_ip" {
  source = "./iac/tf"

  metadata = {
    name = "prod-gateway-pip"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region                  = "eastus"
    resource_group          = "network-rg"
    name                    = "prod-gateway-pip"
    zones                   = ["1", "2", "3"]
    domain_name_label       = "prod-gateway"
    domain_name_label_scope = "TENANT_REUSE"
    idle_timeout_in_minutes = 15
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/publicIPAddresses/write`
on the resource group -- held via Network Contributor, Contributor, or Owner.
