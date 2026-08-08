# Create the virtual network gateway -- the managed appliance that
# terminates hybrid connectivity (site-to-site VPN, point-to-site VPN,
# VNet-to-VNet, ExpressRoute) in the virtual network's dedicated
# "GatewaySubnet" (an ARM name contract the referenced AzureSubnet must
# satisfy).
#
# Gateways provision SLOWLY (25-45 minutes) and bill hourly per SKU. The
# ForceNew surface (name, region, type, vpn_type, generation, edge zone,
# private-IP enablement, every ip_configuration) is therefore expensive --
# design changes to avoid replacement.
resource "azurerm_virtual_network_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Type, routing model, and SKU are always sent explicitly
  # (Vpn/RouteBased when unspecified) -- deterministic payloads on both
  # engines. The SKU cross-checks (type/generation pairing) are
  # spec-validated.
  type     = local.gateway_type
  vpn_type = local.vpn_type
  sku      = local.sku

  # Omission lets Azure pick the SKU's default generation (the provider
  # treats the field as Computed).
  generation = local.generation

  # Each configuration binds the GatewaySubnet and, on VPN gateways, a
  # public IP (spec-validated pairing; ExpressRoute gateways must not
  # carry one). FIXED AT CREATION. The provider defaults an empty name to
  # "vnetGatewayConfig" (the portal's name) -- mirrored here so both
  # engines produce an identical payload.
  dynamic "ip_configuration" {
    for_each = var.spec.ip_configurations
    content {
      name = (
        ip_configuration.value.name != null && ip_configuration.value.name != ""
      ) ? ip_configuration.value.name : "vnetGatewayConfig"
      subnet_id = ip_configuration.value.subnet_id
      public_ip_address_id = (
        ip_configuration.value.public_ip_address_id != null && ip_configuration.value.public_ip_address_id != ""
      ) ? ip_configuration.value.public_ip_address_id : null
      private_ip_address_allocation = lookup(
        local.allocation_wire,
        coalesce(ip_configuration.value.private_ip_address_allocation, "DYNAMIC"),
        "Dynamic"
      )
    }
  }

  active_active              = var.spec.active_active
  private_ip_address_enabled = var.spec.private_ip_address_enabled
  edge_zone                  = var.spec.edge_zone != "" ? var.spec.edge_zone : null

  bgp_enabled = var.spec.bgp_enabled

  # The gateway's BGP speaker: ASN, route weight, and per-configuration
  # APIPA peering addresses (link-local endpoints some peers -- e.g. AWS
  # site-to-site VPN -- require).
  dynamic "bgp_settings" {
    for_each = var.spec.bgp_settings != null ? [var.spec.bgp_settings] : []
    content {
      asn         = bgp_settings.value.asn
      peer_weight = bgp_settings.value.peer_weight

      dynamic "peering_addresses" {
        for_each = bgp_settings.value.peering_addresses
        content {
          ip_configuration_name = (
            peering_addresses.value.ip_configuration_name != null && peering_addresses.value.ip_configuration_name != ""
          ) ? peering_addresses.value.ip_configuration_name : null
          apipa_addresses = peering_addresses.value.apipa_addresses
        }
      }
    }
  }

  # Custom routes: prefixes the gateway advertises to every connected
  # client and tunnel beyond the VNet's own space.
  dynamic "custom_route" {
    for_each = length(var.spec.custom_route_address_prefixes) > 0 ? [1] : []
    content {
      address_prefixes = var.spec.custom_route_address_prefixes
    }
  }

  # Forced tunneling: the default-route site. References resolve to the
  # local network gateway's ARM id before the module runs.
  default_local_network_gateway_id = (
    var.spec.default_local_network_gateway_id != null && var.spec.default_local_network_gateway_id != ""
  ) ? var.spec.default_local_network_gateway_id : null

  # Point-to-site: address pool, authentication (Entra ID, certificate,
  # RADIUS), tunnel protocols, and per-group routing (VPN gateways only
  # -- spec-validated).
  dynamic "vpn_client_configuration" {
    for_each = var.spec.vpn_client_configuration != null ? [var.spec.vpn_client_configuration] : []
    content {
      address_space = vpn_client_configuration.value.address_spaces

      aad_tenant   = vpn_client_configuration.value.aad_tenant != "" ? vpn_client_configuration.value.aad_tenant : null
      aad_audience = vpn_client_configuration.value.aad_audience != "" ? vpn_client_configuration.value.aad_audience : null
      aad_issuer   = vpn_client_configuration.value.aad_issuer != "" ? vpn_client_configuration.value.aad_issuer : null

      dynamic "root_certificate" {
        for_each = vpn_client_configuration.value.root_certificates
        content {
          name             = root_certificate.value.name
          public_cert_data = root_certificate.value.public_cert_data
        }
      }

      dynamic "revoked_certificate" {
        for_each = vpn_client_configuration.value.revoked_certificates
        content {
          name       = revoked_certificate.value.name
          thumbprint = revoked_certificate.value.thumbprint
        }
      }

      radius_server_address = (
        vpn_client_configuration.value.radius_server_address != ""
      ) ? vpn_client_configuration.value.radius_server_address : null
      radius_server_secret = (
        vpn_client_configuration.value.radius_server_secret != ""
      ) ? vpn_client_configuration.value.radius_server_secret : null

      dynamic "radius_server" {
        for_each = vpn_client_configuration.value.radius_servers
        content {
          address = radius_server.value.address
          secret  = radius_server.value.secret
          score   = radius_server.value.score
        }
      }

      dynamic "ipsec_policy" {
        for_each = vpn_client_configuration.value.ipsec_policy != null ? [vpn_client_configuration.value.ipsec_policy] : []
        content {
          dh_group                  = ipsec_policy.value.dh_group
          ike_encryption            = ipsec_policy.value.ike_encryption
          ike_integrity             = ipsec_policy.value.ike_integrity
          ipsec_encryption          = ipsec_policy.value.ipsec_encryption
          ipsec_integrity           = ipsec_policy.value.ipsec_integrity
          pfs_group                 = ipsec_policy.value.pfs_group
          sa_lifetime_in_seconds    = ipsec_policy.value.sa_lifetime_seconds
          sa_data_size_in_kilobytes = ipsec_policy.value.sa_data_size_kilobytes
        }
      }

      vpn_client_protocols = (
        length(vpn_client_configuration.value.vpn_client_protocols) > 0
      ) ? vpn_client_configuration.value.vpn_client_protocols : null
      vpn_auth_types = (
        length(vpn_client_configuration.value.vpn_auth_types) > 0
      ) ? vpn_client_configuration.value.vpn_auth_types : null

      dynamic "virtual_network_gateway_client_connection" {
        for_each = vpn_client_configuration.value.client_connections
        content {
          name               = virtual_network_gateway_client_connection.value.name
          policy_group_names = virtual_network_gateway_client_connection.value.policy_group_names
          address_prefixes   = virtual_network_gateway_client_connection.value.address_prefixes
        }
      }
    }
  }

  # Policy groups for point-to-site segmentation: members match by Entra
  # ID group, certificate CN, or RADIUS attribute.
  dynamic "policy_group" {
    for_each = var.spec.policy_groups
    content {
      name       = policy_group.value.name
      is_default = policy_group.value.is_default
      priority   = policy_group.value.priority

      dynamic "policy_member" {
        for_each = policy_group.value.policy_members
        content {
          name  = policy_member.value.name
          type  = policy_member.value.type
          value = policy_member.value.value
        }
      }
    }
  }

  bgp_route_translation_for_nat_enabled = var.spec.bgp_route_translation_for_nat_enabled

  # Sent only when enabled -- ARM rejects the parameter on SKUs/types
  # without DNS forwarding support, so omission is the safe default.
  dns_forwarding_enabled = var.spec.dns_forwarding_enabled ? true : null

  # Omitted when null: the provider's default (true) matches the spec's
  # documented default, so both engines send nothing unless the user
  # takes an explicit position.
  ip_sec_replay_protection_enabled = var.spec.ip_sec_replay_protection_enabled

  # ER_GW_SCALE autoscale bounds (spec-validated pairing with the SKU).
  # PARITY-EXCEPTION: the Pulumi engine's classic SDK cannot express
  # these -- ER_GW_SCALE gateways deploy via this engine only.
  minimum_scale_unit = var.spec.minimum_scale_unit
  maximum_scale_unit = var.spec.maximum_scale_unit

  remote_vnet_traffic_enabled = var.spec.remote_vnet_traffic_enabled
  virtual_wan_traffic_enabled = var.spec.virtual_wan_traffic_enabled

  tags = local.final_tags
}

# The composed NAT rules: standalone ARM children of the gateway.
# Connections opt into specific rules via their egress/ingress NAT rule
# id lists; each rule's ARM id surfaces in the nat_rule_ids output under
# its name.
resource "azurerm_virtual_network_gateway_nat_rule" "rules" {
  for_each = { for rule in var.spec.nat_rules : rule.name => rule }

  name                       = each.value.name
  resource_group_name        = var.spec.resource_group
  virtual_network_gateway_id = azurerm_virtual_network_gateway.main.id

  mode = lookup(local.nat_mode_wire, coalesce(each.value.mode, "EGRESS_SNAT"), "EgressSnat")
  type = lookup(local.nat_type_wire, coalesce(each.value.type, "STATIC_NAT"), "Static")

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

  ip_configuration_id = each.value.ip_configuration_id != "" ? each.value.ip_configuration_id : null
}
