# A targeted WARP device profile: the default profile's settings body,
# applied only to the devices matched by the wirefilter expression. A real
# object -- create, update, and delete all do what they say (deleting the
# profile returns its devices to the default profile).
resource "cloudflare_zero_trust_device_custom_profile" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name
  match      = var.spec.match
  precedence = var.spec.precedence

  enabled     = try(var.spec.enabled, null)
  description = try(var.spec.description, "") != "" ? var.spec.description : null

  allow_mode_switch              = try(var.spec.allow_mode_switch, null)
  allow_updates                  = try(var.spec.allow_updates, null)
  allowed_to_leave               = try(var.spec.allowed_to_leave, null)
  auto_connect                   = try(var.spec.auto_connect, null)
  captive_portal                 = try(var.spec.captive_portal, null)
  disable_auto_fallback          = try(var.spec.disable_auto_fallback, null)
  exclude_office_ips             = try(var.spec.exclude_office_ips, null)
  register_interface_ip_with_dns = try(var.spec.register_interface_ip_with_dns, null)
  sccm_vpn_boundary_support      = try(var.spec.sccm_vpn_boundary_support, null)
  switch_locked                  = try(var.spec.switch_locked, null)
  support_url                    = try(var.spec.support_url, "") != "" ? var.spec.support_url : null
  tunnel_protocol                = try(var.spec.tunnel_protocol, "") != "" ? var.spec.tunnel_protocol : null
  lan_allow_minutes              = try(var.spec.lan_allow_minutes, null)
  lan_allow_subnet_size          = try(var.spec.lan_allow_subnet_size, null)

  # Split tunnel: exclude and include are mutually exclusive (spec validation
  # enforces it). Each entry carries exactly one of address or host; the
  # unset one is never sent.
  exclude = length(try(var.spec.exclude, [])) > 0 ? [
    for entry in var.spec.exclude : {
      address     = entry.address != "" ? entry.address : null
      host        = entry.host != "" ? entry.host : null
      description = try(entry.description, "") != "" ? entry.description : null
    }
  ] : null

  include = length(try(var.spec.include, [])) > 0 ? [
    for entry in var.spec.include : {
      address     = entry.address != "" ? entry.address : null
      host        = entry.host != "" ? entry.host : null
      description = try(entry.description, "") != "" ? entry.description : null
    }
  ] : null

  service_mode_v2 = try(var.spec.service_mode_v2, null) != null ? {
    mode = var.spec.service_mode_v2.mode
    port = try(var.spec.service_mode_v2.port, null)
  } : null

  # The provider's attribute for the default virtual network is named
  # "default"; the spec names it default_virtual_network_id for clarity.
  virtual_networks = try(var.spec.virtual_networks, null) != null ? {
    allowed = var.spec.virtual_networks.allowed
    default = var.spec.virtual_networks.default_virtual_network_id
  } : null

  dns_search_suffixes = length(try(var.spec.dns_search_suffixes, [])) > 0 ? [
    for row in var.spec.dns_search_suffixes : {
      suffix      = row.suffix
      description = try(row.description, "") != "" ? row.description : null
    }
  ] : null
}

# This profile's local-DNS fallback list. The profile resource reports
# fallback_domains READ-ONLY; this dedicated per-profile companion is the
# only write path, and it replaces the whole list on every apply
# (declarative: what the spec lists is exactly what exists). The rows ride
# the profile -- deleting the profile retires them with it, which is why the
# companion's own no-op destroy is harmless here. Deployed only when the
# spec declares rows.
resource "cloudflare_zero_trust_device_custom_profile_local_domain_fallback" "main" {
  count = length(try(var.spec.fallback_domains, [])) > 0 ? 1 : 0

  account_id = var.spec.account_id
  policy_id  = cloudflare_zero_trust_device_custom_profile.main.id

  domains = [
    for row in var.spec.fallback_domains : {
      suffix      = row.suffix
      description = try(row.description, "") != "" ? row.description : null
      dns_server  = length(try(row.dns_server, [])) > 0 ? row.dns_server : null
    }
  ]
}
