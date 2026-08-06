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
  description = "Specification for the Cloud Armor security policy"
  type = object({
    # StringValueOrRef fields arrive as PLAIN STRINGS: the tfvars converter
    # flattens refs before the module ever sees them.
    project_id = optional(string, "")

    policy_name = optional(string, "")
    description = optional(string, "")
    type        = optional(string, "")

    adaptive_protection_config = optional(object({
      enable_layer_7_ddos_defense = optional(bool, false)
      rule_visibility             = optional(string, "")
      threshold_configs = optional(list(object({
        name                                    = string
        auto_deploy_confidence_threshold        = optional(number)
        auto_deploy_impacted_baseline_threshold = optional(number)
        auto_deploy_load_threshold              = optional(number)
        auto_deploy_expiration_sec              = optional(number)
        detection_absolute_qps                  = optional(number)
        detection_load_threshold                = optional(number)
        detection_relative_to_baseline_qps      = optional(number)
        traffic_granularity_configs = optional(list(object({
          type                     = string
          value                    = optional(string, "")
          enable_each_unique_value = optional(bool, false)
        })), [])
      })), [])
    }), null)

    advanced_options_config = optional(object({
      json_parsing            = optional(string, "")
      log_level               = optional(string, "")
      user_ip_request_headers = optional(list(string), [])
      json_custom_config = optional(object({
        content_types = list(string)
      }), null)
    }), null)

    recaptcha_options_config = optional(object({
      redirect_site_key = string
    }), null)

    rules = optional(list(object({
      action      = string
      priority    = number
      description = optional(string, "")
      preview     = optional(bool, false)

      match = object({
        versioned_expr = optional(string, "")
        src_ip_ranges  = optional(list(string), [])
        expression     = optional(string, "")
        expr_options = optional(object({
          action_token_site_keys  = optional(list(string), [])
          session_token_site_keys = optional(list(string), [])
        }), null)
      })

      rate_limit_options = optional(object({
        conform_action      = string
        exceed_action       = string
        enforce_on_key      = optional(string, "")
        enforce_on_key_name = optional(string, "")
        enforce_on_key_configs = optional(list(object({
          enforce_on_key_type = string
          enforce_on_key_name = optional(string, "")
        })), [])
        rate_limit_threshold = object({
          count        = number
          interval_sec = number
        })
        ban_threshold = optional(object({
          count        = number
          interval_sec = number
        }), null)
        ban_duration_sec = optional(number, 0)
        exceed_redirect_options = optional(object({
          type   = string
          target = optional(string, "")
        }), null)
      }), null)

      redirect_options = optional(object({
        type   = string
        target = optional(string, "")
      }), null)

      header_action = optional(object({
        request_headers_to_adds = list(object({
          header_name  = string
          header_value = optional(string, "")
        }))
      }), null)

      preconfigured_waf_config = optional(object({
        exclusions = list(object({
          target_rule_set = string
          target_rule_ids = optional(list(string), [])
          request_headers = optional(list(object({
            operator = string
            value    = optional(string, "")
          })), [])
          request_cookies = optional(list(object({
            operator = string
            value    = optional(string, "")
          })), [])
          request_uris = optional(list(object({
            operator = string
            value    = optional(string, "")
          })), [])
          request_query_params = optional(list(object({
            operator = string
            value    = optional(string, "")
          })), [])
        }))
      }), null)
    })), [])
  })
}
