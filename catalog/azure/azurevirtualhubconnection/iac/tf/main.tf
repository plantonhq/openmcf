# Create the Virtual Hub Connection -- the attachment that joins one
# spoke VNet to a Virtual WAN hub. The connection is free; what it
# unlocks is the hub's routing. Deleting the hub requires this
# connection to be gone first (the runner's reverse teardown handles
# the ordering).
resource "azurerm_virtual_hub_connection" "main" {
  name                      = var.spec.name
  virtual_hub_id            = var.spec.virtual_hub_id
  remote_virtual_network_id = var.spec.remote_virtual_network_id

  # Off is ARM's default: the spoke keeps its own internet egress. On,
  # the hub advertises 0.0.0.0/0 to this connection (typically paired
  # with a hub firewall via routing intent).
  internet_security_enabled = var.spec.internet_security_enabled

  # Unset routing applies ARM's default behavior: associate with and
  # propagate to the hub's built-in default route table (any-to-any).
  # The spec guarantees a configured block carries at least one of the
  # provider's at-least-one-of trio (association, propagation, static
  # routes).
  dynamic "routing" {
    for_each = var.spec.routing != null ? [var.spec.routing] : []
    content {
      associated_route_table_id = routing.value.associated_route_table_id != "" ? routing.value.associated_route_table_id : null
      inbound_route_map_id      = routing.value.inbound_route_map_id != "" ? routing.value.inbound_route_map_id : null
      outbound_route_map_id     = routing.value.outbound_route_map_id != "" ? routing.value.outbound_route_map_id : null

      # ARM fixes the criteria once the connection is created
      # (ForceNew); Contains is ARM's default.
      static_vnet_local_route_override_criteria = (
        routing.value.static_vnet_local_route_override_criteria == null
        ? "Contains"
        : lookup(local.override_criteria_wire, routing.value.static_vnet_local_route_override_criteria, routing.value.static_vnet_local_route_override_criteria)
      )

      # ARM defaults propagation of static routes ON; unset mirrors it,
      # matching the Pulumi module's nil handling.
      static_vnet_propagate_static_routes_enabled = (
        routing.value.static_vnet_propagate_static_routes_enabled == null
        ? true
        : routing.value.static_vnet_propagate_static_routes_enabled
      )

      dynamic "propagated_route_table" {
        for_each = routing.value.propagated_route_table != null ? [routing.value.propagated_route_table] : []
        content {
          labels          = propagated_route_table.value.labels
          route_table_ids = propagated_route_table.value.route_table_ids
        }
      }

      dynamic "static_vnet_route" {
        for_each = routing.value.static_vnet_routes
        content {
          name                = static_vnet_route.value.name
          address_prefixes    = static_vnet_route.value.address_prefixes
          next_hop_ip_address = static_vnet_route.value.next_hop_ip_address
        }
      }
    }
  }
}
