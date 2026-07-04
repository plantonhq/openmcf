variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Application Gateway specification"
  type = object({
    # The Azure region the gateway lives in.
    region = string

    # The resource group the gateway lives in. References are resolved to
    # a literal name by the platform before the module runs.
    resource_group = string

    # The gateway's name, unique within the resource group.
    name = string

    # The gateway's DEDICATED subnet, as a resolved ARM ID.
    subnet_id = string

    # The SKU, as the spec enum's name string (BASIC, STANDARD_V2,
    # WAF_V2). Name and tier move together on the v2 platform.
    sku = string

    # Fixed instance count XOR autoscale (spec validation enforces
    # exactly one).
    capacity = optional(number)
    autoscale = optional(object({
      min_capacity = number
      max_capacity = optional(number)
    }))

    # Availability zones the gateway spans.
    zones = optional(list(string), [])

    # The gateway's managed identity; type arrives as the spec enum's
    # name string, identity_ids as resolved ARM IDs.
    identity = optional(object({
      type         = string
      identity_ids = optional(list(string), [])
    }))

    # Frontend IP configurations: public (a resolved public-IP ARM ID) or
    # private (a subnet + allocation).
    frontend_ip_configurations = list(object({
      name                            = string
      public_ip_address_id            = optional(string)
      subnet_id                       = optional(string)
      private_ip_address              = optional(string)
      private_ip_address_allocation   = optional(string)
      private_link_configuration_name = optional(string)
    }))

    # Frontend ports, declared once and referenced by name.
    frontend_ports = list(object({
      name = string
      port = number
    }))

    # Backend pools (FQDNs and/or IPs; or joined member-side through the
    # backend_address_pool_ids map output).
    backend_address_pools = list(object({
      name         = string
      fqdns        = optional(list(string), [])
      ip_addresses = optional(list(string), [])
    }))

    # L7 backend settings; protocol arrives as the spec enum's name
    # string (HTTP/HTTPS).
    backend_http_settings = optional(list(object({
      name                                = string
      port                                = number
      protocol                            = string
      cookie_based_affinity_enabled       = optional(bool, false)
      affinity_cookie_name                = optional(string)
      path                                = optional(string)
      request_timeout                     = optional(number, 30)
      probe_name                          = optional(string)
      host_name                           = optional(string)
      pick_host_name_from_backend_address = optional(bool, false)
      trusted_root_certificate_names      = optional(list(string), [])
      connection_draining = optional(object({
        enabled           = optional(bool, false)
        drain_timeout_sec = number
      }))
      dedicated_backend_connection_enabled = optional(bool, false)
    })), [])

    # L7 listeners; protocol arrives as HTTP/HTTPS.
    http_listeners = optional(list(object({
      name                           = string
      frontend_ip_configuration_name = string
      frontend_port_name             = string
      protocol                       = string
      host_names                     = optional(list(string), [])
      ssl_certificate_name           = optional(string)
      require_sni                    = optional(bool, false)
      ssl_profile_name               = optional(string)
      firewall_policy_id             = optional(string)
      custom_error_configurations = optional(list(object({
        status_code           = string
        custom_error_page_url = string
      })), [])
    })), [])

    # L7 routing rules; rule_type arrives as BASIC_ROUTING /
    # PATH_BASED_ROUTING.
    request_routing_rules = optional(list(object({
      name                        = string
      rule_type                   = string
      http_listener_name          = string
      priority                    = number
      backend_address_pool_name   = optional(string)
      backend_http_settings_name  = optional(string)
      url_path_map_name           = optional(string)
      redirect_configuration_name = optional(string)
      rewrite_rule_set_name       = optional(string)
    })), [])

    # URL path maps for path-based routing.
    url_path_maps = optional(list(object({
      name                                = string
      default_backend_address_pool_name   = optional(string)
      default_backend_http_settings_name  = optional(string)
      default_redirect_configuration_name = optional(string)
      default_rewrite_rule_set_name       = optional(string)
      path_rules = list(object({
        name                        = string
        paths                       = list(string)
        backend_address_pool_name   = optional(string)
        backend_http_settings_name  = optional(string)
        redirect_configuration_name = optional(string)
        rewrite_rule_set_name       = optional(string)
        firewall_policy_id          = optional(string)
      }))
    })), [])

    # Custom health probes; protocol arrives as HTTP/HTTPS/TCP/TLS.
    probes = optional(list(object({
      name                                      = string
      protocol                                  = string
      host                                      = optional(string)
      pick_host_name_from_backend_http_settings = optional(bool, false)
      path                                      = optional(string)
      interval                                  = number
      timeout                                   = number
      unhealthy_threshold                       = number
      port                                      = optional(number)
      minimum_servers                           = optional(number)
      proxy_protocol_header_enabled             = optional(bool, false)
      match = optional(object({
        status_codes = list(string)
        body         = optional(string)
      }))
    })), [])

    # TLS certificates: Key Vault secret ID XOR inline PFX data.
    ssl_certificates = optional(list(object({
      name                = string
      key_vault_secret_id = optional(string)
      data                = optional(string)
      password            = optional(string)
    })), [])

    # Backend-trust CA bundles: Key Vault secret ID XOR inline data.
    trusted_root_certificates = optional(list(object({
      name                = string
      key_vault_secret_id = optional(string)
      data                = optional(string)
    })), [])

    # Client-CA bundles for mutual TLS.
    trusted_client_certificates = optional(list(object({
      name = string
      data = string
    })), [])

    # Named TLS postures; the nested ssl_policy carries enum name strings
    # (policy_type PREDEFINED/CUSTOM/CUSTOM_V2, protocols TLS_V1_0-TLS_V1_3).
    ssl_profiles = optional(list(object({
      name                                 = string
      trusted_client_certificate_names     = optional(list(string), [])
      verify_client_certificate_issuer_dn  = optional(bool, false)
      verify_client_certificate_revocation = optional(string)
      ssl_policy = optional(object({
        policy_type          = optional(string)
        policy_name          = optional(string)
        min_protocol_version = optional(string)
        cipher_suites        = optional(list(string), [])
        disabled_protocols   = optional(list(string), [])
      }))
    })), [])

    # The gateway-wide TLS policy.
    ssl_policy = optional(object({
      policy_type          = optional(string)
      policy_name          = optional(string)
      min_protocol_version = optional(string)
      cipher_suites        = optional(list(string), [])
      disabled_protocols   = optional(list(string), [])
    }))

    # Redirect definitions; redirect_type arrives as
    # PERMANENT/FOUND/SEE_OTHER/TEMPORARY.
    redirect_configurations = optional(list(object({
      name                 = string
      redirect_type        = string
      target_listener_name = optional(string)
      target_url           = optional(string)
      include_path         = optional(bool, false)
      include_query_string = optional(bool, false)
    })), [])

    # Rewrite rule sets; url.components arrives as
    # PATH_ONLY/QUERY_STRING_ONLY.
    rewrite_rule_sets = optional(list(object({
      name = string
      rewrite_rules = list(object({
        name          = string
        rule_sequence = number
        conditions = optional(list(object({
          variable    = string
          pattern     = string
          ignore_case = optional(bool, false)
          negate      = optional(bool, false)
        })), [])
        request_header_configurations = optional(list(object({
          header_name  = string
          header_value = optional(string, "")
        })), [])
        response_header_configurations = optional(list(object({
          header_name  = string
          header_value = optional(string, "")
        })), [])
        url = optional(object({
          path         = optional(string)
          query_string = optional(string)
          components   = optional(string)
          reroute      = optional(bool, false)
        }))
      }))
    })), [])

    # Layer-4 (TCP/TLS proxy) trio; protocols arrive as TCP/TLS.
    listeners = optional(list(object({
      name                           = string
      frontend_ip_configuration_name = string
      frontend_port_name             = string
      protocol                       = string
      host_names                     = optional(list(string), [])
      ssl_certificate_name           = optional(string)
      ssl_profile_name               = optional(string)
    })), [])

    backends = optional(list(object({
      name                           = string
      port                           = number
      protocol                       = string
      client_ip_preservation_enabled = optional(bool, false)
      host_name                      = optional(string)
      probe_name                     = optional(string)
      timeout_in_seconds             = optional(number, 30)
      trusted_root_certificate_names = optional(list(string), [])
    })), [])

    routing_rules = optional(list(object({
      name                      = string
      listener_name             = string
      backend_address_pool_name = string
      backend_settings_name     = string
      priority                  = number
    })), [])

    # The WAF policy (WAF_V2 only), as a resolved ARM ID.
    firewall_policy_id                = optional(string)
    force_firewall_policy_association = optional(bool, false)

    # Gateway-wide custom error pages; status_code arrives as
    # HTTP_STATUS_NNN.
    custom_error_configurations = optional(list(object({
      status_code           = string
      custom_error_page_url = string
    })), [])

    # Private Link configurations.
    private_link_configurations = optional(list(object({
      name = string
      ip_configurations = list(object({
        name                          = string
        subnet_id                     = string
        private_ip_address            = optional(string)
        private_ip_address_allocation = string
        primary                       = optional(bool, false)
      }))
    })), [])

    # HTTP/2 for client connections (Azure's default is false).
    http2_enabled = optional(bool, false)

    # FIPS 140-2 validated crypto modules.
    fips_enabled = optional(bool, false)

    # Request/response buffering (both fields required when declared).
    global_configuration = optional(object({
      request_buffering_enabled  = bool
      response_buffering_enabled = bool
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
