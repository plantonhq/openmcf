# Create the local network gateway -- Azure's DESCRIPTION of the
# on-premises side of a site-to-site VPN: the device's public endpoint
# and the address space behind it. Deploys nothing on-premises and costs
# nothing to keep; a gateway connection points at it, one per site.
resource "azurerm_local_network_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # The on-premises endpoint: exactly one of address or FQDN
  # (spec-validated) -- ARM stores whichever form is supplied and, for
  # FQDNs, re-resolves periodically.
  gateway_address = var.spec.gateway_address != "" ? var.spec.gateway_address : null
  gateway_fqdn    = var.spec.gateway_fqdn != "" ? var.spec.gateway_fqdn : null

  # Static routing: the prefixes Azure routes into the tunnel (a SET in
  # the provider -- order is not significant). Empty is legal only
  # alongside BGP (spec-validated) -- learned routes carry the site
  # instead.
  address_space = var.spec.address_spaces

  # Dynamic routing: the on-premises BGP speaker. The peering address
  # lives INSIDE the tunnel (the device's tunnel interface), not the
  # device's public endpoint.
  dynamic "bgp_settings" {
    for_each = var.spec.bgp_settings != null ? [var.spec.bgp_settings] : []
    content {
      asn                 = bgp_settings.value.asn
      bgp_peering_address = bgp_settings.value.bgp_peering_address
      peer_weight         = bgp_settings.value.peer_weight
    }
  }

  tags = local.final_tags
}
