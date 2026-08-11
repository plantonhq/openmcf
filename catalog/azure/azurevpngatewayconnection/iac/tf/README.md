# AzureVpnGatewayConnection Terraform Module

## Overview

This Terraform module provisions a VPN Gateway Connection using the
`azurerm` provider. It creates a single `azurerm_vpn_gateway_connection`
-- the tunnel bundle joining one branch (a VPN Site) to a Virtual WAN
hub's VPN gateway, with per-tunnel IPsec, BGP, and NAT choices.

## Resources Created

- `azurerm_vpn_gateway_connection.main` -- the connection (an ARM child
  of the gateway) carrying one `vpn_link` block per tunnel

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | VPN Gateway Connection specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Unique on the gateway (ForceNew) |
| `vpn_gateway_id` | yes | The owning gateway (ForceNew) |
| `remote_vpn_site_id` | yes | The branch being connected (ForceNew) |
| `vpn_links` | yes (min 1) | One tunnel per site link; `vpn_site_link_id` and `bgp_enabled` are ForceNew per tunnel |
| `routing` | no | Unset = ARM's default table behavior; a configured block must name its association |
| `traffic_selector_policies` | no | CIDR pairings the tunnels are restricted to |

## Outputs

| Output | Description |
|--------|-------------|
| `connection_id` | Full ARM ID of the connection |
| `connection_name` | The connection's name |

## Usage

```hcl
module "branch_london" {
  source = "./iac/tf"

  metadata = { name = "branch-london", org = "mycompany", env = "production" }

  spec = {
    name               = "branch-london"
    vpn_gateway_id     = "/subscriptions/.../vpnGateways/hub-vpn-gateway"
    remote_vpn_site_id = "/subscriptions/.../vpnSites/branch-london"
    vpn_links = [
      {
        name             = "primary-isp"
        vpn_site_link_id = "/subscriptions/.../vpnSites/branch-london/vpnSiteLinks/primary-isp"
      }
    ]
  }
}
```

## Behavior Notes

- The connection carries no tags, region, or resource group -- ARM
  derives them through the gateway; the provider's schema has no tags
  argument.
- ARM's per-link defaults are rendered explicitly (bandwidth 10,
  protocol "IKEv2", mode "Default") so plans show the real values;
  `dpd_timeout_seconds` is omitted when unset (ARM's default is 45).
- An empty `shared_key` is emitted as null -- Azure generates a key;
  the value is sensitive and never appears in plan output.
- A tunnel provisions even when the branch device is absent -- ARM
  state is not tunnel state (provisioned-is-not-connected).

## Required Permissions

The deploying credential needs
`Microsoft.Network/vpnGateways/vpnConnections/write` plus read on the
VPN site -- held via Network Contributor, Contributor, or Owner.
