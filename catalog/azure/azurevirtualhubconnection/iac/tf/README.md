# AzureVirtualHubConnection Terraform Module

## Overview

This Terraform module provisions a Virtual Hub Connection using the
`azurerm` provider. It creates a single `azurerm_virtual_hub_connection`
-- the attachment joining one spoke virtual network to a Virtual WAN
hub, with the routing block where WAN topologies (isolation, shared
services, service chaining) are expressed.

## Resources Created

- `azurerm_virtual_hub_connection.main` -- the attachment

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Virtual Hub Connection specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | 2-80 chars, the provider's own regex; ForceNew |
| `virtual_hub_id` | yes | The hub being joined (ForceNew) |
| `remote_virtual_network_id` | yes | The spoke VNet being attached (ForceNew) |
| `internet_security_enabled` | no | Off by default (the spoke keeps its own egress) |
| `routing` | no | Association, propagation, static routes; unset = ARM's any-to-any default |

## Outputs

| Output | Description |
|--------|-------------|
| `virtual_hub_connection_id` | Full ARM ID -- what a hub BGP peering references |
| `virtual_hub_connection_name` | The connection's name |

## Usage

```hcl
module "spoke_app" {
  source = "./iac/tf"

  metadata = { name = "spoke-app", org = "mycompany", env = "production" }

  spec = {
    name                      = "spoke-app"
    virtual_hub_id            = "/subscriptions/.../virtualHubs/hub-eastus"
    remote_virtual_network_id = "/subscriptions/.../virtualNetworks/app-vnet"
  }
}
```

## Behavior Notes

- The connection carries no tags -- the provider's schema has none.
- A configured routing block must configure something (association,
  propagation, or static routes) -- the spec enforces the provider's
  at-least-one-of rule upfront.
- `static_vnet_local_route_override_criteria` is fixed once created
  (ARM replaces the connection to change it).
- The hub cannot be deleted while this connection exists.

## Required Permissions

The deploying credential needs
`Microsoft.Network/virtualHubs/hubVirtualNetworkConnections/write` plus
join permission on the spoke VNet
(`Microsoft.Network/virtualNetworks/peer/action` class) -- held via
Network Contributor, Contributor, or Owner on both sides.
