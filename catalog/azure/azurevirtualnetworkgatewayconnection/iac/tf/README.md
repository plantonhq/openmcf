# AzureVirtualNetworkGatewayConnection Terraform Module

## Overview

This Terraform module provisions a gateway connection using the
`azurerm` provider. It creates an
`azurerm_virtual_network_gateway_connection` -- the tunnel object
joining a virtual network gateway to an on-premises device (IPsec),
another gateway (Vnet2Vnet), or an ExpressRoute circuit.

The connection provisions in seconds; its far-side pairing rules
(IPSEC needs the local network gateway, VNET_TO_VNET the peer gateway,
EXPRESS_ROUTE the circuit) are spec-validated before the module runs.
Enum fields arrive as proto enum value names (IPSEC, IKE_V2, DEFAULT);
`locals.tf` maps them to azurerm's exact vocabulary. ARM provisioning
success does NOT mean the tunnel is Connected -- the far side must
negotiate.

## Resources Created

- `azurerm_virtual_network_gateway_connection.main` -- the connection itself

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Gateway connection specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The connection's ARM identity; all ForceNew |
| `type` | yes | `IPSEC`, `VNET_TO_VNET`, or `EXPRESS_ROUTE`; decides the required far side |
| `virtual_network_gateway_id` | yes | The owning gateway (ForceNew) |
| `local_network_gateway_id` | IPSEC | The on-premises site's description |
| `peer_virtual_network_gateway_id` | VNET_TO_VNET | The far gateway (mirror connection required there) |
| `express_route_circuit_id` | EXPRESS_ROUTE | The circuit (ForceNew) |
| `shared_key` | no | Sensitive; omitted lets Azure generate one |
| `ipsec_policy` | no | Pinned IPsec/IKE proposal (six algorithms together) |
| `bgp_enabled` / `custom_bgp_addresses` | no | Dynamic routing; APIPA endpoints per tunnel |
| `egress_nat_rule_ids` / `ingress_nat_rule_ids` | no | Gateway NAT rules this tunnel opts into |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `connection_id` | Full ARM ID of the connection |
| `connection_name` | The connection's name as deployed |

## Usage

```hcl
module "site_to_site" {
  source = "./iac/tf"

  metadata = { name = "hq-to-azure", org = "mycompany", env = "production" }

  spec = {
    region                     = "eastus"
    resource_group             = "network-rg"
    name                       = "hq-to-azure"
    type                       = "IPSEC"
    virtual_network_gateway_id = "/subscriptions/xxx/.../virtualNetworkGateways/hub-vpn"
    local_network_gateway_id   = "/subscriptions/xxx/.../localNetworkGateways/hq"
    shared_key                 = var.tunnel_psk
  }
}
```

## Behavior Notes

- **PARITY-EXCEPTION (this engine renders every selector)**: the Pulumi
  classic SDK models exactly one `traffic_selector_policies` entry --
  multi-selector connections deploy via this engine only.
- `shared_key` and `authorization_key` are omitted when empty: Azure
  generates the pre-shared key when absent (readable back via the
  shared-key API).
- `connection_protocol` is omitted when unset so Azure applies its
  default (IKEv2); the provider treats it as Computed.
