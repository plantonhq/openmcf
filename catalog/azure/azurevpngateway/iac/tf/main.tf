# Create the Virtual WAN VPN gateway -- the managed site-to-site VPN
# terminator inside a virtual hub (ARM allows one per hub). The
# gateway bills from creation (~$0.36/hr per scale unit class) and is
# a SLOW resource: creates run 30-45 minutes, deletes 10-20 -- the
# provider's own timeout class is 90 minutes. Deleting it requires its
# connections to be gone first (the runner's reverse teardown handles
# the ordering).
resource "azurerm_vpn_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  virtual_hub_id      = var.spec.virtual_hub_id

  # ARM's default ("Microsoft Network") rendered explicitly so the plan
  # shows the real value -- mirroring the Pulumi module's nil handling.
  # ForceNew: changing it replaces the gateway.
  routing_preference = local.routing_preference

  # 500 Mbps per unit across the managed active-active pair. The
  # provider's default of 1 rendered explicitly.
  scale_unit = var.spec.scale_unit == null ? 1 : var.spec.scale_unit

  # Only meaningful when NAT rules are configured on BGP-enabled
  # tunnels; off is ARM's default.
  bgp_route_translation_for_nat_enabled = var.spec.bgp_route_translation_for_nat_enabled

  # asn/peer_weight go in the create; the custom APIPA addresses are
  # the part ARM only accepts AFTER the gateway exists -- the provider
  # applies them in a second call, which is why they update in place
  # while asn/peer_weight are ForceNew.
  dynamic "bgp_settings" {
    for_each = var.spec.bgp_settings != null ? [var.spec.bgp_settings] : []
    content {
      asn         = bgp_settings.value.asn
      peer_weight = bgp_settings.value.peer_weight

      dynamic "instance_0_bgp_peering_address" {
        for_each = bgp_settings.value.instance_0_bgp_peering_address != null ? [bgp_settings.value.instance_0_bgp_peering_address] : []
        content {
          custom_ips = instance_0_bgp_peering_address.value.custom_ips
        }
      }

      dynamic "instance_1_bgp_peering_address" {
        for_each = bgp_settings.value.instance_1_bgp_peering_address != null ? [bgp_settings.value.instance_1_bgp_peering_address] : []
        content {
          custom_ips = instance_1_bgp_peering_address.value.custom_ips
        }
      }
    }
  }

  tags = local.final_tags
}

# The composed NAT rules: standalone ARM children of the gateway, one
# per spec entry, keyed by name (how overlapping branch address spaces
# are translated). Tunnels opt in via a connection link's
# egress/ingress NAT rule id lists -- the nat_rule_ids output publishes
# each rule's ARM id for exactly that.
resource "azurerm_vpn_gateway_nat_rule" "nat_rules" {
  for_each = { for nat_rule in var.spec.nat_rules : nat_rule.name => nat_rule }

  name           = each.value.name
  vpn_gateway_id = azurerm_vpn_gateway.main.id

  # ARM's defaults ("EgressSnat"/"Static") rendered explicitly so the
  # plan shows the real values. Both are ForceNew on the rule.
  mode = lookup(local.nat_rule_mode_wire, coalesce(each.value.mode, "EGRESS_SNAT"), "EgressSnat")
  type = lookup(local.nat_rule_type_wire, coalesce(each.value.type, "STATIC_NAT"), "Static")

  dynamic "external_mapping" {
    for_each = each.value.external_mappings
    content {
      address_space = external_mapping.value.address_space
      port_range    = external_mapping.value.port_range != "" ? external_mapping.value.port_range : null
    }
  }

  dynamic "internal_mapping" {
    for_each = each.value.internal_mappings
    content {
      address_space = internal_mapping.value.address_space
      port_range    = internal_mapping.value.port_range != "" ? internal_mapping.value.port_range : null
    }
  }

  # Unspecified (rendered as "" by the tfvars converter) applies the
  # rule on both gateway instances (ARM's default) -- emit null, not a
  # value.
  ip_configuration_id = lookup(local.nat_rule_ip_configuration_wire, each.value.ip_configuration, null)
}
