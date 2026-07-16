# The Front Door WAF policy -- a GLOBAL resource (ARM fixes its
# location; the provider sends no region) scoped to a resource group,
# NOT nested under a profile. It enforces nothing on its own: an
# azurerm_cdn_frontdoor_security_policy (the AzureFrontDoorSecurityPolicy
# kind) associates it with a profile's domains. The policy's sku must
# match the sku of every profile it gets associated with.
resource "azurerm_cdn_frontdoor_firewall_policy" "main" {
  name                = var.spec.policy_name
  resource_group_name = var.spec.resource_group
  sku_name            = local.sku_name
  mode                = local.mode_map[var.spec.mode]

  # Provider defaults (both true) apply when the spec is silent -- the
  # tfvars wire format drops unset fields, so null lets the provider
  # default through.
  enabled                    = var.spec.enabled
  request_body_check_enabled = var.spec.request_body_check_enabled

  # Sent only when set: an empty redirect_url would fail the provider's
  # URL validator, and a zero status code its in-list validator.
  redirect_url                      = var.spec.redirect_url != null && var.spec.redirect_url != "" ? var.spec.redirect_url : null
  custom_block_response_status_code = var.spec.custom_block_response_status_code
  custom_block_response_body        = var.spec.custom_block_response_body != null && var.spec.custom_block_response_body != "" ? var.spec.custom_block_response_body : null

  # The JS-challenge and CAPTCHA lifetimes exist only on Premium: Azure
  # ALWAYS enables both policies there (rejecting them on Standard), so
  # on Premium the module pins the documented default of 30 minutes when
  # the spec is silent -- sending nothing would leave the value to drift
  # with Azure's server-side default instead of the declared contract.
  # The spec's Premium-only CELs keep the Standard path from ever
  # carrying a value.
  js_challenge_cookie_expiration_in_minutes = local.is_premium ? coalesce(var.spec.js_challenge_cookie_expiration_in_minutes, 30) : null
  captcha_cookie_expiration_in_minutes      = local.is_premium ? coalesce(var.spec.captcha_cookie_expiration_in_minutes, 30) : null

  # Custom rules, evaluated before the managed sets by ascending
  # priority. Enum vocabulary maps from the spec's value names to ARM's
  # casing in locals.
  dynamic "custom_rule" {
    for_each = local.custom_rules
    content {
      name   = custom_rule.value.name
      type   = local.custom_rule_type_map[custom_rule.value.rule_type]
      action = local.custom_rule_action_map[custom_rule.value.action]

      # Provider defaults (enabled true, priority 1, rate-limit window
      # 1 / threshold 10) apply when the spec is silent. The rate-limit
      # pair is harmless on MatchRule rules (the provider always sends
      # its defaults; ARM ignores them there).
      enabled                        = custom_rule.value.enabled
      priority                       = custom_rule.value.priority
      rate_limit_duration_in_minutes = custom_rule.value.rate_limit_duration_in_minutes
      rate_limit_threshold           = custom_rule.value.rate_limit_threshold

      dynamic "match_condition" {
        for_each = coalesce(custom_rule.value.match_conditions, [])
        content {
          match_variable = local.match_variable_map[match_condition.value.match_variable]
          operator       = local.operator_map[match_condition.value.operator]
          match_values   = match_condition.value.match_values
          # Selector only for the keyed variables (Cookies, PostArgs,
          # QueryString, RequestHeader) -- ARM rejects it elsewhere, so
          # empty is not sent.
          selector           = match_condition.value.selector != null && match_condition.value.selector != "" ? match_condition.value.selector : null
          negation_condition = coalesce(match_condition.value.negate_condition, false)
          transforms         = match_condition.value.transforms != null ? [for t in match_condition.value.transforms : local.transform_map[t]] : null
        }
      }
    }
  }

  # Microsoft's managed rule sets -- PREMIUM only (spec-enforced;
  # azurerm rejects them on Standard too).
  dynamic "managed_rule" {
    for_each = local.managed_rules
    content {
      type    = managed_rule.value.type
      version = managed_rule.value.version
      action  = local.managed_rule_set_action_map[managed_rule.value.action]

      dynamic "exclusion" {
        for_each = coalesce(managed_rule.value.exclusions, [])
        content {
          match_variable = local.exclusion_match_variable_map[exclusion.value.match_variable]
          operator       = local.selector_operator_map[exclusion.value.operator]
          selector       = exclusion.value.selector
        }
      }

      dynamic "override" {
        for_each = coalesce(managed_rule.value.overrides, [])
        content {
          rule_group_name = override.value.rule_group_name

          dynamic "exclusion" {
            for_each = coalesce(override.value.exclusions, [])
            content {
              match_variable = local.exclusion_match_variable_map[exclusion.value.match_variable]
              operator       = local.selector_operator_map[exclusion.value.operator]
              selector       = exclusion.value.selector
            }
          }

          dynamic "rule" {
            for_each = coalesce(override.value.rules, [])
            content {
              rule_id = rule.value.rule_id
              action  = local.managed_rule_override_action_map[rule.value.action]
              # azurerm's default here is FALSE (listing a rule disables
              # it -- the common tuning gesture); null lets that default
              # through.
              enabled = rule.value.enabled

              dynamic "exclusion" {
                for_each = coalesce(rule.value.exclusions, [])
                content {
                  match_variable = local.exclusion_match_variable_map[exclusion.value.match_variable]
                  operator       = local.selector_operator_map[exclusion.value.operator]
                  selector       = exclusion.value.selector
                }
              }
            }
          }
        }
      }
    }
  }

  # Log scrubbing: the operator/selector pairing contracts (EqualsAny
  # for IP/URI, selector XOR EqualsAny) are spec-enforced; the provider
  # re-validates them at plan time.
  dynamic "log_scrubbing" {
    for_each = var.spec.log_scrubbing != null ? [var.spec.log_scrubbing] : []
    content {
      enabled = log_scrubbing.value.enabled

      dynamic "scrubbing_rule" {
        for_each = log_scrubbing.value.scrubbing_rules
        content {
          enabled        = scrubbing_rule.value.enabled
          match_variable = local.scrubbing_match_variable_map[scrubbing_rule.value.match_variable]
          # Absent means Equals (the provider's default) -- null lets it
          # through.
          operator = scrubbing_rule.value.operator != null ? local.selector_operator_map[scrubbing_rule.value.operator] : null
          selector = scrubbing_rule.value.selector != null && scrubbing_rule.value.selector != "" ? scrubbing_rule.value.selector : null
        }
      }
    }
  }

  tags = local.final_tags
}
