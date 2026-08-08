# Create the ExpressRoute circuit peering -- the BGP routing
# configuration that makes routes flow through the circuit. The circuit
# must be provider-provisioned (state "Provisioned") before ARM accepts
# this configuration; on a fresh unprovisioned circuit the create fails
# with the provisioning-state error, not a validation error.
resource "azurerm_express_route_circuit_peering" "main" {
  peering_type               = lookup(local.peering_type_wire, var.spec.peering_type, var.spec.peering_type)
  express_route_circuit_name = var.spec.express_route_circuit_name
  resource_group_name        = var.spec.resource_group

  vlan_id = var.spec.vlan_id

  # The /30 pair travels together (spec-validated); private peering
  # accepts a peering created without session addressing.
  primary_peer_address_prefix   = var.spec.primary_peer_address_prefix != "" ? var.spec.primary_peer_address_prefix : null
  secondary_peer_address_prefix = var.spec.secondary_peer_address_prefix != "" ? var.spec.secondary_peer_address_prefix : null

  # Disabling withdraws the IPv4 routes while keeping the configuration
  # (maps to the peering's Enabled/Disabled state in ARM).
  ipv4_enabled = coalesce(var.spec.ipv4_enabled, true)

  # 0 lets Azure record the ASN the router presents; the provider treats
  # peer_asn as optional-computed.
  peer_asn = var.spec.peer_asn > 0 ? var.spec.peer_asn : null

  # The BGP MD5 hash key. ARM never returns it on reads.
  shared_key = var.spec.shared_key != "" ? var.spec.shared_key : null

  # MICROSOFT_PEERING only (spec-validated): without a route filter,
  # Microsoft peering advertises nothing.
  route_filter_id = var.spec.route_filter_id != "" ? var.spec.route_filter_id : null

  # MICROSOFT_PEERING only (spec-validated): the public-prefix
  # advertisement contract Microsoft validates ownership against.
  dynamic "microsoft_peering_config" {
    for_each = var.spec.microsoft_peering_config != null ? [var.spec.microsoft_peering_config] : []
    content {
      advertised_public_prefixes = microsoft_peering_config.value.advertised_public_prefixes
      customer_asn               = microsoft_peering_config.value.customer_asn
      routing_registry_name      = coalesce(microsoft_peering_config.value.routing_registry_name, "NONE")
      advertised_communities     = microsoft_peering_config.value.advertised_communities
    }
  }

  # The IPv6 half: its own /126 pairs and (Microsoft peering) its own
  # advertisement contract. Forbidden on the deprecated public peering
  # (spec-validated).
  dynamic "ipv6" {
    for_each = var.spec.ipv6 != null ? [var.spec.ipv6] : []
    content {
      primary_peer_address_prefix   = ipv6.value.primary_peer_address_prefix
      secondary_peer_address_prefix = ipv6.value.secondary_peer_address_prefix
      enabled                       = coalesce(ipv6.value.enabled, true)
      route_filter_id               = ipv6.value.route_filter_id != "" ? ipv6.value.route_filter_id : null

      dynamic "microsoft_peering" {
        for_each = ipv6.value.microsoft_peering != null ? [ipv6.value.microsoft_peering] : []
        content {
          advertised_public_prefixes = microsoft_peering.value.advertised_public_prefixes
          customer_asn               = microsoft_peering.value.customer_asn
          routing_registry_name      = coalesce(microsoft_peering.value.routing_registry_name, "NONE")
          advertised_communities     = microsoft_peering.value.advertised_communities
        }
      }
    }
  }
}

# The composed Global Reach connections: ARM children of this peering,
# joining it to other circuits' private peerings so the on-premises
# sites behind both circuits reach each other over the Microsoft
# backbone. The near-side peering_id is wired to the peering above --
# only the far side is user surface.
resource "azurerm_express_route_circuit_connection" "connections" {
  for_each = { for connection in var.spec.connections : connection.name => connection }

  name            = each.value.name
  peering_id      = azurerm_express_route_circuit_peering.main.id
  peer_peering_id = each.value.peer_peering_id

  address_prefix_ipv4 = each.value.address_prefix_ipv4
  # ARM-enforced at deploy time: IPv6 is rejected when the parent
  # circuit is ExpressRoute-Direct (port) based.
  address_prefix_ipv6 = each.value.address_prefix_ipv6 != "" ? each.value.address_prefix_ipv6 : null

  # Redeemed when the far circuit belongs to another subscription. ARM
  # masks it on reads, so an imported connection legitimately plans an
  # in-place update on it.
  authorization_key = each.value.authorization_key != "" ? each.value.authorization_key : null
}
