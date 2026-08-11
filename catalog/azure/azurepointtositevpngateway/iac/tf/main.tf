# Create the point-to-site VPN gateway -- the managed receiver inside
# a virtual hub that individual devices dial into (ARM allows one per
# hub, a slot separate from the hub's site-to-site VPN gateway). The
# gateway bills from creation and is a SLOW resource: creates run
# 30-45 minutes -- the provider's own timeout class is 90 minutes. HOW
# users authenticate lives on the referenced VPN server configuration;
# WHAT addresses clients get comes from the connection configurations
# here.
resource "azurerm_point_to_site_vpn_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Both ForceNew: the gateway is born in its hub, pointing at its
  # authentication policy.
  virtual_hub_id              = var.spec.virtual_hub_id
  vpn_server_configuration_id = var.spec.vpn_server_configuration_id

  # 500 concurrent connections per unit across the managed instance
  # pair. The provider REQUIRES an explicit value (no provider
  # default); the spec's unset applies 1 -- rendered explicitly so the
  # plan shows the real value, mirroring the Pulumi module's nil
  # handling.
  scale_unit = var.spec.scale_unit == null ? 1 : var.spec.scale_unit

  # Off is ARM's default (client internet traffic rides Microsoft's
  # backbone). ForceNew: changing it replaces the gateway.
  routing_preference_internet_enabled = var.spec.routing_preference_internet_enabled

  # Pushed to connecting clients. Emitted only when configured; NOTE:
  # the provider cannot CLEAR a previously set list (its update path
  # skips empty lists) -- removing servers requires replacing the
  # gateway, which the spec documents on the field.
  dns_servers = length(var.spec.dns_servers) > 0 ? var.spec.dns_servers : null

  # One block per named client address pool. Most gateways carry
  # exactly one; multiple pools require OpenVPN on the server
  # configuration and are matched to users via its policy groups.
  dynamic "connection_configuration" {
    for_each = var.spec.connection_configurations
    content {
      name = connection_configuration.value.name

      # The provider nests the pool under vpn_client_address_pool; the
      # spec flattens it to address_prefixes (recorded in the parity
      # manifest).
      vpn_client_address_pool {
        address_prefixes = connection_configuration.value.address_prefixes
      }

      # Off is ARM's default: clients keep their local internet
      # egress. On, the hub advertises 0.0.0.0/0 into the tunnel.
      internet_security_enabled = connection_configuration.value.internet_security_enabled

      # Unset routing applies ARM's default behavior: associate with
      # and propagate to the hub's built-in default route table. A
      # configured block carries its association (the spec requires it
      # -- the provider's own contract). The provider names the
      # propagation targets `ids`; the spec says what they ARE
      # (route_table_ids -- recorded in the parity manifest).
      dynamic "route" {
        for_each = connection_configuration.value.route != null ? [connection_configuration.value.route] : []
        content {
          associated_route_table_id = route.value.associated_route_table_id
          inbound_route_map_id      = route.value.inbound_route_map_id != "" ? route.value.inbound_route_map_id : null
          outbound_route_map_id     = route.value.outbound_route_map_id != "" ? route.value.outbound_route_map_id : null

          dynamic "propagated_route_table" {
            for_each = route.value.propagated_route_table != null ? [route.value.propagated_route_table] : []
            content {
              ids    = propagated_route_table.value.route_table_ids
              labels = propagated_route_table.value.labels
            }
          }
        }
      }
    }
  }

  tags = local.final_tags
}
