# AzurePointToSiteVpnGateway Terraform Module

## Overview

This Terraform module provisions a Point-to-Site VPN Gateway using the
`azurerm` provider. It creates a single
`azurerm_point_to_site_vpn_gateway` -- the managed receiver inside a
Virtual WAN hub that individual devices dial into, authenticated per
the VPN Server Configuration it references.

## Resources Created

- `azurerm_point_to_site_vpn_gateway.main` -- the gateway with its
  connection configurations (client address pools, per-pool routing)
  as inline blocks

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Point-to-Site VPN Gateway specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Must match the hub's region (ForceNew) |
| `resource_group` | yes | Resource group name (ForceNew) |
| `name` | yes | The gateway's name (ForceNew) |
| `virtual_hub_id` | yes | The hub the gateway deploys into -- one P2S gateway per hub (ForceNew) |
| `vpn_server_configuration_id` | yes | The authentication policy (ForceNew) |
| `connection_configurations` | yes | Named client address pools (at least one) |
| `scale_unit` | no | 500 connections each; unset applies 1 (the provider requires an explicit value) |
| `routing_preference_internet_enabled` | no | Hot-potato internet egress for the gateway's interface (ForceNew) |
| `dns_servers` | no | Pushed to clients; cannot be cleared in place once set |

## Outputs

| Output | Description |
|--------|-------------|
| `point_to_site_vpn_gateway_id` | Full ARM ID of the gateway |
| `point_to_site_vpn_gateway_name` | The gateway's name |

## Usage

```hcl
module "remote_users_gw" {
  source = "./iac/tf"

  metadata = { name = "remote-users-gw", org = "mycompany", env = "production" }

  spec = {
    name                        = "remote-users-gw"
    region                      = "eastus"
    resource_group              = "network-rg"
    virtual_hub_id              = "/subscriptions/.../virtualHubs/hub-eastus"
    vpn_server_configuration_id = "/subscriptions/.../vpnServerConfigurations/remote-workforce"
    connection_configurations = [
      { name = "default-clients", address_prefixes = ["172.16.201.0/24"] }
    ]
  }
}
```

## Behavior Notes

- The gateway is a SLOW, billing resource: creates run 30-45 minutes
  (the provider's timeout class is 90) and billing starts at creation.
- `scale_unit` is Required on the provider with no default; the module
  renders the spec's unset as an explicit 1 (recorded in the parity
  manifest).
- A configured `route` block requires its `associated_route_table_id`
  (the provider's contract, front-loaded in the spec); unset routing
  applies ARM's default association/propagation.
- The provider cannot CLEAR `dns_servers` once set (its update path
  skips empty lists) -- removing servers requires replacing the
  gateway.

## Required Permissions

The deploying credential needs
`Microsoft.Network/p2sVpnGateways/write` plus read on the hub and VPN
server configuration -- held via Network Contributor, Contributor, or
Owner on the resource group.
