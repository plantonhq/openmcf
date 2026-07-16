# Create the load balancer with its frontend IP configurations.
#
# Lifecycle notes worth knowing before operating this resource:
# - SKU, SKU tier, and edge zone are fixed at creation -- changing any of
#   them replaces the load balancer. Frontends, pools, probes, and rules
#   all update in place.
# - Azure does not allow removing ALL frontends from an existing load
#   balancer; going from some to none replaces the resource.
# - Changing a frontend's zones replaces that frontend (and briefly its
#   address) -- pick the zone posture up front.
resource "azurerm_lb" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  sku                 = local.sku
  sku_tier            = local.sku_tier
  edge_zone           = var.spec.edge_zone
  tags                = local.final_tags

  # Each frontend is public (public IP / prefix) or internal (subnet).
  # The allocation is derived: a pinned private address means Static, an
  # internal frontend without a pin stays Dynamic, and public frontends
  # carry no allocation at all.
  dynamic "frontend_ip_configuration" {
    for_each = local.frontend_ip_configurations
    content {
      name                          = frontend_ip_configuration.value.name
      subnet_id                     = frontend_ip_configuration.value.subnet_id
      public_ip_address_id          = frontend_ip_configuration.value.public_ip_address_id
      public_ip_prefix_id           = frontend_ip_configuration.value.public_ip_prefix_id
      private_ip_address            = frontend_ip_configuration.value.private_ip_address
      private_ip_address_allocation = frontend_ip_configuration.value.allocation
      private_ip_address_version    = frontend_ip_configuration.value.version
      zones                         = length(frontend_ip_configuration.value.zones) > 0 ? frontend_ip_configuration.value.zones : null

      gateway_load_balancer_frontend_ip_configuration_id = frontend_ip_configuration.value.gateway_lb_fip_id
    }
  }
}

# Backend address pools. NIC-based membership joins from the member side
# (via the exported pool IDs); only vnet-scoped IP-based membership is
# declared here. Tunnel interfaces exist only on GATEWAY-SKU pools (the
# NVA chaining contract -- spec-level validation pairs them with the SKU).
resource "azurerm_lb_backend_address_pool" "pools" {
  for_each = { for pool in var.spec.backend_pools : pool.name => pool }

  name               = each.value.name
  loadbalancer_id    = azurerm_lb.main.id
  virtual_network_id = each.value.virtual_network_id
  # The tfvars wire format carries the FULL proto value name; map it to
  # ARM's casing (an unmapped pass-through would fail the azurerm
  # provider's validation).
  synchronous_mode = each.value.synchronous_mode != null ? local.synchronous_mode_map[each.value.synchronous_mode] : null

  dynamic "tunnel_interface" {
    for_each = each.value.tunnel_interfaces
    content {
      identifier = tunnel_interface.value.identifier
      port       = tunnel_interface.value.port
      protocol   = local.tunnel_protocol_map[tunnel_interface.value.protocol]
      type       = local.tunnel_type_map[tunnel_interface.value.type]
    }
  }
}

# IP-based backend members: appliances/servers addressed by IP (REGIONAL
# tier) or regional load balancer frontends (GLOBAL tier). NIC-based
# members never appear here -- they associate from the NIC side.
resource "azurerm_lb_backend_address_pool_address" "addresses" {
  for_each = local.pool_addresses

  name                    = each.value.name
  backend_address_pool_id = azurerm_lb_backend_address_pool.pools[each.value.pool_name].id

  # REGIONAL member: ip_address + the pool's virtual network. GLOBAL
  # member: the regional LB frontend ID instead (the two are mutually
  # exclusive; spec-level validation enforces it).
  ip_address                          = each.value.ip_address != null && each.value.ip_address != "" ? each.value.ip_address : null
  virtual_network_id                  = each.value.ip_address != null && each.value.ip_address != "" ? each.value.vnet_id : null
  backend_address_ip_configuration_id = each.value.lb_fip_id != null && each.value.lb_fip_id != "" ? each.value.lb_fip_id : null
}

