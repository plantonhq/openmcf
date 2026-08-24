# AzureVirtualHub Terraform Module

## Overview

This Terraform module provisions a Virtual Hub using the `azurerm`
provider: the managed regional router of a Virtual WAN, plus its
composed routing children -- custom route tables (with inline routes),
route maps, BGP connections with NVAs, and the hub's routing intent.

## Resources Created

- `azurerm_virtual_hub.main` -- the hub
- `azurerm_virtual_hub_route_table.route_tables` -- one per
  `spec.route_tables` entry, keyed by name (routes managed INLINE;
  never mix in the provider's standalone route resource)
- `azurerm_route_map.route_maps` -- one per `spec.route_maps` entry,
  keyed by name
- `azurerm_virtual_hub_bgp_connection.bgp_connections` -- one per
  `spec.bgp_connections` entry, keyed by name
- `azurerm_virtual_hub_routing_intent.routing_intent` -- at most one,
  keyed by the intent's name

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Virtual Hub specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The hub's ARM identity; all ForceNew |
| `virtual_wan_id` | yes | The WAN the hub belongs to (WAN hubs only; ForceNew) |
| `address_prefix` | yes | The hub's private CIDR (min /24, recommend /23; ForceNew) |
| `sku` | no | STANDARD (default) or BASIC (ForceNew) |
| `hub_routing_preference` | no | EXPRESS_ROUTE (default) / VPN_GATEWAY / AS_PATH |
| `virtual_router_auto_scale_min_capacity` | no | Router capacity floor; ARM's default is 2 |
| `routes` | no | Classic inline routes on the default table |
| `route_tables` / `route_maps` / `bgp_connections` / `routing_intent` | no | The composed routing children |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `virtual_hub_id` | Full ARM ID -- what connections and gateways reference |
| `virtual_hub_name` | The hub's name |
| `default_route_table_id` | The built-in default route table's ARM ID |
| `virtual_router_asn` | The hub router's BGP ASN (always 65515) |
| `virtual_router_ips` | The hub router's peering IPv4 addresses |
| `route_table_ids` | Custom route table ARM IDs, keyed by name |
| `route_map_ids` | Route map ARM IDs, keyed by name |
| `bgp_connection_ids` | BGP connection ARM IDs, keyed by name |
| `routing_intent_id` | The routing intent's ARM ID ("" when none) |

## Usage

```hcl
module "hub_eastus" {
  source = "./iac/tf"

  metadata = { name = "hub-eastus", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "hub-eastus"
    virtual_wan_id = "/subscriptions/.../virtualWans/global-wan"
    address_prefix = "10.100.0.0/23"
  }
}
```

## Behavior Notes

- A Standard hub bills hourly from creation; ARM takes
  15-30 minutes to bring the router to a Provisioned routing state.
- Routing intent and per-connection route-table customization are
  mutually exclusive on ARM's side.
- ARM refuses to delete a hub that still has gateways or connections --
  destroy those first.
