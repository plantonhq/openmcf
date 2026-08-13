# AzureExpressRouteGateway Terraform Module

## Overview

This Terraform module provisions an ExpressRoute Gateway using the
`azurerm` provider: the Virtual WAN on-ramp for ExpressRoute circuits,
plus its composed connections -- one ARM child per spec entry, each
joining a circuit's private peering to the hub.

## Resources Created

- `azurerm_express_route_gateway.main` -- the gateway
- `azurerm_express_route_connection.connections` -- one per
  `spec.connections` entry, keyed by name

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | ExpressRoute Gateway specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The gateway's ARM identity; all ForceNew |
| `virtual_hub_id` | yes | The hub the gateway deploys into (ForceNew; one gateway per hub) |
| `scale_units` | yes | The capacity floor (1-10, ~2 Gbps each); updatable in place |
| `allow_non_virtual_wan_traffic` | no | Off by default (WAN networks only) |
| `connections` | no | The composed circuit-peering connections |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `express_route_gateway_id` | Full ARM ID of the gateway |
| `express_route_gateway_name` | The gateway's name |
| `connection_ids` | Connection ARM IDs, keyed by name |

## Usage

```hcl
module "hub_er_gateway" {
  source = "./iac/tf"

  metadata = { name = "hub-er-gateway", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "hub-er-gateway"
    virtual_hub_id = "/subscriptions/.../virtualHubs/hub-eastus"
    scale_units    = 1
  }
}
```

## Behavior Notes

- The gateway bills ~$0.42/hr per scale unit FROM CREATION and takes
  roughly 30 minutes to provision.
- ARM accepts a connection only when the circuit's provider side is
  PROVISIONED -- a rejection on provisioning state is a prerequisite
  problem, not a module defect.
- `authorization_key` is sensitive and never returned by ARM -- an
  imported connection configured with one legitimately plans an
  in-place update on it.
- Deletion is bottom-up: connections → gateway → hub.

## Required Permissions

The deploying credential needs
`Microsoft.Network/expressRouteGateways/write` (plus
`expressRouteGateways/expressRouteConnections/write` for connections)
-- held via Network Contributor, Contributor, or Owner.
