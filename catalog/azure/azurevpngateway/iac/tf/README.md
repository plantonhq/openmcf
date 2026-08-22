# AzureVpnGateway Terraform Module

## Overview

This Terraform module provisions a Virtual WAN VPN Gateway using the
`azurerm` provider. It creates an `azurerm_vpn_gateway` (the managed
site-to-site terminator inside a virtual hub) plus one
`azurerm_vpn_gateway_nat_rule` per spec entry -- the composed children
tunnels opt into for overlapping-address translation.

## Resources Created

- `azurerm_vpn_gateway.main` -- the gateway (30-45 min create; bills
  from creation)
- `azurerm_vpn_gateway_nat_rule.nat_rules` -- one per spec `nat_rules`
  entry, keyed by rule name

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | VPN Gateway specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Must match the hub's region (ForceNew) |
| `resource_group` | yes | Resource group name (ForceNew) |
| `name` | yes | The gateway's name (ForceNew) |
| `virtual_hub_id` | yes | The hub to deploy into -- one VPN gateway per hub (ForceNew) |
| `routing_preference` | no | `MICROSOFT_NETWORK` (default) or `INTERNET` (ForceNew) |
| `scale_unit` | no | 500 Mbps units across the pair; default 1 (updates in place) |
| `bgp_settings` | no | asn/peer_weight (ForceNew) + per-instance custom APIPA IPs (update in place) |
| `nat_rules` | no | Composed children; each needs external AND internal mappings |

## Outputs

| Output | Description |
|--------|-------------|
| `vpn_gateway_id` | Full ARM ID -- what a connection's `vpn_gateway_id` references |
| `vpn_gateway_name` | The gateway's name |
| `bgp_asn` | The gateway's ASN (65515 on today's Virtual WAN) |
| `public_ip_addresses` | Each instance's public IPv4 -- what branch devices dial |
| `private_ip_addresses` | Each instance's private IPv4 |
| `nat_rule_ids` | Each NAT rule's ARM ID keyed by rule name |

## Usage

```hcl
module "hub_vpn_gateway" {
  source = "./iac/tf"

  metadata = { name = "hub-vpn-gateway", org = "mycompany", env = "production" }

  spec = {
    name           = "hub-vpn-gateway"
    region         = "eastus"
    resource_group = "network-rg"
    virtual_hub_id = "/subscriptions/.../virtualHubs/hub-eastus"
  }
}
```

## Behavior Notes

- ARM's defaults are rendered explicitly (routing preference
  "Microsoft Network" -- note the space -- and scale unit 1) so plans
  show the real values.
- The BGP block's asn/peer_weight are ForceNew; the custom APIPA
  addresses are applied by the provider in a SECOND call after the
  gateway exists, which is why they update in place.
- NAT rule mode/type default to "EgressSnat"/"Static" explicitly; an
  unspecified instance pin emits null (the rule applies on both
  instances).
- Deleting the gateway requires its connections to be gone first.
