# AzureExpressRoutePort Terraform Module

## Overview

This Terraform module provisions an ExpressRoute Port using the
`azurerm` provider. It creates an `azurerm_express_route_port` (your
own physical port pair on a Microsoft edge router -- ExpressRoute
Direct) plus one `azurerm_express_route_port_authorization` per spec
`authorizations` entry -- the keys other subscriptions redeem to build
circuits on this port's capacity.

**The port bills its full monthly rate from creation** (one of the most
expensive single objects in Azure networking), and some subscriptions
need Microsoft enrollment for ExpressRoute Direct before ARM accepts
the create. The MACsec contracts (keys travel together; a user-assigned
identity must be present) are spec-validated before the module runs.

## Resources Created

- `azurerm_express_route_port.main` -- the port pair
- `azurerm_express_route_port_authorization.authorizations` -- one per
  spec entry, keyed by name (`for_each`)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | ExpressRoute Port specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The port's ARM identity; all ForceNew |
| `peering_location` | yes | The colocation facility (ExpressRoute Direct vocabulary); ForceNew |
| `bandwidth_in_gbps` | yes | 10 or 100 at real locations; ForceNew |
| `encapsulation` | yes | DOT1Q or QINQ; ForceNew |
| `billing_type` | no | METERED_DATA (default) or UNLIMITED_DATA |
| `identity` | no | Managed identity; USER_ASSIGNED required for MACsec |
| `link1` / `link2` | no | Manipulate the fixed physical pair (admin state, MACsec) |
| `authorizations` | no | Keys this port ISSUES, by name |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `express_route_port_id` | Full ARM ID -- what Direct circuits reference |
| `express_route_port_name` | The port's name |
| `guid` / `ethertype` / `mtu` | Port-level physical facts |
| `system_assigned_identity_principal_id` | Populated when the identity type includes SYSTEM_ASSIGNED |
| `link{1,2}_id` / `_router_name` / `_interface_name` / `_patch_panel_id` / `_rack_id` / `_connector_type` | The per-link letter-of-authorization facts for the facility's cross-connect order |
| `authorization_keys` | Name-keyed generated keys (sensitive) |

## Usage

```hcl
module "hq_port" {
  source = "./iac/tf"

  metadata = { name = "hq-port", org = "mycompany", env = "production" }

  spec = {
    region            = "eastus"
    resource_group    = "network-rg"
    name              = "hq-port"
    peering_location  = "Equinix-Ashburn-DC2"
    bandwidth_in_gbps = 10
    encapsulation     = "DOT1Q"
    link1             = { admin_enabled = true }
    link2             = { admin_enabled = true }
  }
}
```

## Behavior Notes

- ARM always creates the link PAIR with the port; the `link1`/`link2`
  blocks only manipulate the existing links. The provider applies link
  configuration in a second call after the port exists, so a fresh
  port briefly reports links disabled even when `admin_enabled` is
  true.
- Links start administratively DISABLED by default -- enable them once
  the facility completes the physical cross-connect.
- The provider serializes the port and its authorization children with
  its own lock (ARM allows one port mutation at a time).
