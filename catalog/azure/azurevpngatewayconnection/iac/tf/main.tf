# Create the VPN gateway connection -- the tunnel bundle joining one
# branch (a VPN site) to a hub's VPN gateway. ARM addresses it as a
# child of the gateway. The object is free and provisions in minutes;
# each tunnel reaches Connected only when the branch device
# negotiates (provisioned-is-not-connected -- a Succeeded deployment
# with a tunnel in Connecting means the far side disagrees, not that
# the deployment failed).
resource "azurerm_vpn_gateway_connection" "main" {
  name               = var.spec.name
  vpn_gateway_id     = var.spec.vpn_gateway_id
  remote_vpn_site_id = var.spec.remote_vpn_site_id

  # Off is ARM's default: the branch keeps its own internet egress. On,
  # the hub advertises 0.0.0.0/0 to this branch (typically paired with
  # a hub firewall via routing intent).
  internet_security_enabled = var.spec.internet_security_enabled

  # Unset routing applies ARM's default behavior: associate with and
  # propagate to the hub's built-in default route table. A configured
  # block carries its association (the spec requires it -- the
  # provider's own contract).
  dynamic "routing" {
    for_each = var.spec.routing != null ? [var.spec.routing] : []
    content {
      associated_route_table = routing.value.associated_route_table_id
      inbound_route_map_id   = routing.value.inbound_route_map_id != "" ? routing.value.inbound_route_map_id : null
      outbound_route_map_id  = routing.value.outbound_route_map_id != "" ? routing.value.outbound_route_map_id : null

      dynamic "propagated_route_table" {
        for_each = routing.value.propagated_route_table != null ? [routing.value.propagated_route_table] : []
        content {
          route_table_ids = propagated_route_table.value.route_table_ids
          labels          = propagated_route_table.value.labels
        }
      }
    }
  }

  # One tunnel per site link being connected. vpn_site_link_id and
  # bgp_enabled are ForceNew on each tunnel; everything else updates in
  # place.
  dynamic "vpn_link" {
    for_each = var.spec.vpn_links
    content {
      name             = vpn_link.value.name
      vpn_site_link_id = vpn_link.value.vpn_site_link_id

      # ARM's defaults rendered explicitly so the plan shows the real
      # values -- mirroring the Pulumi module's nil handling.
      bandwidth_mbps  = vpn_link.value.bandwidth_mbps == null ? 10 : vpn_link.value.bandwidth_mbps
      protocol        = lookup(local.protocol_wire, coalesce(vpn_link.value.protocol, "IKE_V2"), "IKEv2")
      connection_mode = lookup(local.connection_mode_wire, coalesce(vpn_link.value.connection_mode, "DEFAULT"), "Default")
      route_weight    = vpn_link.value.route_weight

      # Omitted (not 0) when unset -- the provider sends DPD only when
      # configured; ARM's default is 45 seconds.
      dpd_timeout_seconds = vpn_link.value.dpd_timeout_seconds

      # Omit to let Azure generate a key. Sensitive: the value never
      # appears in plan output.
      shared_key = vpn_link.value.shared_key != "" ? vpn_link.value.shared_key : null

      bgp_enabled                           = vpn_link.value.bgp_enabled
      ratelimit_enabled                     = vpn_link.value.ratelimit_enabled
      local_azure_ip_address_enabled        = vpn_link.value.local_azure_ip_address_enabled
      policy_based_traffic_selector_enabled = vpn_link.value.policy_based_traffic_selector_enabled

      egress_nat_rule_ids  = length(vpn_link.value.egress_nat_rule_ids) > 0 ? vpn_link.value.egress_nat_rule_ids : null
      ingress_nat_rule_ids = length(vpn_link.value.ingress_nat_rule_ids) > 0 ? vpn_link.value.ingress_nat_rule_ids : null

      # The spec requires every field of a configured proposal (the
      # provider's contract -- no partial pinning).
      dynamic "ipsec_policy" {
        for_each = vpn_link.value.ipsec_policies
        content {
          sa_lifetime_sec          = ipsec_policy.value.sa_lifetime_sec
          sa_data_size_kb          = ipsec_policy.value.sa_data_size_kb
          encryption_algorithm     = ipsec_policy.value.encryption_algorithm
          integrity_algorithm      = ipsec_policy.value.integrity_algorithm
          ike_encryption_algorithm = ipsec_policy.value.ike_encryption_algorithm
          ike_integrity_algorithm  = ipsec_policy.value.ike_integrity_algorithm
          dh_group                 = ipsec_policy.value.dh_group
          pfs_group                = ipsec_policy.value.pfs_group
        }
      }

      # Which of the gateway's custom APIPA addresses each instance
      # peers from on this tunnel ("Instance0"/"Instance1" are ARM's
      # identifiers -- the spec validates the vocabulary).
      dynamic "custom_bgp_address" {
        for_each = vpn_link.value.custom_bgp_addresses
        content {
          ip_address          = custom_bgp_address.value.ip_address
          ip_configuration_id = custom_bgp_address.value.ip_configuration_id
        }
      }
    }
  }

  # Most connections leave this empty (routing comes from the site's
  # prefixes or BGP).
  dynamic "traffic_selector_policy" {
    for_each = var.spec.traffic_selector_policies
    content {
      local_address_ranges  = traffic_selector_policy.value.local_address_cidrs
      remote_address_ranges = traffic_selector_policy.value.remote_address_cidrs
    }
  }
}
