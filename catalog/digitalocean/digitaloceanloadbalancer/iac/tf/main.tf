# DigitalOcean Load Balancer
#
# Provisions a load balancer on DigitalOcean, modeling the complete
# digitalocean_loadbalancer resource surface: regional and global types,
# sizing (slug or scaling units), forwarding rules with TLS termination or
# passthrough, health checks with full threshold tuning, cookie-based sticky
# sessions, backend targeting by Droplet IDs or tag, VPC/subnet placement,
# network visibility and stack, an LB-level firewall, HTTPS redirect, PROXY
# protocol, backend keepalive, idle-timeout tuning, TLS cipher policy,
# project placement, bring-your-own-IP, and the global balancer's domains,
# target balancers, CDN, and regional failover settings.

resource "digitalocean_loadbalancer" "this" {
  name   = var.spec.load_balancer_name
  region = local.region
  type   = local.type

  # Optional VPC / subnet placement (both create-only). Unset VPC means
  # the region's default; GLOBAL balancers take none.
  vpc_uuid    = local.vpc_uuid
  subnet_uuid = local.subnet_uuid

  # Sizing: the provider's request builder is strictly either/or and
  # prefers size_unit. Spec CEL forbids carrying both.
  size      = local.size
  size_unit = local.size_unit

  redirect_http_to_https           = var.spec.redirect_http_to_https
  enable_proxy_protocol            = var.spec.enable_proxy_protocol
  enable_backend_keepalive         = var.spec.enable_backend_keepalive
  disable_lets_encrypt_dns_records = var.spec.disable_lets_encrypt_dns_records
  http_idle_timeout_seconds        = local.http_idle_timeout_seconds
  tls_cipher_policy                = local.tls_cipher_policy
  network                          = local.network
  network_stack                    = local.network_stack
  project_id                       = local.project_id

  # BYOIP: create-only. When unset DigitalOcean allocates the address.
  ip = local.ip

  # Backend targeting. droplet_ids is sent as null when empty so
  # tag-managed membership never diffs against an empty computed set.
  droplet_tag = local.droplet_tag
  droplet_ids = local.droplet_ids_or_null

  target_load_balancer_ids = local.target_load_balancer_ids

  dynamic "forwarding_rule" {
    for_each = coalesce(var.spec.forwarding_rules, [])
    content {
      entry_port      = forwarding_rule.value.entry_port
      entry_protocol  = lower(forwarding_rule.value.entry_protocol)
      target_port     = forwarding_rule.value.target_port
      target_protocol = lower(forwarding_rule.value.target_protocol)
      tls_passthrough = forwarding_rule.value.tls_passthrough
      # Flattened StringValueOrRef: already the certificate NAME.
      certificate_name = try(forwarding_rule.value.certificate_name, "") != "" ? forwarding_rule.value.certificate_name : null
    }
  }

  dynamic "healthcheck" {
    for_each = var.spec.health_check != null ? [var.spec.health_check] : []
    content {
      port     = healthcheck.value.port
      protocol = lower(healthcheck.value.protocol)
      # Path is required for http/https and forbidden for tcp; send only
      # when set so a TCP check never carries an empty path.
      path                     = try(healthcheck.value.path, "") != "" ? healthcheck.value.path : null
      check_interval_seconds   = try(healthcheck.value.check_interval_sec, 0) > 0 ? healthcheck.value.check_interval_sec : null
      response_timeout_seconds = try(healthcheck.value.response_timeout_seconds, 0) > 0 ? healthcheck.value.response_timeout_seconds : null
      unhealthy_threshold      = try(healthcheck.value.unhealthy_threshold, 0) > 0 ? healthcheck.value.unhealthy_threshold : null
      healthy_threshold        = try(healthcheck.value.healthy_threshold, 0) > 0 ? healthcheck.value.healthy_threshold : null
    }
  }

  dynamic "sticky_sessions" {
    for_each = var.spec.sticky_sessions != null ? [var.spec.sticky_sessions] : []
    content {
      type               = sticky_sessions.value.type
      cookie_name        = try(sticky_sessions.value.cookie_name, "") != "" ? sticky_sessions.value.cookie_name : null
      cookie_ttl_seconds = try(sticky_sessions.value.cookie_ttl_seconds, 0) > 0 ? sticky_sessions.value.cookie_ttl_seconds : null
    }
  }

  dynamic "firewall" {
    for_each = var.spec.firewall != null ? [var.spec.firewall] : []
    content {
      allow = length(coalesce(firewall.value.allow, [])) > 0 ? firewall.value.allow : null
      deny  = length(coalesce(firewall.value.deny, [])) > 0 ? firewall.value.deny : null
    }
  }

  dynamic "domains" {
    for_each = coalesce(var.spec.domains, [])
    content {
      name             = domains.value.name
      is_managed       = domains.value.is_managed
      certificate_name = try(domains.value.certificate_name, "") != "" ? domains.value.certificate_name : null
    }
  }

  dynamic "glb_settings" {
    for_each = var.spec.glb_settings != null ? [var.spec.glb_settings] : []
    content {
      target_protocol    = glb_settings.value.target_protocol
      target_port        = glb_settings.value.target_port
      region_priorities  = length(coalesce(glb_settings.value.region_priorities, {})) > 0 ? glb_settings.value.region_priorities : null
      failover_threshold = try(glb_settings.value.failover_threshold, 0) > 0 ? glb_settings.value.failover_threshold : null

      dynamic "cdn" {
        for_each = glb_settings.value.cdn != null ? [glb_settings.value.cdn] : []
        content {
          is_enabled = cdn.value.is_enabled
        }
      }
    }
  }
}
