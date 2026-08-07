# The rule set -- the named container routes attach by ARM id. It
# carries no properties of its own; the delivery policy lives in the
# rules below. No Azure tags: ARM does not support tags on rule sets.
resource "azurerm_cdn_frontdoor_rule_set" "main" {
  name                     = var.spec.rule_set_name
  cdn_frontdoor_profile_id = var.spec.profile_id
}

# One provider resource per rule (ARM models rules as children of the
# set), keyed by the rule's spec-unique name. Evaluation position is the
# rule's own `order`, so resource ordering here carries no meaning.
#
# The conditions/actions blocks below are exhaustive over the provider's
# surface: every typed condition and action the spec models renders
# through a dynamic block, mapped from the spec enum vocabulary to ARM's
# casing in locals. Absent optional enums are NOT sent, so Azure's own
# defaults apply (behavior_on_match Continue, redirect scheme
# MatchRequest) -- except the address conditions' operator, whose
# documented IP_MATCH default is materialized here because tfvars drops
# unset fields.
resource "azurerm_cdn_frontdoor_rule" "main" {
  for_each = local.rules_by_name

  name                      = each.value.name
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.main.id
  order                     = coalesce(each.value.order, 0)
  behavior_on_match         = each.value.behavior_on_match != null ? local.behavior_on_match_map[each.value.behavior_on_match] : null

  actions {
    dynamic "url_redirect_action" {
      for_each = each.value.actions.url_redirect != null ? [each.value.actions.url_redirect] : []
      content {
        redirect_type     = local.redirect_type_map[url_redirect_action.value.redirect_type]
        redirect_protocol = url_redirect_action.value.redirect_protocol != null ? local.redirect_protocol_map[url_redirect_action.value.redirect_protocol] : null
        # Empty strings are meaningful: they preserve the corresponding
        # part of the incoming request (host, path, query, fragment).
        destination_hostname = url_redirect_action.value.destination_hostname
        destination_path     = url_redirect_action.value.destination_path
        query_string         = url_redirect_action.value.query_string
        destination_fragment = url_redirect_action.value.destination_fragment
      }
    }

    dynamic "url_rewrite_action" {
      for_each = each.value.actions.url_rewrite != null ? [each.value.actions.url_rewrite] : []
      content {
        source_pattern          = url_rewrite_action.value.source_pattern
        destination             = url_rewrite_action.value.destination
        preserve_unmatched_path = coalesce(url_rewrite_action.value.preserve_unmatched_path, false)
      }
    }

    dynamic "request_header_action" {
      for_each = each.value.actions.request_headers
      content {
        header_action = local.header_action_map[request_header_action.value.header_action]
        header_name   = request_header_action.value.header_name
        # DELETE carries no value (spec-enforced); the provider rejects
        # an empty value on Append/Overwrite, so null is sent instead.
        value = request_header_action.value.value != null && request_header_action.value.value != "" ? request_header_action.value.value : null
      }
    }

    dynamic "response_header_action" {
      for_each = each.value.actions.response_headers
      content {
        header_action = local.header_action_map[response_header_action.value.header_action]
        header_name   = response_header_action.value.header_name
        value         = response_header_action.value.value != null && response_header_action.value.value != "" ? response_header_action.value.value : null
      }
    }

    dynamic "route_configuration_override_action" {
      for_each = each.value.actions.route_configuration_override != null ? [each.value.actions.route_configuration_override] : []
      content {
        cdn_frontdoor_origin_group_id = route_configuration_override_action.value.origin_group_id
        # The override dialect of the shared protocol enum
        # (HttpOnly/HttpsOnly vs the redirect's Http/Https).
        forwarding_protocol           = route_configuration_override_action.value.forwarding_protocol != null ? local.override_forwarding_protocol_map[route_configuration_override_action.value.forwarding_protocol] : null
        cache_behavior                = local.cache_behavior_map[route_configuration_override_action.value.cache_behavior]
        cache_duration                = route_configuration_override_action.value.cache_duration
        query_string_caching_behavior = route_configuration_override_action.value.query_string_caching_behavior != null ? local.query_string_caching_behavior_map[route_configuration_override_action.value.query_string_caching_behavior] : null
        query_string_parameters       = route_configuration_override_action.value.query_string_parameters
        compression_enabled           = route_configuration_override_action.value.compression_enabled
      }
    }
  }

  dynamic "conditions" {
    for_each = each.value.conditions != null ? [each.value.conditions] : []
    content {
      dynamic "remote_address_condition" {
        for_each = conditions.value.remote_address
        content {
          # Absent operator means IP_MATCH -- materialized here because
          # tfvars drops unset fields.
          operator         = remote_address_condition.value.operator != null ? local.operator_map[remote_address_condition.value.operator] : "IPMatch"
          negate_condition = coalesce(remote_address_condition.value.negate_condition, false)
          match_values     = remote_address_condition.value.match_values
        }
      }

      dynamic "request_method_condition" {
        for_each = conditions.value.request_method
        content {
          # Methods are already ARM's vocabulary -- passed through.
          negate_condition = coalesce(request_method_condition.value.negate_condition, false)
          match_values     = request_method_condition.value.match_values
        }
      }

      dynamic "query_string_condition" {
        for_each = conditions.value.query_string
        content {
          operator         = local.operator_map[query_string_condition.value.operator]
          negate_condition = coalesce(query_string_condition.value.negate_condition, false)
          match_values     = query_string_condition.value.match_values
          transforms       = query_string_condition.value.transforms != null ? [for t in query_string_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "post_args_condition" {
        for_each = conditions.value.post_args
        content {
          post_args_name   = post_args_condition.value.post_args_name
          operator         = local.operator_map[post_args_condition.value.operator]
          negate_condition = coalesce(post_args_condition.value.negate_condition, false)
          match_values     = post_args_condition.value.match_values
          transforms       = post_args_condition.value.transforms != null ? [for t in post_args_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_uri_condition" {
        for_each = conditions.value.request_uri
        content {
          operator         = local.operator_map[request_uri_condition.value.operator]
          negate_condition = coalesce(request_uri_condition.value.negate_condition, false)
          match_values     = request_uri_condition.value.match_values
          transforms       = request_uri_condition.value.transforms != null ? [for t in request_uri_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_header_condition" {
        for_each = conditions.value.request_header
        content {
          header_name      = request_header_condition.value.header_name
          operator         = local.operator_map[request_header_condition.value.operator]
          negate_condition = coalesce(request_header_condition.value.negate_condition, false)
          match_values     = request_header_condition.value.match_values
          transforms       = request_header_condition.value.transforms != null ? [for t in request_header_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_body_condition" {
        for_each = conditions.value.request_body
        content {
          operator         = local.operator_map[request_body_condition.value.operator]
          negate_condition = coalesce(request_body_condition.value.negate_condition, false)
          match_values     = request_body_condition.value.match_values
          transforms       = request_body_condition.value.transforms != null ? [for t in request_body_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_scheme_condition" {
        for_each = conditions.value.request_scheme
        content {
          negate_condition = coalesce(request_scheme_condition.value.negate_condition, false)
          # Single scheme; absent means HTTP (Azure's default),
          # materialized because tfvars drops unset fields.
          match_values = [coalesce(request_scheme_condition.value.match_value, "HTTP")]
        }
      }

      dynamic "url_path_condition" {
        for_each = conditions.value.url_path
        content {
          operator         = local.operator_map[url_path_condition.value.operator]
          negate_condition = coalesce(url_path_condition.value.negate_condition, false)
          match_values     = url_path_condition.value.match_values
          transforms       = url_path_condition.value.transforms != null ? [for t in url_path_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "url_file_extension_condition" {
        for_each = conditions.value.url_file_extension
        content {
          operator         = local.operator_map[url_file_extension_condition.value.operator]
          negate_condition = coalesce(url_file_extension_condition.value.negate_condition, false)
          match_values     = url_file_extension_condition.value.match_values
          transforms       = url_file_extension_condition.value.transforms != null ? [for t in url_file_extension_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "url_filename_condition" {
        for_each = conditions.value.url_filename
        content {
          operator         = local.operator_map[url_filename_condition.value.operator]
          negate_condition = coalesce(url_filename_condition.value.negate_condition, false)
          match_values     = url_filename_condition.value.match_values
          transforms       = url_filename_condition.value.transforms != null ? [for t in url_filename_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "http_version_condition" {
        for_each = conditions.value.http_version
        content {
          # Versions are already ARM's vocabulary -- passed through.
          negate_condition = coalesce(http_version_condition.value.negate_condition, false)
          match_values     = http_version_condition.value.match_values
        }
      }

      dynamic "cookies_condition" {
        for_each = conditions.value.cookies
        content {
          cookie_name      = cookies_condition.value.cookie_name
          operator         = local.operator_map[cookies_condition.value.operator]
          negate_condition = coalesce(cookies_condition.value.negate_condition, false)
          match_values     = cookies_condition.value.match_values
          transforms       = cookies_condition.value.transforms != null ? [for t in cookies_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "is_device_condition" {
        for_each = conditions.value.is_device
        content {
          negate_condition = coalesce(is_device_condition.value.negate_condition, false)
          # Single device class; already ARM's vocabulary.
          match_values = [is_device_condition.value.match_value]
        }
      }

      dynamic "socket_address_condition" {
        for_each = conditions.value.socket_address
        content {
          operator         = socket_address_condition.value.operator != null ? local.operator_map[socket_address_condition.value.operator] : "IPMatch"
          negate_condition = coalesce(socket_address_condition.value.negate_condition, false)
          match_values     = socket_address_condition.value.match_values
        }
      }

      dynamic "client_port_condition" {
        for_each = conditions.value.client_port
        content {
          operator         = local.operator_map[client_port_condition.value.operator]
          negate_condition = coalesce(client_port_condition.value.negate_condition, false)
          match_values     = client_port_condition.value.match_values
        }
      }

      dynamic "server_port_condition" {
        for_each = conditions.value.server_port
        content {
          operator         = local.operator_map[server_port_condition.value.operator]
          negate_condition = coalesce(server_port_condition.value.negate_condition, false)
          match_values     = server_port_condition.value.match_values
        }
      }

      dynamic "host_name_condition" {
        for_each = conditions.value.host_name
        content {
          operator         = local.operator_map[host_name_condition.value.operator]
          negate_condition = coalesce(host_name_condition.value.negate_condition, false)
          match_values     = host_name_condition.value.match_values
          transforms       = host_name_condition.value.transforms != null ? [for t in host_name_condition.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "ssl_protocol_condition" {
        for_each = conditions.value.ssl_protocol
        content {
          # TLS versions are already ARM's vocabulary -- passed through.
          negate_condition = coalesce(ssl_protocol_condition.value.negate_condition, false)
          match_values     = ssl_protocol_condition.value.match_values
        }
      }
    }
  }
}
