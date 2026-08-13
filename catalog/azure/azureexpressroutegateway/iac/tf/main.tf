# Create the ExpressRoute Gateway -- the Virtual WAN on-ramp for
# ExpressRoute circuits. The gateway bills per scale unit (~$0.42/hr
# per unit) FROM CREATION, and ARM takes roughly 30 minutes to
# provision one. A hub holds at most one ExpressRoute gateway.
resource "azurerm_express_route_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  virtual_hub_id      = var.spec.virtual_hub_id

  # The MINIMUM scale units (1-10, spec-validated); each unit carries
  # ~2 Gbps and ARM auto-scales above this floor under load.
  scale_units = var.spec.scale_units

  # Off is ARM's default: only Virtual WAN networks may ride the
  # gateway. On, classic VNets connected to the same circuit may too.
  allow_non_virtual_wan_traffic = var.spec.allow_non_virtual_wan_traffic

  tags = local.final_tags
}

# The composed connections: standalone ARM children of the gateway, one
# per spec entry, keyed by name. Each joins one ExpressRoute circuit
# PEERING to the hub -- and ARM accepts it only when the circuit's
# provider side is PROVISIONED (a live carrier or an ExpressRoute
# Direct port behind it).
resource "azurerm_express_route_connection" "connections" {
  for_each = { for connection in var.spec.connections : connection.name => connection }

  name                             = each.value.name
  express_route_gateway_id         = azurerm_express_route_gateway.main.id
  express_route_circuit_peering_id = each.value.express_route_circuit_peering_id

  # The authorization key (a UUID, sensitive) for a circuit in ANOTHER
  # subscription; empty means the circuit is in this subscription.
  authorization_key = each.value.authorization_key != "" ? each.value.authorization_key : null

  internet_security_enabled            = each.value.internet_security_enabled
  express_route_gateway_bypass_enabled = each.value.express_route_gateway_bypass_enabled

  # 0 is ARM's default; higher weight wins when the same prefix is
  # reachable over multiple connections.
  routing_weight = each.value.routing_weight

  # Unset routing applies ARM's default behavior: associate with and
  # propagate to the hub's built-in default route table. The spec
  # guarantees a configured block carries an association or a
  # propagation (the provider's at-least-one-of pair).
  dynamic "routing" {
    for_each = each.value.routing != null ? [each.value.routing] : []
    content {
      associated_route_table_id = routing.value.associated_route_table_id != "" ? routing.value.associated_route_table_id : null
      inbound_route_map_id      = routing.value.inbound_route_map_id != "" ? routing.value.inbound_route_map_id : null
      outbound_route_map_id     = routing.value.outbound_route_map_id != "" ? routing.value.outbound_route_map_id : null

      dynamic "propagated_route_table" {
        for_each = routing.value.propagated_route_table != null ? [routing.value.propagated_route_table] : []
        content {
          labels          = propagated_route_table.value.labels
          route_table_ids = propagated_route_table.value.route_table_ids
        }
      }
    }
  }
}
