# AzureExpressRouteCircuit Terraform Module

## Overview

This Terraform module provisions an ExpressRoute circuit using the
`azurerm` provider. It creates an `azurerm_express_route_circuit` (the
billing/identity object for a dedicated private connection to
Microsoft) plus one `azurerm_express_route_circuit_authorization` per
spec `authorizations` entry -- the keys other subscriptions redeem to
connect their gateways to this circuit.

The provisioning-mode contract (service-provider trio XOR ExpressRoute
Direct pair) and SKU requirements are spec-validated before the module
runs. **Billing starts at creation** -- Azure meters the circuit from
service-key issuance, even while the provider side is unprovisioned.

## Resources Created

- `azurerm_express_route_circuit.main` -- the circuit
- `azurerm_express_route_circuit_authorization.authorizations` -- one
  per spec entry, keyed by name (`for_each`)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | ExpressRoute circuit specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` / `resource_group` / `name` | yes | The circuit's ARM identity; all ForceNew |
| `sku_tier` / `sku_family` | yes | BASIC/LOCAL/STANDARD/PREMIUM x METERED_DATA/UNLIMITED_DATA |
| `service_provider_name` + `peering_location` + `bandwidth_in_mbps` | mode 1 | The provider trio (bandwidth grows in place, never shrinks) |
| `express_route_port_id` + `bandwidth_in_gbps` | mode 2 | The ExpressRoute Direct pair |
| `rate_limiting_enabled` | no | Direct circuits only (ARM EnableDirectPortRateLimit) |
| `authorization_key` | no | The key this circuit REDEEMS (sensitive) |
| `authorizations` | no | Keys this circuit ISSUES, by name |
| `tags` | no | User tags, merged over metadata-derived tags |

## Outputs

| Output | Description |
|--------|-------------|
| `express_route_circuit_id` | Full ARM ID |
| `express_route_circuit_name` | The join key peerings reference |
| `service_key` | The provisioning credential for the provider (sensitive) |
| `service_provider_provisioning_state` | NotProvisioned / Provisioning / Provisioned / Deprovisioning |
| `authorization_keys` | Name-keyed generated keys (sensitive) |

## Usage

```hcl
module "hq_circuit" {
  source = "./iac/tf"

  metadata = { name = "hq-circuit", org = "mycompany", env = "production" }

  spec = {
    region                = "eastus"
    resource_group        = "network-rg"
    name                  = "hq-circuit"
    sku_tier              = "STANDARD"
    sku_family            = "METERED_DATA"
    service_provider_name = "Equinix"
    peering_location      = "Washington DC"
    bandwidth_in_mbps     = 1000
  }
}
```

## Behavior Notes

- A fresh circuit sits in `NotProvisioned` until the connectivity
  provider completes the cross-connect with the service key --
  peerings cannot be configured before that.
- The provider sets the redeemed `authorization_key` in a SECOND
  create call after the circuit exists (an ARM sequencing quirk).
- Decreasing `bandwidth_in_mbps` replaces the circuit; increasing it
  updates in place.

## Required Permissions

The deploying credential needs
`Microsoft.Network/expressRouteCircuits/write` (and
`.../authorizations/write` when issuing keys) -- held via Network
Contributor, Contributor, or Owner.
