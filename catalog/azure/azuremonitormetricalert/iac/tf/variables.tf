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
  description = "Azure Monitor metric alert specification"
  type = object({
    # The Azure Resource Group name (references resolved by the platform
    # before the module runs). Metric alerts are global -- no region.
    resource_group = string

    # The alert rule name (ForceNew)
    alert_name = string

    # The ARM ids whose metrics are evaluated (references resolved by
    # the platform). Multi-scope requires target_resource_type/location.
    scopes = list(string)

    # Notification context (runbook links, on-call notes)
    description = optional(string)

    # Whether the rule evaluates at all
    enabled = optional(bool, true)

    # Stateful auto-resolve behavior (Azure's default)
    auto_mitigate = optional(bool, true)

    # Severity 0 (critical) - 4 (verbose)
    severity = optional(number, 3)

    # Evaluation cadence and rolling window, ISO 8601 durations
    # (spec-closed vocabularies)
    frequency   = optional(string, "PT1M")
    window_size = optional(string, "PT5M")

    # Static threshold conditions (AND-combined). Exactly one condition
    # family is set (spec-enforced).
    static_criteria = optional(list(object({
      metric_namespace = string
      metric_name      = string
      # Enum value names (spec-closed): AVERAGE/COUNT/MINIMUM/MAXIMUM/TOTAL
      aggregation = string
      # Enum value names; static direction subset (spec-enforced)
      operator  = string
      threshold = number
      dimensions = optional(list(object({
        name = string
        # INCLUDE / EXCLUDE / STARTS_WITH
        operator = string
        values   = list(string)
      })), [])
      skip_metric_validation = optional(bool, false)
    })), [])

    # Dynamic (machine-learning) threshold condition -- at most one.
    dynamic_criteria = optional(object({
      metric_namespace = string
      metric_name      = string
      aggregation      = string
      # GREATER_THAN / LESS_THAN / GREATER_OR_LESS_THAN (spec-enforced)
      operator = string
      # LOW / MEDIUM / HIGH
      alert_sensitivity        = string
      evaluation_total_count   = optional(number, 4)
      evaluation_failure_count = optional(number, 4)
      ignore_data_before       = optional(string)
      dimensions = optional(list(object({
        name     = string
        operator = string
        values   = list(string)
      })), [])
      skip_metric_validation = optional(bool, false)
    }))

    # Application Insights web-test availability condition.
    web_test_availability_criteria = optional(object({
      web_test_id           = string
      component_id          = string
      failed_location_count = number
    }))

    # Required for multi-resource / group / subscription scopes
    target_resource_type     = optional(string)
    target_resource_location = optional(string)

    # Action groups fired on trigger/resolve
    actions = optional(list(object({
      action_group_id    = string
      webhook_properties = optional(map(string), {})
    })), [])

    # User tags, merged over the metadata-derived tags (user wins)
    tags = optional(map(string), {})
  })
}
