# Gateway DNS location: a named entry point (office, site, network) whose
# DNS traffic Gateway filters. A plain CRUD resource (real
# create/update/delete; only the account forces replacement).
#
# Provider behaviors the mapping honors:
#   - Update is a full replace at the API: max_ttl is sent whenever the spec
#     declares it, and its omission genuinely resets the behavior to inherit
#     (documented on the spec field).
#   - dns_destination_ips_id is only sent when set -- unset lets Cloudflare
#     auto-assign the shared IPv4 destination pair.
#   - max_ttl and every networks list are ALWAYS sent as known values
#     (max_ttl coalesces to the documented server default {mode: inherit};
#     an absent networks list is sent as []). At provider v5.23.0/v5.24.0
#     the Go model types these computed-optional attributes as raw
#     pointers that cannot hold "unknown", so leaving any of them null in
#     config plans it as unknown and CRASHES the apply-time conversion
#     ("Value Conversion Error ... target type cannot handle unknown
#     values"; measured live 2026-08-26, unfixed on provider main). The
#     explicit server defaults keep the planned value known and change
#     nothing semantically.
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
      networks = [
        for row in try(var.spec.endpoints.doh.networks, []) : { network = row.network }
      ]
    }
    dot = {
      enabled = try(var.spec.endpoints.dot.enabled, null)
      networks = [
        for row in try(var.spec.endpoints.dot.networks, []) : { network = row.network }
      ]
    }
    ipv4 = {
      enabled = try(var.spec.endpoints.ipv4.enabled, null)
    }
    ipv6 = {
      enabled = try(var.spec.endpoints.ipv6.enabled, null)
      networks = [
        for row in try(var.spec.endpoints.ipv6.networks, []) : { network = row.network }
      ]
    }
  } : null

  networks = [
    for row in try(var.spec.networks, []) : { network = row.network }
  ]

  max_ttl = try(var.spec.max_ttl, null) != null ? {
    mode     = var.spec.max_ttl.mode
    ttl_secs = try(var.spec.max_ttl.ttl_secs, null)
    } : {
    mode     = "inherit"
    ttl_secs = null
  }
}
