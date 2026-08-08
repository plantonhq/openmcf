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
# defaults apply (behaviour_on_match Continue, redirect scheme
# MatchRequest).
#
# The provider expresses negation as `Not`-prefixed operator values
# (NotEqual, NotIPMatch, ...) instead of a separate flag, and requires an
# operator on every condition -- including the equality-only ones (method,
# scheme, HTTP version, device, TLS). The spec keeps negation orthogonal
# (a `negate_condition` bool per condition), so every operator below is
# composed as: optional "Not" prefix + the mapped base operator.
resource "azurerm_cdn_frontdoor_rule" "main" {
  for_each = local.rules_by_name

  name                      = each.value.name
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.main.id
  order                     = coalesce(each.value.order, 0)
  behaviour_on_match        = each.value.behavior_on_match != null ? local.behavior_on_match_map[each.value.behavior_on_match] : null

  actions {
    dynamic "url_redirect" {
      for_each = each.value.actions.url_redirect != null ? [each.value.actions.url_redirect] : []
      content {
        redirect_type     = local.redirect_type_map[url_redirect.value.redirect_type]
        redirect_protocol = url_redirect.value.redirect_protocol != null ? local.redirect_protocol_map[url_redirect.value.redirect_protocol] : null
        # Omitted attributes preserve the corresponding part of the
        # incoming request (host, path, query, fragment); the provider
        # rejects empty strings, so empties are sent as null.
        destination_host_name = url_redirect.value.destination_hostname != "" ? url_redirect.value.destination_hostname : null
        destination_path      = url_redirect.value.destination_path != "" ? url_redirect.value.destination_path : null
        query_string          = url_redirect.value.query_string != "" ? url_redirect.value.query_string : null
        destination_fragment  = url_redirect.value.destination_fragment != "" ? url_redirect.value.destination_fragment : null
      }
    }

    dynamic "url_rewrite" {
      for_each = each.value.actions.url_rewrite != null ? [each.value.actions.url_rewrite] : []
      content {
        source_pattern                  = url_rewrite.value.source_pattern
        destination_path                = url_rewrite.value.destination
        preserve_unmatched_path_enabled = coalesce(url_rewrite.value.preserve_unmatched_path, false)
      }
    }

    dynamic "modify_request_header" {
      for_each = each.value.actions.request_headers
      content {
        operator    = local.header_action_map[modify_request_header.value.header_action]
        header_name = modify_request_header.value.header_name
        # DELETE carries no value (spec-enforced); the provider rejects
        # an empty value on Append/Overwrite, so null is sent instead.
        header_value = modify_request_header.value.value != null && modify_request_header.value.value != "" ? modify_request_header.value.value : null
      }
    }

    dynamic "modify_response_header" {
      for_each = each.value.actions.response_headers
      content {
        operator     = local.header_action_map[modify_response_header.value.header_action]
        header_name  = modify_response_header.value.header_name
        header_value = modify_response_header.value.value != null && modify_response_header.value.value != "" ? modify_response_header.value.value : null
      }
    }

    dynamic "route_configuration_override" {
      for_each = each.value.actions.route_configuration_override != null ? [each.value.actions.route_configuration_override] : []
      content {
        # WHERE the request goes: only rendered when an origin override
        # is chosen (the spec pairs origin_group_id with
        # forwarding_protocol -- the override dialect of the shared
        # protocol enum: HttpOnly/HttpsOnly vs the redirect's
        # Http/Https).
        dynamic "origin_group" {
          for_each = route_configuration_override.value.origin_group_id != null ? [route_configuration_override.value] : []
          content {
            cdn_frontdoor_origin_group_id = origin_group.value.origin_group_id
            forwarding_protocol           = local.override_forwarding_protocol_map[origin_group.value.forwarding_protocol]
          }
        }

        # HOW it caches: the provider requires the block -- mirroring
        # the spec, where every override makes an explicit cache
        # decision (DISABLED is the "no caching" choice).
        caching {
          behaviour               = local.cache_behavior_map[route_configuration_override.value.cache_behavior]
          duration                = route_configuration_override.value.cache_duration != null && route_configuration_override.value.cache_duration != "" ? route_configuration_override.value.cache_duration : null
          query_string_behaviour  = route_configuration_override.value.query_string_caching_behavior != null ? local.query_string_caching_behavior_map[route_configuration_override.value.query_string_caching_behavior] : null
          query_string_parameters = route_configuration_override.value.query_string_parameters
          compression_enabled     = route_configuration_override.value.compression_enabled
        }
      }
    }
  }

  dynamic "conditions" {
    for_each = each.value.conditions != null ? [each.value.conditions] : []
    content {
      dynamic "remote_address" {
        for_each = conditions.value.remote_address
        content {
          # Absent operator means IP_MATCH -- materialized here because
          # tfvars drops unset fields.
          operator = "${coalesce(remote_address.value.negate_condition, false) ? "Not" : ""}${remote_address.value.operator != null ? local.operator_map[remote_address.value.operator] : "IPMatch"}"
          values   = remote_address.value.match_values
        }
      }

      dynamic "request_method" {
        for_each = conditions.value.request_method
        content {
          # Equality is the only comparison ARM supports here; the
          # provider still requires it spelled out.
          operator = coalesce(request_method.value.negate_condition, false) ? "NotEqual" : "Equal"
          # Methods are already ARM's vocabulary -- passed through.
          values = request_method.value.match_values
        }
      }

      dynamic "query_string" {
        for_each = conditions.value.query_string
        content {
          operator   = "${coalesce(query_string.value.negate_condition, false) ? "Not" : ""}${local.operator_map[query_string.value.operator]}"
          values     = query_string.value.match_values
          transforms = query_string.value.transforms != null ? [for t in query_string.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "post_argument" {
        for_each = conditions.value.post_args
        content {
          name       = post_argument.value.post_args_name
          operator   = "${coalesce(post_argument.value.negate_condition, false) ? "Not" : ""}${local.operator_map[post_argument.value.operator]}"
          values     = post_argument.value.match_values
          transforms = post_argument.value.transforms != null ? [for t in post_argument.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_url" {
        for_each = conditions.value.request_uri
        content {
          operator   = "${coalesce(request_url.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_url.value.operator]}"
          values     = request_url.value.match_values
          transforms = request_url.value.transforms != null ? [for t in request_url.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_header" {
        for_each = conditions.value.request_header
        content {
          name       = request_header.value.header_name
          operator   = "${coalesce(request_header.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_header.value.operator]}"
          values     = request_header.value.match_values
          transforms = request_header.value.transforms != null ? [for t in request_header.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_body" {
        for_each = conditions.value.request_body
        content {
          operator   = "${coalesce(request_body.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_body.value.operator]}"
          values     = request_body.value.match_values
          transforms = request_body.value.transforms != null ? [for t in request_body.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_scheme" {
        for_each = conditions.value.request_scheme
        content {
          operator = coalesce(request_scheme.value.negate_condition, false) ? "NotEqual" : "Equal"
          # Single scheme; absent means HTTP (Azure's default),
          # materialized because tfvars drops unset fields.
          values = [coalesce(request_scheme.value.match_value, "HTTP")]
        }
      }

      dynamic "request_path" {
        for_each = conditions.value.url_path
        content {
          operator   = "${coalesce(request_path.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_path.value.operator]}"
          values     = request_path.value.match_values
          transforms = request_path.value.transforms != null ? [for t in request_path.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_file_extension" {
        for_each = conditions.value.url_file_extension
        content {
          operator   = "${coalesce(request_file_extension.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_file_extension.value.operator]}"
          values     = request_file_extension.value.match_values
          transforms = request_file_extension.value.transforms != null ? [for t in request_file_extension.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "request_filename" {
        for_each = conditions.value.url_filename
        content {
          operator   = "${coalesce(request_filename.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_filename.value.operator]}"
          values     = request_filename.value.match_values
          transforms = request_filename.value.transforms != null ? [for t in request_filename.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "http_version" {
        for_each = conditions.value.http_version
        content {
          operator = coalesce(http_version.value.negate_condition, false) ? "NotEqual" : "Equal"
          # Versions are already ARM's vocabulary -- passed through.
          values = http_version.value.match_values
        }
      }

      dynamic "request_cookies" {
        for_each = conditions.value.cookies
        content {
          name       = request_cookies.value.cookie_name
          operator   = "${coalesce(request_cookies.value.negate_condition, false) ? "Not" : ""}${local.operator_map[request_cookies.value.operator]}"
          values     = request_cookies.value.match_values
          transforms = request_cookies.value.transforms != null ? [for t in request_cookies.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "device_type" {
        for_each = conditions.value.is_device
        content {
          operator = coalesce(device_type.value.negate_condition, false) ? "NotEqual" : "Equal"
          # Single device class; already ARM's vocabulary.
          values = [device_type.value.match_value]
        }
      }

      dynamic "socket_address" {
        for_each = conditions.value.socket_address
        content {
          operator = "${coalesce(socket_address.value.negate_condition, false) ? "Not" : ""}${socket_address.value.operator != null ? local.operator_map[socket_address.value.operator] : "IPMatch"}"
          values   = socket_address.value.match_values
        }
      }

      dynamic "client_port" {
        for_each = conditions.value.client_port
        content {
          operator = "${coalesce(client_port.value.negate_condition, false) ? "Not" : ""}${local.operator_map[client_port.value.operator]}"
          values   = client_port.value.match_values
        }
      }

      dynamic "server_port" {
        for_each = conditions.value.server_port
        content {
          operator = "${coalesce(server_port.value.negate_condition, false) ? "Not" : ""}${local.operator_map[server_port.value.operator]}"
          values   = server_port.value.match_values
        }
      }

      dynamic "host_name" {
        for_each = conditions.value.host_name
        content {
          operator   = "${coalesce(host_name.value.negate_condition, false) ? "Not" : ""}${local.operator_map[host_name.value.operator]}"
          values     = host_name.value.match_values
          transforms = host_name.value.transforms != null ? [for t in host_name.value.transforms : local.transform_map[t]] : null
        }
      }

      dynamic "ssl_protocol" {
        for_each = conditions.value.ssl_protocol
        content {
          operator = coalesce(ssl_protocol.value.negate_condition, false) ? "NotEqual" : "Equal"
          # TLS versions are already ARM's vocabulary -- passed through.
          values = ssl_protocol.value.match_values
        }
      }
    }
  }
}
