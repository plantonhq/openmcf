# The Application Gateway -- one atomic ARM resource carrying every
# sub-object (frontends, ports, listeners, pools, settings, rules, path
# maps, probes, certificates, profiles, redirects, rewrites). Sub-objects
# wire to each other BY NAME within this resource; what other resources
# need to reach is exported as name-keyed map outputs. Applies routinely
# run 15-25 minutes -- Azure's slowest networking resource.
resource "azurerm_application_gateway" "main" {
  name                = var.spec.name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # SKU name and tier carry the same value on the v2 platform (and
  # Basic). Capacity is fixed sizing; autoscale replaces it (spec
  # validation guarantees exactly one).
  sku {
    name     = local.sku_map[var.spec.sku]
    tier     = local.sku_map[var.spec.sku]
    capacity = var.spec.capacity
  }

  dynamic "autoscale_configuration" {
    for_each = var.spec.autoscale != null ? [var.spec.autoscale] : []
    content {
      min_capacity = autoscale_configuration.value.min_capacity
      max_capacity = autoscale_configuration.value.max_capacity
    }
  }

  # The gateway's dedicated-subnet anchor -- pure ARM plumbing derived
  # from the spec's subnet_id; users never name it.
  gateway_ip_configuration {
    name      = local.gateway_ip_configuration_name
    subnet_id = var.spec.subnet_id
  }

  zones         = var.spec.zones
  fips_enabled  = var.spec.fips_enabled
  http2_enabled = var.spec.http2_enabled

  # The WAF policy attachment (WAF_v2 only; per-listener and per-path
  # overrides live on their blocks).
  firewall_policy_id                = var.spec.firewall_policy_id
  force_firewall_policy_association = var.spec.force_firewall_policy_association

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = identity.value.identity_ids
    }
  }

  # Frontends: public (a referenced Standard public IP) or private (an
  # address in the gateway subnet). Each frontend's ARM ID is exported in
  # the frontend_ip_configuration_ids map output.
  dynamic "frontend_ip_configuration" {
    for_each = var.spec.frontend_ip_configurations
    content {
      name                 = frontend_ip_configuration.value.name
      public_ip_address_id = frontend_ip_configuration.value.public_ip_address_id
      subnet_id            = frontend_ip_configuration.value.subnet_id
      private_ip_address   = frontend_ip_configuration.value.private_ip_address
      private_ip_address_allocation = (
        frontend_ip_configuration.value.private_ip_address_allocation == null || frontend_ip_configuration.value.private_ip_address_allocation == "" ? null :
        local.ip_allocation_map[frontend_ip_configuration.value.private_ip_address_allocation]
      )
      private_link_configuration_name = frontend_ip_configuration.value.private_link_configuration_name
    }
  }

  dynamic "frontend_port" {
    for_each = var.spec.frontend_ports
    content {
      name = frontend_port.value.name
      port = frontend_port.value.port
    }
  }

  # Backend pools. Membership can also arrive member-side: NICs and scale
  # sets join through the backend_address_pool_ids map output.
  dynamic "backend_address_pool" {
    for_each = var.spec.backend_address_pools
    content {
      name         = backend_address_pool.value.name
      fqdns        = backend_address_pool.value.fqdns
      ip_addresses = backend_address_pool.value.ip_addresses
    }
  }

  # L7 backend settings. cookie_based_affinity is ARM's Enabled/Disabled
  # string behind the spec's boolean.
  dynamic "backend_http_settings" {
    for_each = var.spec.backend_http_settings
    content {
      name                                 = backend_http_settings.value.name
      port                                 = backend_http_settings.value.port
      protocol                             = local.protocol_map[backend_http_settings.value.protocol]
      cookie_based_affinity                = backend_http_settings.value.cookie_based_affinity_enabled ? "Enabled" : "Disabled"
      affinity_cookie_name                 = backend_http_settings.value.affinity_cookie_name
      path                                 = backend_http_settings.value.path
      request_timeout                      = backend_http_settings.value.request_timeout
      probe_name                           = backend_http_settings.value.probe_name
      host_name                            = backend_http_settings.value.host_name
      pick_host_name_from_backend_address  = backend_http_settings.value.pick_host_name_from_backend_address
      trusted_root_certificate_names       = backend_http_settings.value.trusted_root_certificate_names
      dedicated_backend_connection_enabled = backend_http_settings.value.dedicated_backend_connection_enabled

      dynamic "connection_draining" {
        for_each = backend_http_settings.value.connection_draining != null ? [backend_http_settings.value.connection_draining] : []
        content {
          enabled           = connection_draining.value.enabled
          drain_timeout_sec = connection_draining.value.drain_timeout_sec
        }
      }
    }
  }

  # L7 listeners.
  dynamic "http_listener" {
    for_each = var.spec.http_listeners
    content {
      name                           = http_listener.value.name
      frontend_ip_configuration_name = http_listener.value.frontend_ip_configuration_name
      frontend_port_name             = http_listener.value.frontend_port_name
      protocol                       = local.protocol_map[http_listener.value.protocol]
      host_names                     = http_listener.value.host_names
      ssl_certificate_name           = http_listener.value.ssl_certificate_name
      require_sni                    = http_listener.value.require_sni
      ssl_profile_name               = http_listener.value.ssl_profile_name
      firewall_policy_id             = http_listener.value.firewall_policy_id

      dynamic "custom_error_configuration" {
        for_each = http_listener.value.custom_error_configurations
        content {
          status_code           = local.status_code_map[custom_error_configuration.value.status_code]
          custom_error_page_url = custom_error_configuration.value.custom_error_page_url
        }
      }
    }
  }

  # L7 routing rules.
  dynamic "request_routing_rule" {
    for_each = var.spec.request_routing_rules
    content {
      name                        = request_routing_rule.value.name
      rule_type                   = local.rule_type_map[request_routing_rule.value.rule_type]
      http_listener_name          = request_routing_rule.value.http_listener_name
      priority                    = request_routing_rule.value.priority
      backend_address_pool_name   = request_routing_rule.value.backend_address_pool_name
      backend_http_settings_name  = request_routing_rule.value.backend_http_settings_name
      url_path_map_name           = request_routing_rule.value.url_path_map_name
      redirect_configuration_name = request_routing_rule.value.redirect_configuration_name
      rewrite_rule_set_name       = request_routing_rule.value.rewrite_rule_set_name
    }
  }

  # Path-based routing.
  dynamic "url_path_map" {
    for_each = var.spec.url_path_maps
    content {
      name                                = url_path_map.value.name
      default_backend_address_pool_name   = url_path_map.value.default_backend_address_pool_name
      default_backend_http_settings_name  = url_path_map.value.default_backend_http_settings_name
      default_redirect_configuration_name = url_path_map.value.default_redirect_configuration_name
      default_rewrite_rule_set_name       = url_path_map.value.default_rewrite_rule_set_name

      dynamic "path_rule" {
        for_each = url_path_map.value.path_rules
        content {
          name                        = path_rule.value.name
          paths                       = path_rule.value.paths
          backend_address_pool_name   = path_rule.value.backend_address_pool_name
          backend_http_settings_name  = path_rule.value.backend_http_settings_name
          redirect_configuration_name = path_rule.value.redirect_configuration_name
          rewrite_rule_set_name       = path_rule.value.rewrite_rule_set_name
          firewall_policy_id          = path_rule.value.firewall_policy_id
        }
      }
    }
  }

  # Custom health probes (HTTP/HTTPS GET probes or TCP/TLS connection
  # probes for layer-4 backends).
  dynamic "probe" {
    for_each = var.spec.probes
    content {
      name                                      = probe.value.name
      protocol                                  = local.protocol_map[probe.value.protocol]
      host                                      = probe.value.host
      pick_host_name_from_backend_http_settings = probe.value.pick_host_name_from_backend_http_settings
      path                                      = probe.value.path
      interval                                  = probe.value.interval
      timeout                                   = probe.value.timeout
      unhealthy_threshold                       = probe.value.unhealthy_threshold
      port                                      = probe.value.port
      minimum_servers                           = probe.value.minimum_servers
      proxy_protocol_header_enabled             = probe.value.proxy_protocol_header_enabled

      dynamic "match" {
        for_each = probe.value.match != null ? [probe.value.match] : []
        content {
          status_code = match.value.status_codes
          body        = match.value.body
        }
      }
    }
  }

  # TLS certificates: Key Vault (renewals propagate) or inline PFX.
  dynamic "ssl_certificate" {
    for_each = var.spec.ssl_certificates
    content {
      name                = ssl_certificate.value.name
      key_vault_secret_id = ssl_certificate.value.key_vault_secret_id
      data                = ssl_certificate.value.data
      password            = ssl_certificate.value.password
    }
  }

  dynamic "trusted_root_certificate" {
    for_each = var.spec.trusted_root_certificates
    content {
      name                = trusted_root_certificate.value.name
      key_vault_secret_id = trusted_root_certificate.value.key_vault_secret_id
      data                = trusted_root_certificate.value.data
    }
  }

  dynamic "trusted_client_certificate" {
    for_each = var.spec.trusted_client_certificates
    content {
      name = trusted_client_certificate.value.name
      data = trusted_client_certificate.value.data
    }
  }

  # Named TLS postures (mutual TLS + per-profile policy).
  dynamic "ssl_profile" {
    for_each = var.spec.ssl_profiles
    content {
      name                                = ssl_profile.value.name
      trusted_client_certificate_names    = ssl_profile.value.trusted_client_certificate_names
      verify_client_certificate_issuer_dn = ssl_profile.value.verify_client_certificate_issuer_dn
      verify_client_certificate_revocation = (
        ssl_profile.value.verify_client_certificate_revocation == null || ssl_profile.value.verify_client_certificate_revocation == "" ? null :
        local.revocation_check_map[ssl_profile.value.verify_client_certificate_revocation]
      )

      dynamic "ssl_policy" {
        for_each = ssl_profile.value.ssl_policy != null ? [ssl_profile.value.ssl_policy] : []
        content {
          policy_type = (
            ssl_policy.value.policy_type == null || ssl_policy.value.policy_type == "" ? null :
            local.ssl_policy_type_map[ssl_policy.value.policy_type]
          )
          policy_name = ssl_policy.value.policy_name
          min_protocol_version = (
            ssl_policy.value.min_protocol_version == null || ssl_policy.value.min_protocol_version == "" ? null :
            local.tls_protocol_map[ssl_policy.value.min_protocol_version]
          )
          cipher_suites      = ssl_policy.value.cipher_suites
          disabled_protocols = [for protocol in ssl_policy.value.disabled_protocols : local.tls_protocol_map[protocol]]
        }
      }
    }
  }

  # The gateway-wide TLS policy (listeners with an ssl_profile use the
  # profile's policy instead).
  dynamic "ssl_policy" {
    for_each = var.spec.ssl_policy != null ? [var.spec.ssl_policy] : []
    content {
      policy_type = (
        ssl_policy.value.policy_type == null || ssl_policy.value.policy_type == "" ? null :
        local.ssl_policy_type_map[ssl_policy.value.policy_type]
      )
      policy_name = ssl_policy.value.policy_name
      min_protocol_version = (
        ssl_policy.value.min_protocol_version == null || ssl_policy.value.min_protocol_version == "" ? null :
        local.tls_protocol_map[ssl_policy.value.min_protocol_version]
      )
      cipher_suites      = ssl_policy.value.cipher_suites
      disabled_protocols = [for protocol in ssl_policy.value.disabled_protocols : local.tls_protocol_map[protocol]]
    }
  }

  dynamic "redirect_configuration" {
    for_each = var.spec.redirect_configurations
    content {
      name                 = redirect_configuration.value.name
      redirect_type        = local.redirect_type_map[redirect_configuration.value.redirect_type]
      target_listener_name = redirect_configuration.value.target_listener_name
      target_url           = redirect_configuration.value.target_url
      include_path         = redirect_configuration.value.include_path
      include_query_string = redirect_configuration.value.include_query_string
    }
  }

  # Rewrite rule sets (header edits + URL rewrites).
  dynamic "rewrite_rule_set" {
    for_each = var.spec.rewrite_rule_sets
    content {
      name = rewrite_rule_set.value.name

      dynamic "rewrite_rule" {
        for_each = rewrite_rule_set.value.rewrite_rules
        content {
          name          = rewrite_rule.value.name
          rule_sequence = rewrite_rule.value.rule_sequence

          dynamic "condition" {
            for_each = rewrite_rule.value.conditions
            content {
              variable    = condition.value.variable
              pattern     = condition.value.pattern
              ignore_case = condition.value.ignore_case
              negate      = condition.value.negate
            }
          }

          dynamic "request_header_configuration" {
            for_each = rewrite_rule.value.request_header_configurations
            content {
              header_name  = request_header_configuration.value.header_name
              header_value = request_header_configuration.value.header_value
            }
          }

          dynamic "response_header_configuration" {
            for_each = rewrite_rule.value.response_header_configurations
            content {
              header_name  = response_header_configuration.value.header_name
              header_value = response_header_configuration.value.header_value
            }
          }

          dynamic "url" {
            for_each = rewrite_rule.value.url != null ? [rewrite_rule.value.url] : []
            content {
              path         = url.value.path
              query_string = url.value.query_string
              components = (
                url.value.components == null || url.value.components == "" ? null :
                local.url_component_map[url.value.components]
              )
              reroute = url.value.reroute
            }
          }
        }
      }
    }
  }

  # Layer-4 (TCP/TLS proxy) trio.
  dynamic "listener" {
    for_each = var.spec.listeners
    content {
      name                           = listener.value.name
      frontend_ip_configuration_name = listener.value.frontend_ip_configuration_name
      frontend_port_name             = listener.value.frontend_port_name
      protocol                       = local.protocol_map[listener.value.protocol]
      host_names                     = listener.value.host_names
      ssl_certificate_name           = listener.value.ssl_certificate_name
      ssl_profile_name               = listener.value.ssl_profile_name
    }
  }

  dynamic "backend" {
    for_each = var.spec.backends
    content {
      name                           = backend.value.name
      port                           = backend.value.port
      protocol                       = local.protocol_map[backend.value.protocol]
      client_ip_preservation_enabled = backend.value.client_ip_preservation_enabled
      host_name                      = backend.value.host_name
      probe_name                     = backend.value.probe_name
      timeout_in_seconds             = backend.value.timeout_in_seconds
      trusted_root_certificate_names = backend.value.trusted_root_certificate_names
    }
  }

  dynamic "routing_rule" {
    for_each = var.spec.routing_rules
    content {
      name                      = routing_rule.value.name
      listener_name             = routing_rule.value.listener_name
      backend_address_pool_name = routing_rule.value.backend_address_pool_name
      backend_name              = routing_rule.value.backend_settings_name
      priority                  = routing_rule.value.priority
    }
  }

  # Gateway-wide custom error pages.
  dynamic "custom_error_configuration" {
    for_each = var.spec.custom_error_configurations
    content {
      status_code           = local.status_code_map[custom_error_configuration.value.status_code]
      custom_error_page_url = custom_error_configuration.value.custom_error_page_url
    }
  }

  # Private Link exposure of frontends.
  dynamic "private_link_configuration" {
    for_each = var.spec.private_link_configurations
    content {
      name = private_link_configuration.value.name

      dynamic "ip_configuration" {
        for_each = private_link_configuration.value.ip_configurations
        content {
          name                          = ip_configuration.value.name
          subnet_id                     = ip_configuration.value.subnet_id
          private_ip_address            = ip_configuration.value.private_ip_address
          private_ip_address_allocation = local.ip_allocation_map[ip_configuration.value.private_ip_address_allocation]
          primary                       = ip_configuration.value.primary
        }
      }
    }
  }

  # Request/response buffering (Azure defaults both to true when the
  # block is absent).
  dynamic "global" {
    for_each = var.spec.global_configuration != null ? [var.spec.global_configuration] : []
    content {
      request_buffering_enabled  = global.value.request_buffering_enabled
      response_buffering_enabled = global.value.response_buffering_enabled
    }
  }

  tags = local.final_tags
}
