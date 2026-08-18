# AzureVpnSite Terraform Module

## Overview

This Terraform module provisions a VPN Site using the `azurerm`
provider. It creates a single `azurerm_vpn_site` -- the Virtual WAN
address-book entry for one branch location, with its links, reachable
address space, device metadata, and O365 breakout policy.

## Resources Created

- `azurerm_vpn_site.main` -- the branch description (links are inline
  blocks; ARM assigns each an ID)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | VPN Site specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | The site object's region (ForceNew) |
| `resource_group` | yes | Resource group name (ForceNew) |
| `name` | yes | Must not contain `' < > % & : ? / +` (ForceNew) |
| `virtual_wan_id` | yes | The WAN the site belongs to (ForceNew) |
| `address_cidrs` | no | Prefixes routed into the tunnels; empty only when every link speaks BGP |
| `links` | no | The connectable endpoints; each needs an IP or FQDN |
| `o365_policy` | no | O365 breakout categories (SD-WAN partners read it) |

## Outputs

| Output | Description |
|--------|-------------|
| `vpn_site_id` | Full ARM ID -- what a connection's `remote_vpn_site_id` references |
| `vpn_site_name` | The site's name |
| `link_ids` | Each link's ARM ID keyed by link name -- what a connection's tunnels pin to |

## Usage

```hcl
module "branch_london" {
  source = "./iac/tf"

  metadata = { name = "branch-london", org = "mycompany", env = "production" }

  spec = {
    name           = "branch-london"
    region         = "eastus"
    resource_group = "network-rg"
    virtual_wan_id = "/subscriptions/.../virtualWans/corp-wan"
    address_cidrs  = ["192.168.10.0/24"]
    links = [
      { name = "primary-isp", ip_address = "203.0.113.10", speed_in_mbps = 200 }
    ]
  }
}
```

## Behavior Notes

- Empty optional strings (`ip_address`, `fqdn`, `provider_name`,
  `device_vendor`, `device_model`) are emitted as null -- the provider
  validates configured values as non-empty.
- The spec guarantees every link carries at least one endpoint (IP or
  FQDN) -- ARM's create-time rule, front-loaded.
- The site is free and provisions in seconds; deleting it requires the
  connections pointing at it to be gone first.
