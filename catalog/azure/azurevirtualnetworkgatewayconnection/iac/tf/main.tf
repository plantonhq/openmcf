# Create the gateway connection -- the tunnel object joining a virtual
# network gateway to its far side: an on-premises device (IPsec, via a
# local network gateway), another gateway (Vnet2Vnet), or an ExpressRoute
# circuit. The type-specific far-side requirements are spec-validated.
#
# PROVISIONED IS NOT CONNECTED: ARM provisions the connection object as
# soon as the parameters are valid; the tunnel reaches Connected only
# when the far side negotiates successfully. A Succeeded deployment with
# a tunnel stuck in Connecting means the far side, the shared key, or the
# IPsec parameters disagree -- not a failed deployment.
resource "azurerm_virtual_network_gateway_connection" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  type                       = local.connection_type
  virtual_network_gateway_id = var.spec.virtual_network_gateway_id

  # The far side, by connection type (spec-validated pairings).
  # References resolve to ARM ids before the module runs.
  local_network_gateway_id = (
    var.spec.local_network_gateway_id != null && var.spec.local_network_gateway_id != ""
  ) ? var.spec.local_network_gateway_id : null
  peer_virtual_network_gateway_id = (
    var.spec.peer_virtual_network_gateway_id != null && var.spec.peer_virtual_network_gateway_id != ""
  ) ? var.spec.peer_virtual_network_gateway_id : null
  express_route_circuit_id = (
    var.spec.express_route_circuit_id != null && var.spec.express_route_circuit_id != ""
  ) ? var.spec.express_route_circuit_id : null

  # The pre-shared key: omitted when unset so Azure generates one
  # (readable back from the connection's shared-key API). Secrets are
  # reference-resolved at deploy time -- never stored in manifests.
  shared_key = (
    var.spec.shared_key != null && var.spec.shared_key != ""
  ) ? var.spec.shared_key : null
  authorization_key = (
    var.spec.authorization_key != null && var.spec.authorization_key != ""
  ) ? var.spec.authorization_key : null

  bgp_enabled = var.spec.bgp_enabled

  # Custom APIPA BGP endpoints (IPsec + BGP + active-active gateway --
  # the gateway-side prerequisites are documented on the spec field; ARM
  # enforces them live).
  dynamic "custom_bgp_addresses" {
    for_each = var.spec.custom_bgp_addresses != null ? [var.spec.custom_bgp_addresses] : []
    content {
      primary   = custom_bgp_addresses.value.primary
      secondary = custom_bgp_addresses.value.secondary != "" ? custom_bgp_addresses.value.secondary : null
    }
  }

  dpd_timeout_seconds = var.spec.dpd_timeout_seconds
  connection_protocol = local.connection_protocol
  connection_mode     = local.connection_mode
  routing_weight      = var.spec.routing_weight

  # Gateway NAT rules this connection opts into, by ARM id (the owning
  # gateway publishes them in its nat_rule_ids output).
  egress_nat_rule_ids  = length(var.spec.egress_nat_rule_ids) > 0 ? var.spec.egress_nat_rule_ids : null
  ingress_nat_rule_ids = length(var.spec.ingress_nat_rule_ids) > 0 ? var.spec.ingress_nat_rule_ids : null

  use_policy_based_traffic_selectors = var.spec.use_policy_based_traffic_selectors
  express_route_gateway_bypass       = var.spec.express_route_gateway_bypass
  private_link_fast_path_enabled     = var.spec.private_link_fast_path_enabled
  local_azure_ip_address_enabled     = var.spec.local_azure_ip_address_enabled

  # PARITY-EXCEPTION: this engine renders every selector; the Pulumi
  # engine's classic SDK models exactly one -- manifests with several
  # deploy via this engine only.
  dynamic "traffic_selector_policy" {
    for_each = var.spec.traffic_selector_policies
    content {
      local_address_cidrs  = traffic_selector_policy.value.local_address_cidrs
      remote_address_cidrs = traffic_selector_policy.value.remote_address_cidrs
    }
  }

  # The custom IPsec/IKE proposal: the six algorithm fields are
  # spec-required together; the SA bounds are omitted when unset so Azure
  # fills its defaults.
  dynamic "ipsec_policy" {
    for_each = var.spec.ipsec_policy != null ? [var.spec.ipsec_policy] : []
    content {
      dh_group         = ipsec_policy.value.dh_group
      ike_encryption   = ipsec_policy.value.ike_encryption
      ike_integrity    = ipsec_policy.value.ike_integrity
      ipsec_encryption = ipsec_policy.value.ipsec_encryption
      ipsec_integrity  = ipsec_policy.value.ipsec_integrity
      pfs_group        = ipsec_policy.value.pfs_group
      sa_datasize      = ipsec_policy.value.sa_datasize
      sa_lifetime      = ipsec_policy.value.sa_lifetime
    }
  }

  tags = local.final_tags
}
