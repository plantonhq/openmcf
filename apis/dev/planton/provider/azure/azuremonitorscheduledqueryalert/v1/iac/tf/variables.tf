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
  description = "Azure Monitor scheduled query alert specification"
  type = object({
    # The Azure region -- must match the queried workspace's region
    region = string

    # The Azure Resource Group name (references resolved by the platform
    # before the module runs)
    resource_group = string

    # The alert rule name (ForceNew)
    alert_name = string

    # The ARM id of the queried resource -- a Log Analytics Workspace or
    # an Application Insights resource (references resolved by the
    # platform). ForceNew.
    scope = string

    # Portal / notification display name
    display_name = optional(string)

    # Notification context (runbook links, on-call notes)
    description = optional(string)

    # Whether the rule evaluates at all
    enabled = optional(bool, true)

    # Severity 0 (critical) - 4 (verbose)
    severity = optional(number, 3)

    # Evaluation cadence and query window, ISO 8601 durations
    # (spec-closed vocabularies)
    evaluation_frequency = optional(string, "PT5M")
    window_duration      = optional(string, "PT5M")

    # Conditions -- each evaluates independently; the rule fires when
    # any holds.
    criteria = list(object({
      query = string
      # Enum value names (spec-closed): COUNT/AVERAGE/MINIMUM/MAXIMUM/TOTAL
      time_aggregation_method = string
      # Enum value names: EQUAL/GREATER_THAN/GREATER_THAN_OR_EQUAL/
      # LESS_THAN/LESS_THAN_OR_EQUAL
      operator              = string
      threshold             = number
      metric_measure_column = optional(string)
      resource_id_column    = optional(string)
      dimensions = optional(list(object({
        name = string
        # INCLUDE / EXCLUDE
        operator = string
        values   = list(string)
      })), [])
      failing_periods = optional(object({
        minimum_failing_periods_to_trigger_alert = number
        number_of_evaluation_periods             = number
      }))
    }))

    # Query lookback override (spec-closed vocabulary)
    query_time_range_override = optional(string)

    # Stateful auto-resolve; mutually exclusive with the mute duration
    # (spec-enforced)
    auto_mitigation_enabled = optional(bool, false)

    # Suppress repeat firings for this duration (spec-closed vocabulary)
    mute_actions_after_alert_duration = optional(string)

    # Verify workspace-alerts storage is configured
    workspace_alerts_storage_enabled = optional(bool, false)

    # Skip query validation at create (custom-log tables that appear
    # only after data flows)
    skip_query_validation = optional(bool, false)

    # Managed identity: type is the spec enum's value name
    # (SYSTEM_ASSIGNED / USER_ASSIGNED); ids carry resolved ARM ids.
    identity = optional(object({
      type                       = string
      user_assigned_identity_ids = optional(list(string), [])
    }))

    # Resource types targeted when the query projects a resource-id column
    target_resource_types = optional(list(string), [])

    # Action groups + payload customization
    action = optional(object({
      action_group_ids  = optional(list(string), [])
      custom_properties = optional(map(string), {})
      email_subject     = optional(string)
    }))

    # User tags, merged over the metadata-derived tags (user wins)
    tags = optional(map(string), {})
  })
}
