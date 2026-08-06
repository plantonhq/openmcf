# AzurePublicIpPrefix Terraform Module

## Overview

This Terraform module provisions an Azure Public IP Prefix using the
`azurerm` provider. It creates a single `azurerm_public_ip_prefix` -- a
reserved, contiguous range of public IP addresses that public IPs allocate
from and NAT gateways associate for outbound SNAT.

The prefix is essentially immutable: everything except tags is fixed at
creation, and replacing it changes the ACTUAL reserved range -- everything
partners have allowlisted. Treat replacement as a coordinated migration,
never a casual update. The prefix cannot be deleted while any of its
addresses are in use.

Only explicit choices are sent to Azure: an unspecified prefix length,
version, SKU, or tier falls through to Azure's defaults, so an unspecified
spec deploys identically on both engines.

## Resources Created

- `azurerm_public_ip_prefix.main` -- the reserved address range

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Public IP prefix specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; addresses can only be allocated by resources in the same region |
| `resource_group` | yes | Resource group name |
| `name` | yes | Prefix name, unique within the resource group |
| `prefix_length` | no | CIDR length; Azure default 28 (16 addresses), /31 reserves 2 -- smaller numbers reserve bigger ranges and bill for every address |
| `ip_version` | no | `IPV4` (Azure's default) or `IPV6` |
| `sku` / `sku_tier` | no | `STANDARD`/`REGIONAL` are Azure's defaults; `STANDARD_V2` pairs with StandardV2 NAT gateways; `GLOBAL` requires `STANDARD` |
| `zones` | no | Availability zones anchoring the range; multiple zones make it zone-redundant |
| `custom_ip_prefix_id` | no | BYOIP custom prefix ARM ID to carve the range from |
| `tags` | no | User tags, merged over metadata-derived tags; the only field that updates in place |

## Outputs

| Output | Description |
|--------|-------------|
| `public_ip_prefix_id` | Full ARM ID of the prefix |
| `ip_prefix` | The actual reserved CIDR (e.g. `20.42.0.16/28`), known only after creation |
| `public_ip_prefix_name` | The prefix's name as deployed |

## Usage

```hcl
module "public_ip_prefix" {
  source = "./iac/tf"

  metadata = {
    name = "prod-egress"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "prod-egress"
    prefix_length  = 28
    zones          = ["1", "2", "3"]
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Network/publicIPPrefixes/write`
on the resource group -- held via Network Contributor, Contributor, or
Owner.
