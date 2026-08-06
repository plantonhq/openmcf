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
  description = "Azure Front Door rule set specification"
  type = object({
    # The Front Door profile the rule set lives in, by ARM ID.
    # References are resolved to a literal ID by the platform before the
    # module runs. ForceNew.
    profile_id = string

    # The rule set's name -- unique within the profile; letters and
    # digits only. ForceNew (replaces the set AND every rule in it).
    rule_set_name = string

    # The ordered delivery policy. Each rule is its own Azure resource
    # keyed by name; evaluation follows each rule's `order`, ascending.
    rules = optional(list(object({
      # Unique within the set. ForceNew per rule.
      name = string

      # Evaluation position; 0 or greater. tfvars drops zero values, so
      # the attribute defaults the documented 0 back in.
      order = optional(number, 0)

      # CONTINUE (default) / STOP -- spec enum value names.
      behavior_on_match = optional(string)

      # Match conditions, grouped by the request attribute inspected.
      # ALL conditions must match (values within one condition OR
      # together). Absent conditions = the rule applies to every
      # request. Operators and transforms arrive as the spec enum's
      # FULL value names and are mapped to ARM's casing in locals.
      conditions = optional(object({
        remote_address = optional(list(object({
          # ANY / GEO_MATCH / IP_MATCH; absent means IP_MATCH.
          operator         = optional(string)
          negate_condition = optional(bool)
          match_values     = optional(list(string))
        })), [])
        request_method = optional(list(object({
          negate_condition = optional(bool)
          # GET / HEAD / OPTIONS / TRACE / POST / PUT / DELETE --
          # already ARM's vocabulary, passed through.
          match_values = list(string)
        })), [])
        query_string = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        post_args = optional(list(object({
          post_args_name   = string
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        request_uri = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        request_header = optional(list(object({
          header_name      = string
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        request_body = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = list(string)
          transforms       = optional(list(string))
        })), [])
        request_scheme = optional(list(object({
          negate_condition = optional(bool)
          # HTTP / HTTPS; absent means HTTP (Azure's default).
          match_value = optional(string)
        })), [])
        url_path = optional(list(object({
          # The standard operator set plus WILDCARD.
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        url_file_extension = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = list(string)
          transforms       = optional(list(string))
        })), [])
        url_filename = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        http_version = optional(list(object({
          negate_condition = optional(bool)
          # "2.0" / "1.1" / "1.0" / "0.9" -- ARM's own values.
          match_values = list(string)
        })), [])
        cookies = optional(list(object({
          cookie_name      = string
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        is_device = optional(list(object({
          negate_condition = optional(bool)
          # Desktop / Mobile -- ARM's own values.
          match_value = string
        })), [])
        socket_address = optional(list(object({
          # ANY / IP_MATCH; absent means IP_MATCH.
          operator         = optional(string)
          negate_condition = optional(bool)
          match_values     = optional(list(string))
        })), [])
        client_port = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
        })), [])
        server_port = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          # "80" / "443" -- the only ports Front Door serves.
          match_values = list(string)
        })), [])
        host_name = optional(list(object({
          operator         = string
          negate_condition = optional(bool)
          match_values     = optional(list(string))
          transforms       = optional(list(string))
        })), [])
        ssl_protocol = optional(list(object({
          negate_condition = optional(bool)
          # TLSv1 / TLSv1.1 / TLSv1.2 -- ARM's own values.
          match_values = list(string)
        })), [])
      }))

      # The actions applied on match. The spec guarantees at least one
      # action, at most five, and redirect XOR rewrite.
      actions = object({
        url_redirect = optional(object({
          # MOVED / FOUND / TEMPORARY_REDIRECT / PERMANENT_REDIRECT.
          redirect_type = string
          # MATCH_REQUEST (default) / HTTP_ONLY / HTTPS_ONLY -- mapped
          # to the redirect dialect (Http/Https/MatchRequest).
          redirect_protocol = optional(string)
          # Empty preserves the incoming host/path/query/fragment.
          destination_hostname = optional(string, "")
          destination_path     = optional(string, "")
          query_string         = optional(string, "")
          destination_fragment = optional(string, "")
        }))
        url_rewrite = optional(object({
          source_pattern          = string
          destination             = string
          preserve_unmatched_path = optional(bool)
        }))
        request_headers = optional(list(object({
          # APPEND / OVERWRITE / DELETE.
          header_action = string
          header_name   = string
          # Required for APPEND/OVERWRITE, empty for DELETE
          # (spec-enforced).
          value = optional(string)
        })), [])
        response_headers = optional(list(object({
          header_action = string
          header_name   = string
          value         = optional(string)
        })), [])
        route_configuration_override = optional(object({
          # Set together with forwarding_protocol, or not at all.
          origin_group_id = optional(string)
          # MATCH_REQUEST / HTTP_ONLY / HTTPS_ONLY -- mapped to the
          # override dialect (HttpOnly/HttpsOnly/MatchRequest).
          forwarding_protocol = optional(string)
          # HONOR_ORIGIN / OVERRIDE_ALWAYS / OVERRIDE_IF_ORIGIN_MISSING
          # / DISABLED -- required by the spec.
          cache_behavior = string
          cache_duration = optional(string)
          # IGNORE_QUERY_STRING / USE_QUERY_STRING /
          # IGNORE_SPECIFIED_QUERY_STRINGS /
          # INCLUDE_SPECIFIED_QUERY_STRINGS.
          query_string_caching_behavior = optional(string)
          query_string_parameters       = optional(list(string))
          compression_enabled           = optional(bool)
        }))
      })
    })), [])
  })
}