# Health probes. probe_threshold is the flap dampener: consecutive
# successes required before a recovered instance is re-admitted.
resource "azurerm_lb_probe" "probes" {
  for_each = { for probe in var.spec.health_probes : probe.name => probe }

  name                = each.value.name
  loadbalancer_id     = azurerm_lb.main.id
  protocol            = lookup(local.probe_protocol_map, coalesce(each.value.protocol, "PROBE_TCP"), "Tcp")
  port                = each.value.port
  request_path        = each.value.request_path != null && each.value.request_path != "" ? each.value.request_path : null
  interval_in_seconds = each.value.interval_in_seconds
  number_of_probes    = each.value.number_of_probes
  probe_threshold     = each.value.probe_threshold
}

# Load-balancing rules. Pools and probes are resolved by name through the
# module-local maps; the frontend defaults to the sole declared frontend
# when the rule omits it (spec-level validation guarantees the omission
# only happens with exactly one frontend).
resource "azurerm_lb_rule" "rules" {
  for_each = { for rule in var.spec.rules : rule.name => rule }

  name            = each.value.name
  loadbalancer_id = azurerm_lb.main.id

  frontend_ip_configuration_name = coalesce(each.value.frontend_ip_configuration_name, local.default_frontend_name)

  protocol      = local.transport_protocol_map[each.value.protocol]
  frontend_port = each.value.frontend_port
  backend_port  = each.value.backend_port

  backend_address_pool_ids = [for n in each.value.backend_pool_names : local.backend_pool_map[n]]
  probe_id                 = each.value.probe_name != null && each.value.probe_name != "" ? local.probe_map[each.value.probe_name] : null

  load_distribution       = each.value.load_distribution != null && each.value.load_distribution != "" ? local.load_distribution_map[each.value.load_distribution] : null
  idle_timeout_in_minutes = each.value.idle_timeout_in_minutes
  floating_ip_enabled     = each.value.floating_ip_enabled
  tcp_reset_enabled       = each.value.tcp_reset_enabled
  disable_outbound_snat   = each.value.disable_outbound_snat
}

# Inbound NAT rules. Single-target mode (frontend_port) leaves attachment
# to the member side (a NIC's NAT-rule association referencing the
# exported rule ID); pool-style mode (port range + pool) gives every pool
# member its own frontend port.
resource "azurerm_lb_nat_rule" "nat_rules" {
  for_each = { for rule in var.spec.nat_rules : rule.name => rule }

  name                = each.value.name
  resource_group_name = var.spec.resource_group
  loadbalancer_id     = azurerm_lb.main.id

  frontend_ip_configuration_name = coalesce(each.value.frontend_ip_configuration_name, local.default_frontend_name)

  protocol     = local.transport_protocol_map[each.value.protocol]
  backend_port = each.value.backend_port

  # Exactly one mode is set (spec-level validation): null out the other
  # mode's fields so the provider sees a clean single-target or
  # pool-style rule.
  frontend_port           = each.value.frontend_port > 0 ? each.value.frontend_port : null
  backend_address_pool_id = each.value.backend_pool_name != null && each.value.backend_pool_name != "" ? local.backend_pool_map[each.value.backend_pool_name] : null
  frontend_port_start     = each.value.frontend_port_start > 0 ? each.value.frontend_port_start : null
  frontend_port_end       = each.value.frontend_port_end > 0 ? each.value.frontend_port_end : null

  floating_ip_enabled     = each.value.floating_ip_enabled
  tcp_reset_enabled       = each.value.tcp_reset_enabled
  idle_timeout_in_minutes = each.value.idle_timeout_in_minutes
}

# Outbound rules: explicit SNAT through public frontends. Combine with
# disable_outbound_snat on the load-balancing rules that share the pool.
resource "azurerm_lb_outbound_rule" "outbound_rules" {
  for_each = { for rule in var.spec.outbound_rules : rule.name => rule }

  name                    = each.value.name
  loadbalancer_id         = azurerm_lb.main.id
  backend_address_pool_id = local.backend_pool_map[each.value.backend_pool_name]
  protocol                = local.transport_protocol_map[each.value.protocol]

  allocated_outbound_ports = each.value.allocated_outbound_ports
  tcp_reset_enabled        = each.value.tcp_reset_enabled
  idle_timeout_in_minutes  = each.value.idle_timeout_in_minutes

  dynamic "frontend_ip_configuration" {
    for_each = each.value.frontend_ip_configuration_names
    content {
      name = frontend_ip_configuration.value
    }
  }
}
