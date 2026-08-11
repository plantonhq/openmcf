# AzureExpressRouteCircuitPeering Terraform Module

## Overview

This Terraform module provisions an ExpressRoute circuit peering using
the `azurerm` provider. It creates an
`azurerm_express_route_circuit_peering` (the BGP routing configuration
that makes routes flow through a circuit) plus one
`azurerm_express_route_circuit_connection` per spec `connections` entry
-- Global Reach links to other circuits' private peerings.

The type-dependent contracts (route filter and Microsoft config only on
MICROSOFT_PEERING, no IPv6 on the deprecated public peering,
connections only on private peering, the /30 pair travelling together)
are spec-validated before the module runs.

## Resources Created

- `azurerm_express_route_circuit_peering.main` -- the peering (at most
  one of each type per circuit; the type IS the ARM child name)
- `azurerm_express_route_circuit_connection.connections` -- one per
  spec entry, keyed by name (`for_each`)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Peering specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `resource_group` / `express_route_circuit_name` | yes | The parent circuit (by NAME); ForceNew |
| `peering_type` | yes | AZURE_PRIVATE_PEERING / MICROSOFT_PEERING (/ deprecated public) |
| `vlan_id` | yes | 1-4094, unique on the circuit |
| `primary_peer_address_prefix` + `secondary_peer_address_prefix` | pair | One /30 per physical link |
| `microsoft_peering_config` | MS only | The advertised-public-prefix contract |
| `ipv6` | no | /126 pairs + its own MS contract (not on public peering) |
| `route_filter_id` | MS only | Selects Microsoft service communities |
| `shared_key` | no | BGP MD5 key (sensitive; never read back) |
| `connections` | private only | Global Reach links to other circuits' private peerings |

## Outputs

| Output | Description |
|--------|-------------|
| `express_route_circuit_peering_id` | Full ARM ID -- the far side of another circuit's Global Reach connection |
| `azure_asn` | Microsoft's BGP ASN on the peering |
| `primary_azure_port` / `secondary_azure_port` | Microsoft-edge port identifiers |
| `connection_ids` | Name-keyed Global Reach connection IDs |

## Usage

```hcl
module "hq_private_peering" {
  source = "./iac/tf"

  metadata = { name = "hq-private-peering", org = "mycompany", env = "production" }

  spec = {
    resource_group              = "network-rg"
    express_route_circuit_name  = "hq-circuit"
    peering_type                = "AZURE_PRIVATE_PEERING"
    vlan_id                     = 100
    primary_peer_address_prefix   = "192.168.16.0/30"
    secondary_peer_address_prefix = "192.168.16.4/30"
  }
}
```

## Behavior Notes

- ARM accepts and stores peering configuration even on a fresh
  (NotProvisioned) circuit -- live-verified for private peering. The
  BGP session only establishes once the connectivity provider completes
  the cross-connect; Microsoft peering's server-side public-prefix
  validation has not been exercised on an unprovisioned circuit.
- `shared_key` is write-only in ARM; the provider never reads it back.
- The provider serializes all peering operations per circuit through an
  internal lock -- concurrent peerings on one circuit apply in
  sequence.

## Required Permissions

The deploying credential needs
`Microsoft.Network/expressRouteCircuits/peerings/write` (and
`.../connections/write` for Global Reach) -- held via Network
Contributor, Contributor, or Owner.
