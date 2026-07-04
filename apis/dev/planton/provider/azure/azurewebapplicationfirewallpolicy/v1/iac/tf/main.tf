# The regional Web Application Firewall policy. Custom rules run first
# (ascending priority), then the managed rule sets; policy settings govern
# enforcement mode and body inspection. Application Gateways attach the
# policy by referencing its ID -- gateway-wide, per listener, or per URL
# path rule -- so the policy carries no back-references of its own.
resource "azurerm_web_application_firewall_policy" "main" {
  name                = var.spec.policy_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  # Custom rules: IP/geo allowlists, header exceptions, rate limits. The
  # rate-limit trio is only sent for RATE_LIMIT_RULE rules (spec
  # validation pairs them; azurerm rejects strays).
  dynamic "custom_rules" {
    for_each = var.spec.custom_rules
    content {
      name      = custom_rules.value.name
      priority  = custom_rules.value.priority
      enabled   = custom_rules.value.enabled
      rule_type = local.rule_type_map[custom_rules.value.rule_type]
      action    = local.custom_rule_action_map[custom_rules.value.action]

      rate_limit_duration = (
        custom_rules.value.rate_limit_duration == null || custom_rules.value.rate_limit_duration == "" ? null :
        local.rate_limit_duration_map[custom_rules.value.rate_limit_duration]
      )
      rate_limit_threshold = custom_rules.value.rate_limit_threshold
      group_rate_limit_by = (
        custom_rules.value.group_rate_limit_by == null || custom_rules.value.group_rate_limit_by == "" ? null :
        local.group_rate_limit_by_map[custom_rules.value.group_rate_limit_by]
      )

      dynamic "match_conditions" {
        for_each = custom_rules.value.match_conditions
        content {
          operator           = local.match_operator_map[match_conditions.value.operator]
          match_values       = match_conditions.value.match_values
          negation_condition = match_conditions.value.negation_condition
          transforms         = [for transform in match_conditions.value.transforms : local.transform_map[transform]]

          dynamic "match_variables" {
            for_each = match_conditions.value.match_variables
            content {
              variable_name = local.match_variable_name_map[match_variables.value.variable_name]
              selector      = match_variables.value.selector
            }
          }
        }
      }
    }
  }

  # The managed rule configuration: Microsoft's curated sets, tuned with
  # per-rule overrides and scoped exclusions.
  managed_rules {
    dynamic "exclusion" {
      for_each = var.spec.managed_rules.exclusions
      content {
        match_variable          = local.exclusion_match_variable_map[exclusion.value.match_variable]
        selector_match_operator = local.selector_match_operator_map[exclusion.value.selector_match_operator]
        selector                = exclusion.value.selector

        dynamic "excluded_rule_set" {
          for_each = exclusion.value.excluded_rule_set != null ? [exclusion.value.excluded_rule_set] : []
          content {
            # Unspecified type applies OWASP -- azurerm's own default,
            # materialized so both engines send the same payload.
            type    = excluded_rule_set.value.type == null || excluded_rule_set.value.type == "" ? "OWASP" : local.managed_rule_set_type_map[excluded_rule_set.value.type]
            version = excluded_rule_set.value.version

            dynamic "rule_group" {
              for_each = excluded_rule_set.value.rule_groups
              content {
                rule_group_name = rule_group.value.rule_group_name
                excluded_rules  = rule_group.value.excluded_rules
              }
            }
          }
        }
      }
    }

    dynamic "managed_rule_set" {
      for_each = var.spec.managed_rules.managed_rule_sets
      content {
        type    = managed_rule_set.value.type == null || managed_rule_set.value.type == "" ? "OWASP" : local.managed_rule_set_type_map[managed_rule_set.value.type]
        version = managed_rule_set.value.version

        dynamic "rule_group_override" {
          for_each = managed_rule_set.value.rule_group_overrides
          content {
            rule_group_name = rule_group_override.value.rule_group_name

            dynamic "rule" {
              for_each = rule_group_override.value.rules
              content {
                id      = rule.value.id
                enabled = rule.value.enabled
                action = (
                  rule.value.action == null || rule.value.action == "" ? null :
                  local.rule_override_action_map[rule.value.action]
                )
              }
            }
          }
        }
      }
    }
  }

  # Enforcement mode and body-inspection dials. Omitting the block applies
  # Azure's defaults (enabled, Prevention, body check on, 128 KB limits).
  dynamic "policy_settings" {
    for_each = var.spec.policy_settings != null ? [var.spec.policy_settings] : []
    content {
      enabled = policy_settings.value.enabled
      mode = (
        policy_settings.value.mode == null || policy_settings.value.mode == "" ? "Prevention" :
        local.mode_map[policy_settings.value.mode]
      )
      request_body_check                        = policy_settings.value.request_body_check
      request_body_enforcement                  = policy_settings.value.request_body_enforcement
      request_body_inspect_limit_in_kb          = policy_settings.value.request_body_inspect_limit_in_kb
      max_request_body_size_in_kb               = policy_settings.value.max_request_body_size_in_kb
      file_upload_enforcement                   = policy_settings.value.file_upload_enforcement
      file_upload_limit_in_mb                   = policy_settings.value.file_upload_limit_in_mb
      js_challenge_cookie_expiration_in_minutes = policy_settings.value.js_challenge_cookie_expiration_in_minutes

      dynamic "log_scrubbing" {
        for_each = policy_settings.value.log_scrubbing != null ? [policy_settings.value.log_scrubbing] : []
        content {
          enabled = log_scrubbing.value.enabled

          dynamic "rule" {
            for_each = log_scrubbing.value.rules
            content {
              enabled        = rule.value.enabled
              match_variable = local.scrubbing_match_variable_map[rule.value.match_variable]
              # Unspecified operator means Equals (azurerm's default).
              selector_match_operator = (
                rule.value.selector_match_operator == null || rule.value.selector_match_operator == "" ? "Equals" :
                local.selector_match_operator_map[rule.value.selector_match_operator]
              )
              selector = rule.value.selector
            }
          }
        }
      }
    }
  }

  tags = local.final_tags
}
