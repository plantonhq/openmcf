# Gateway DNS location: a named entry point (office, site, network) whose
# DNS traffic Gateway filters. A plain CRUD resource (real
# create/update/delete; only the account forces replacement).
#
# Two provider behaviors the mapping honors:
#   - Update is a full replace at the API: max_ttl is sent whenever the spec
#     declares it, and its omission genuinely resets the behavior to inherit
#     (documented on the spec field).
#   - dns_destination_ips_id is only sent when set -- unset lets Cloudflare
#     auto-assign the shared IPv4 destination pair.
resource "cloudflare_zero_trust_dns_location" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name

  client_default         = try(var.spec.client_default, null)
  ecs_support            = try(var.spec.ecs_support, null)
  dns_destination_ips_id = try(var.spec.dns_destination_ips_id, "") != "" ? var.spec.dns_destination_ips_id : null

  # The endpoints tree declares ALL FOUR types at once (spec validation
  # mirrors the provider schema's requirement).
  endpoints = try(var.spec.endpoints, null) != null ? {
    doh = {
      enabled       = try(var.spec.endpoints.doh.enabled, null)
      require_token = try(var.spec.endpoints.doh.require_token, null)
      networks = length(try(var.spec.endpoints.doh.networks, [])) > 0 ? [
        for row in var.spec.endpoints.doh.networks : { network = row.network }
      ] : null
    }
    dot = {
      enabled = try(var.spec.endpoints.dot.enabled, null)
      networks = length(try(var.spec.endpoints.dot.networks, [])) > 0 ? [
        for row in var.spec.endpoints.dot.networks : { network = row.network }
      ] : null
    }
    ipv4 = {
      enabled = try(var.spec.endpoints.ipv4.enabled, null)
    }
    ipv6 = {
      enabled = try(var.spec.endpoints.ipv6.enabled, null)
      networks = length(try(var.spec.endpoints.ipv6.networks, [])) > 0 ? [
        for row in var.spec.endpoints.ipv6.networks : { network = row.network }
      ] : null
    }
  } : null

  networks = length(try(var.spec.networks, [])) > 0 ? [
    for row in var.spec.networks : { network = row.network }
  ] : null

  max_ttl = try(var.spec.max_ttl, null) != null ? {
    mode     = var.spec.max_ttl.mode
    ttl_secs = try(var.spec.max_ttl.ttl_secs, null)
  } : null
}
