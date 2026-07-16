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
  description = "Azure Web Application Firewall policy specification"
  type = object({
    # The Azure region the policy lives in (must match the gateways it
    # attaches to).
    region = string

    # The resource group the policy lives in. References are resolved to a
    # literal name by the platform before the module runs.
    resource_group = string

    # The policy's name, unique within the resource group.
    policy_name = string

    # Custom rules, evaluated before the managed sets by ascending
    # priority. Enums arrive as the spec enums' FULL name strings
    # (MATCH_RULE, RATE_LIMIT_RULE; ALLOW/BLOCK/LOG/JS_CHALLENGE;
    # ONE_MIN/FIVE_MINS; CLIENT_ADDR/...; REMOTE_ADDR/...; ANY/IP_MATCH/...;
    # LOWERCASE/...) and are mapped to ARM's values in locals.
    custom_rules = optional(list(object({
      name                 = optional(string)
      priority             = number
      enabled              = optional(bool, true)
      rule_type            = string
      action               = string
      rate_limit_duration  = optional(string)
      rate_limit_threshold = optional(number)
      group_rate_limit_by  = optional(string)
      match_conditions = list(object({
        match_variables = list(object({
          variable_name = string
          selector      = optional(string)
        }))
        operator           = string
        match_values       = optional(list(string), [])
        negation_condition = optional(bool, false)
        transforms         = optional(list(string), [])
      }))
    })), [])

    # The managed (Microsoft-curated) rule configuration -- always present
    # (spec validation requires it).
    managed_rules = object({
      exclusions = optional(list(object({
        match_variable          = string
        selector_match_operator = string
        selector                = string
        excluded_rule_set = optional(object({
          type    = optional(string)
          version = optional(string, "3.2")
          rule_groups = optional(list(object({
            rule_group_name = string
            excluded_rules  = optional(list(string), [])
          })), [])
        }))
      })), [])
      managed_rule_sets = list(object({
        type    = optional(string)
        version = string
        rule_group_overrides = optional(list(object({
          rule_group_name = string
          rules = list(object({
            id      = string
            enabled = optional(bool, false)
            action  = optional(string)
          }))
        })), [])
      }))
    })

    # Enforcement mode and body-inspection dials; mode arrives as the spec
    # enum's name string (PREVENTION/DETECTION).
    policy_settings = optional(object({
      enabled                                   = optional(bool, true)
      mode                                      = optional(string)
      request_body_check                        = optional(bool, true)
      request_body_enforcement                  = optional(bool, true)
      request_body_inspect_limit_in_kb          = optional(number, 128)
      max_request_body_size_in_kb               = optional(number, 128)
      file_upload_enforcement                   = optional(bool)
      file_upload_limit_in_mb                   = optional(number, 100)
      js_challenge_cookie_expiration_in_minutes = optional(number, 30)
      log_scrubbing = optional(object({
        enabled = optional(bool, true)
        rules = list(object({
          enabled                 = optional(bool, true)
          match_variable          = string
          selector_match_operator = optional(string)
          selector                = optional(string)
        }))
      }))
    }))

    # Free-form user tags, merged over the metadata-derived tags (user
    # tags win on key collision).
    tags = optional(map(string), {})
  })
}
