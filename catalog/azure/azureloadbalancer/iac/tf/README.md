# AzureLoadBalancer Terraform Module

## Overview

This Terraform module provisions an Azure Load Balancer using the
`azurerm` provider (`~> 5.0`). It creates the `azurerm_lb` composite --
frontends, backend address pools, health probes, load-balancing rules,
inbound NAT rules, and outbound (SNAT) rules -- as one unit, because
every sub-resource shares the load balancer's lifecycle.

What is deliberately NOT here is pool MEMBERSHIP: a network interface or
scale set joins a pool from the member side, referencing this module's
name-keyed `backend_pool_ids` output (and a NIC's NAT-rule association
references `nat_rule_ids`). Only IP-based members -- appliances or, for
GLOBAL-tier pools, regional load balancer frontends -- are declared
inline via each pool's `addresses`.

Spec-level validation enforces the same rules ARM does -- the frontend
address-source exclusivity, the SKU pairings (GATEWAY/tunnels,
GLOBAL/STANDARD, two-pool rules), the HA-ports zero-port pairing, the
NAT-rule mode XOR, the probe request-path pairing, and every by-name
cross-reference -- so the module maps fields without re-validating them.

Lifecycle notes worth knowing before operating this resource:

- SKU, SKU tier, and edge zone are fixed at creation -- changing any of
  them replaces the load balancer. Frontends, pools, probes, and rules
  all update in place.
- Azure does not allow removing ALL frontends from an existing load
  balancer; going from some to none replaces the resource.
- Changing a frontend's `zones` replaces that frontend (and briefly its
  address) -- pick the zone posture up front.

## Resources Created

- `azurerm_lb.main` -- the load balancer with its frontend IP configurations
- `azurerm_lb_backend_address_pool.pools` -- one per backend pool
- `azurerm_lb_backend_address_pool_address.addresses` -- one per inline IP-based pool member
- `azurerm_lb_probe.probes` -- one per health probe
- `azurerm_lb_rule.rules` -- one per load-balancing rule
- `azurerm_lb_nat_rule.nat_rules` -- one per inbound NAT rule (single-target or pool-style)
- `azurerm_lb_outbound_rule.outbound_rules` -- one per outbound rule

The legacy `azurerm_lb_nat_pool` is not modeled: pool-style NAT rules
(`frontend_port_start`/`_end` + `backend_pool_name`) supersede it.

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Load balancer specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; must match the backend resources' region |
| `resource_group` | yes | Resource group name |
| `name` | yes | Load balancer name, unique within the resource group |
| `sku` | no | `STANDARD` (default) or `GATEWAY` (NVA chaining); fixed at creation. Basic is not modeled (retired September 2025) |
| `sku_tier` | no | `REGIONAL` (default) or `GLOBAL` (cross-region; STANDARD only); fixed at creation |
| `edge_zone` | no | Azure Edge Zone pinning; fixed at creation |
| `frontend_ip_configurations` | yes | At least one. Each is public (`public_ip_address_id` or `public_ip_prefix_id`) or internal (`subnet_id` + optional pinned `private_ip_address`, `private_ip_address_version`, `zones`); `gateway_load_balancer_frontend_ip_configuration_id` chains it behind a Gateway LB |
| `backend_pools` | no | Named pools; `virtual_network_id` + `synchronous_mode` + inline `addresses` for IP-based membership, `tunnel_interfaces` on GATEWAY pools. NIC membership joins member-side |
| `health_probes` | no | `protocol` (`PROBE_TCP`/`PROBE_HTTP`/`PROBE_HTTPS`), `port`, `request_path` (HTTP/HTTPS only), `interval_in_seconds` (default 15), `number_of_probes` (default 2), `probe_threshold` (default 1) |
| `rules` | no | Frontend port/protocol to pool/port mappings: `backend_pool_names` (two only on GATEWAY), `probe_name`, `load_distribution`, `idle_timeout_in_minutes`, `floating_ip_enabled`, `tcp_reset_enabled`, `disable_outbound_snat`. Protocol `ALL` + ports 0 = HA ports |
| `nat_rules` | no | Single-target mode (`frontend_port`; the NIC association completes it) XOR pool-style mode (`backend_pool_name` + `frontend_port_start`/`_end`); plus `backend_port`, `floating_ip_enabled`, `tcp_reset_enabled`, `idle_timeout_in_minutes` (4-30) |
| `outbound_rules` | no | Explicit SNAT: `frontend_ip_configuration_names` (public frontends), `backend_pool_name`, `protocol`, `allocated_outbound_ports` (default 1024; 0 = divide evenly), `tcp_reset_enabled`, `idle_timeout_in_minutes` |
| `tags` | no | User tags, merged over metadata-derived tags (user wins) |

Enum-typed fields arrive as the FULL proto value name and are mapped to
ARM values internally (`STANDARD` becomes `Standard`, `PROBE_HTTP`
becomes `Http`, `SOURCE_IP_PROTOCOL` becomes `SourceIPProtocol`, `VXLAN`
stays `VXLAN`); unset optional enums apply Azure's defaults so an
unspecified spec and Azure's default deploy identically.

Rules and NAT rules may omit `frontend_ip_configuration_name` only when
exactly one frontend is declared -- the module defaults it to that
frontend. Pools and probes are resolved by name through module-local
maps.

## Outputs

| Output | Description |
|--------|-------------|
| `load_balancer_id` | Full ARM ID of the load balancer |
| `load_balancer_name` | The load balancer's name as deployed |
| `private_ip_address` | The first internal frontend's private IP (empty when every frontend is public) |
| `private_ip_addresses` | All internal frontends' private IPs, in declaration order |
| `frontend_ip_configuration_ids` | ARM ID of each frontend, keyed by frontend name -- referenced for gateway chaining and GLOBAL-tier pool members |
| `backend_pool_ids` | ARM ID of each backend pool, keyed by pool name -- THE member-side association seam (NIC and VMSS references) |
| `probe_ids` | ARM ID of each health probe, keyed by probe name -- referenced by a scale set's rolling-upgrade health probe |
| `nat_rule_ids` | ARM ID of each inbound NAT rule, keyed by rule name -- the NIC NAT-rule association seam |

## Usage

```hcl
module "load_balancer" {
  source = "./iac/tf"

  metadata = { name = "web-lb", org = "mycompany", env = "production" }

  spec = {
    region         = "eastus"
    resource_group = "prod-rg"
    name           = "web-lb"

    frontend_ip_configurations = [{
      name                 = "public"
      public_ip_address_id = "/subscriptions/xxx/.../publicIPAddresses/web-lb-ip"
    }]

    backend_pools = [
      { name = "web" }
    ]

    health_probes = [{
      name         = "http-health"
      protocol     = "PROBE_HTTP"
      port         = 8080
      request_path = "/healthz"
    }]

    rules = [{
      name               = "https"
      protocol           = "TCP"
      frontend_port      = 443
      backend_port       = 8443
      backend_pool_names = ["web"]
      probe_name         = "http-health"
      tcp_reset_enabled  = true
    }]
  }
}
```

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege
action manifest the deploying credential needs.
