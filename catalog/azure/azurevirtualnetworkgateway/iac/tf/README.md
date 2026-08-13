# AzureVirtualNetworkGateway Terraform Module

## Overview

This Terraform module provisions an Azure virtual network gateway using
the `azurerm` provider. It creates an `azurerm_virtual_network_gateway`
-- the managed appliance terminating site-to-site VPN, point-to-site
VPN, VNet-to-VNet, or ExpressRoute connectivity -- plus one
`azurerm_virtual_network_gateway_nat_rule` per composed NAT rule.

The gateway is one third of the site-to-site story: it lives in the
referenced "GatewaySubnet" (an ARM name contract) and binds referenced
AzurePublicIp addresses; AzureLocalNetworkGateway describes each
on-premises site and AzureVirtualNetworkGatewayConnection ties a site to
this gateway. Gateways provision in 25-45 minutes and delete in 10-20,
so the ForceNew surface (name, region, type, vpn_type, generation, edge
zone, private-IP enablement, every ip configuration) is expensive --
design changes to avoid replacement.

Enum fields arrive as proto enum value names (VPN_GW_1_AZ,
EXPRESS_ROUTE, EGRESS_SNAT); `locals.tf` maps them to azurerm's exact
vocabulary. The
type/generation/SKU pairing rules and the public-IP contract (required
per configuration on VPN gateways, forbidden on ExpressRoute) are
spec-validated before the module ever runs.

## Resources Created

- `azurerm_virtual_network_gateway.main` -- the gateway itself
- `azurerm_virtual_network_gateway_nat_rule.rules` -- one per `nat_rules` entry (keyed by rule name)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Virtual network gateway specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The gateway's ARM identity; all ForceNew |
| `type` | no | `VPN` (default) or `EXPRESS_ROUTE`; flips the public-IP contract |
| `vpn_type` | no | `ROUTE_BASED` (default) or `POLICY_BASED` (legacy, BASIC only) |
| `sku` | yes | Full SKU vocabulary; pairing rules spec-enforced |
| `generation` | no | `GENERATION1`/`GENERATION2`; omitted lets Azure pick |
| `ip_configurations` | yes | 1-3 bindings of the GatewaySubnet + (VPN) public IPs |
| `active_active` | no | Two-instance pair; needs two ip configurations |
| `bgp_enabled` / `bgp_settings` | no | BGP speaker: ASN, peer weight, APIPA peering addresses |
| `vpn_client_configuration` | no | Point-to-site: pool, Entra ID/cert/RADIUS auth, protocols |
| `nat_rules` | no | Composed NAT rules; ids surface in `nat_rule_ids` |
| `minimum_scale_unit` / `maximum_scale_unit` | no | ER_GW_SCALE autoscale bounds (this engine only) |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `virtual_network_gateway_id` | Full ARM ID of the gateway -- the join key connections attach through |
| `virtual_network_gateway_name` | The gateway's name as deployed |
| `nat_rule_ids` | Map of NAT rule name to ARM id -- connections opt in via their egress/ingress lists |

## Usage

```hcl
module "vpn_gateway" {
  source = "./iac/tf"

  metadata = { name = "hub-vpn", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "network-rg"
    name           = "hub-vpn-gateway"
    sku            = "VPN_GW_1_AZ"
    ip_configurations = [{
      subnet_id            = "/subscriptions/xxx/.../subnets/GatewaySubnet"
      public_ip_address_id = "/subscriptions/xxx/.../publicIPAddresses/vpn-gw-pip"
    }]
    bgp_enabled  = true
    bgp_settings = { asn = 65515 }
  }
}
```

## Behavior Notes

- **PARITY-EXCEPTION (this engine only)**: `minimum_scale_unit` /
  `maximum_scale_unit` (the ER_GW_SCALE autoscaler) are not expressible
  in the Pulumi classic SDK -- the Pulumi module fails loudly when they
  are set, so autoscaling ExpressRoute gateways deploy via Terraform.
- `dns_forwarding_enabled` is sent only when true: ARM rejects the
  parameter on SKUs/types without DNS forwarding support.
- `ip_sec_replay_protection_enabled` is omitted when unset -- the
  provider's default (on) matches the spec's documented default.
- An empty ip configuration name becomes `vnetGatewayConfig`, the name
  the Azure portal uses, on both engines.

## Required Permissions

The deploying credential needs
`Microsoft.Network/virtualNetworkGateways/*` on the resource group,
plus join rights on the gateway subnet and public IPs -- held via
Network Contributor, Contributor, or Owner.
