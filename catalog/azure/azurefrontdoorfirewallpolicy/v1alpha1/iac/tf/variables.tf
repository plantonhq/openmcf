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
  description = "Azure Front Door WAF (firewall) policy specification"
  type = object({
    # The resource group the policy is created in. References are
    # resolved to a literal name by the platform before the module runs.
    # The policy is GLOBAL -- no region; ARM fixes its location.
    resource_group = string

    # The policy's name -- unique within the resource group; 1-128
    # characters, begins with a letter, letters and numbers only.
    # ForceNew.
    policy_name = string

    # STANDARD / PREMIUM -- spec enum value names; absent means
    # STANDARD. ForceNew, and Azure refuses a PREMIUM -> STANDARD
    # downgrade outright. Must match the sku of every profile the
    # policy gets associated with.
    sku = optional(string)

    # DETECTION / PREVENTION -- spec enum value names. Required.
    mode = string

    # Whether the policy is enforced at all. Absent means true (Azure's
    # default).
    enabled = optional(bool)

    # Whether request bodies are inspected. Absent means true (Azure's
    # default).
    request_body_check_enabled = optional(bool)

    # Where clients go when a rule's action is REDIRECT. Required by
    # ARM when any rule redirects (apply-time; cross-field).
    redirect_url = optional(string)

    # Block-response customization: one of Azure's allowed status codes
    # (200/403/405/406/429/990-999) and a base64-encoded body.
    custom_block_response_status_code = optional(number)
    custom_block_response_body        = optional(string)

    # PREMIUM only (spec-enforced): the JS-challenge / CAPTCHA solved
    # lifetimes in minutes (5-1440). On Premium Azure always enables
    # both policies; the module pins the documented default of 30 when
    # unset so the deployed value matches the declared contract.
    js_challenge_cookie_expiration_in_minutes = optional(number)
    captcha_cookie_expiration_in_minutes      = optional(number)

    # Your own match / rate-limit rules, evaluated first by ascending
    # priority. Enums arrive as the spec enum's FULL value names and
    # are mapped to ARM's casing in locals.
    custom_rules = optional(list(object({
      name = string

      # Absent means true (Azure's default).
      enabled = optional(bool)

      # Lower runs first; absent means 1 (Azure's default).
      priority = optional(number)

      # MATCH_RULE / RATE_LIMIT_RULE.
      rule_type = string

      # The rate-limit window/threshold pair; absent means 1 minute /
      # 10 requests (Azure's defaults). Ignored by MATCH_RULE rules.
      rate_limit_duration_in_minutes = optional(number)
      rate_limit_threshold           = optional(number)

      # ALLOW / BLOCK / LOG / REDIRECT / JS_CHALLENGE / CAPTCHA (the
      # challenge pair is PREMIUM-only, spec-enforced).
      action = string

      # Conditions AND together; values within one condition OR
      # together.
      match_conditions = optional(list(object({
        # COOKIES / POST_ARGS / QUERY_STRING / REMOTE_ADDR /
        # REQUEST_BODY / REQUEST_HEADER / REQUEST_METHOD / REQUEST_URI /
        # SOCKET_ADDR.
        match_variable = string

        # The key inside keyed variables (headers, cookies, args);
        # ARM rejects it elsewhere, so empty is not sent.
        selector = optional(string)

        # ANY / BEGINS_WITH / CONTAINS / ENDS_WITH / EQUAL / GEO_MATCH /
        # GREATER_THAN / GREATER_THAN_OR_EQUAL / IP_MATCH / LESS_THAN /
        # LESS_THAN_OR_EQUAL / REG_EX.
        operator = string

        # 1-600 values, each up to 256 characters (OR semantics).
        match_values = list(string)

        # Absent means false (Azure's default).
        negate_condition = optional(bool)

        # LOWERCASE / REMOVE_NULLS / TRIM / UPPERCASE / URL_DECODE /
        # URL_ENCODE; up to 5.
        transforms = optional(list(string))
      })), [])
    })), [])

    # Microsoft's curated rule sets -- PREMIUM only (spec-enforced).
    # type/version are Azure's own strings (Microsoft_DefaultRuleSet
    # 1.1/2.0/2.1, Microsoft_BotManagerRuleSet 1.0/1.1, legacy
    # DefaultRuleSet 1.0/preview-0.1).
    managed_rules = optional(list(object({
      type    = string
      version = string

      # RULE_SET_BLOCK / RULE_SET_LOG / RULE_SET_REDIRECT -- the action
      # when the set (or its 2.x anomaly score) trips.
      action = string

      # Set-wide exclusions: request parts every rule in the set skips.
      exclusions = optional(list(object({
        # EXCLUDE_QUERY_STRING_ARG_NAMES /
        # EXCLUDE_REQUEST_BODY_JSON_ARG_NAMES /
        # EXCLUDE_REQUEST_BODY_POST_ARG_NAMES /
        # EXCLUDE_REQUEST_COOKIE_NAMES / EXCLUDE_REQUEST_HEADER_NAMES.
        match_variable = string

        # SELECTOR_CONTAINS / SELECTOR_ENDS_WITH / SELECTOR_EQUALS /
        # SELECTOR_EQUALS_ANY / SELECTOR_STARTS_WITH.
        operator = string

        selector = string
      })), [])

      # Per-group tuning with per-rule overrides and scoped exclusions.
      overrides = optional(list(object({
        rule_group_name = string

        exclusions = optional(list(object({
          match_variable = string
          operator       = string
          selector       = string
        })), [])

        rules = optional(list(object({
          rule_id = string

          # Absent means FALSE here (azurerm's deliberate default --
          # listing a rule is the disable gesture).
          enabled = optional(bool)

          # OVERRIDE_ALLOW / OVERRIDE_ANOMALY_SCORING / OVERRIDE_BLOCK /
          # OVERRIDE_CAPTCHA / OVERRIDE_JS_CHALLENGE / OVERRIDE_LOG /
          # OVERRIDE_REDIRECT (version/type gates spec-enforced).
          action = string

          exclusions = optional(list(object({
            match_variable = string
            operator       = string
            selector       = string
          })), [])
        })), [])
      })), [])
    })), [])

    # Scrub sensitive request data out of the WAF's logs.
    log_scrubbing = optional(object({
      # Absent means true (Azure's default).
      enabled = optional(bool)

      scrubbing_rules = list(object({
        # Absent means true (Azure's default).
        enabled = optional(bool)

        # SCRUB_QUERY_STRING_ARG_NAMES / SCRUB_REQUEST_BODY_JSON_ARG_NAMES /
        # SCRUB_REQUEST_BODY_POST_ARG_NAMES / SCRUB_REQUEST_COOKIE_NAMES /
        # SCRUB_REQUEST_HEADER_NAMES / SCRUB_REQUEST_IP_ADDRESS /
        # SCRUB_REQUEST_URI.
        match_variable = string

        # SELECTOR_EQUALS / SELECTOR_EQUALS_ANY only here; absent means
        # SELECTOR_EQUALS (Azure's default).
        operator = optional(string)

        # Required with SELECTOR_EQUALS; omitted with
        # SELECTOR_EQUALS_ANY (spec-enforced).
        selector = optional(string)
      }))
    }))

    # Free-form tags merged over the Planton-derived resource tags; a
    # user tag with the same key wins.
    tags = optional(map(string), {})
  })
}
