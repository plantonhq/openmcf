# The scheduled query alert runs KQL against the scoped workspace (or
# Application Insights resource) on a schedule and fires action groups
# when its condition holds -- the alerting half of the logging pipeline.
# The rule is regional and must live in the same region as the resource
# it queries.
#
# The auto-mitigation / mute-duration exclusivity and the failing-periods
# bounds are enforced by the spec before the module runs; the
# metric-measure-column pairing (required for non-Count aggregations,
# forbidden for Count) is Azure's apply-time contract, documented on the
# spec field.
resource "azurerm_monitor_scheduled_query_rules_alert_v2" "main" {
  name                = var.spec.alert_name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region
  scopes              = [var.spec.scope]

  display_name = var.spec.display_name
  description  = var.spec.description
  enabled      = var.spec.enabled
  severity     = var.spec.severity

  evaluation_frequency = var.spec.evaluation_frequency
  window_duration      = var.spec.window_duration

  query_time_range_override         = var.spec.query_time_range_override
  auto_mitigation_enabled           = var.spec.auto_mitigation_enabled
  mute_actions_after_alert_duration = var.spec.mute_actions_after_alert_duration
  workspace_alerts_storage_enabled  = var.spec.workspace_alerts_storage_enabled
  skip_query_validation             = var.spec.skip_query_validation
  target_resource_types             = length(var.spec.target_resource_types) > 0 ? var.spec.target_resource_types : null

  dynamic "criteria" {
    for_each = var.spec.criteria
    content {
      query                   = criteria.value.query
      time_aggregation_method = local.time_aggregation_map[criteria.value.time_aggregation_method]
      operator                = local.operator_map[criteria.value.operator]
      threshold               = criteria.value.threshold
      metric_measure_column   = criteria.value.metric_measure_column
      resource_id_column      = criteria.value.resource_id_column

      dynamic "dimension" {
        for_each = criteria.value.dimensions
        content {
          name     = dimension.value.name
          operator = local.dimension_operator_map[dimension.value.operator]
          values   = dimension.value.values
        }
      }

      # The flap damper: require N of M recent evaluations to breach
      # before firing.
      dynamic "failing_periods" {
        for_each = criteria.value.failing_periods != null ? [criteria.value.failing_periods] : []
        content {
          minimum_failing_periods_to_trigger_alert = failing_periods.value.minimum_failing_periods_to_trigger_alert
          number_of_evaluation_periods             = failing_periods.value.number_of_evaluation_periods
        }
      }
    }
  }

  # The rule's managed identity -- required when the workspace enforces
  # Entra-only query access.
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  dynamic "action" {
    for_each = var.spec.action != null ? [var.spec.action] : []
    content {
      action_groups     = action.value.action_group_ids
      custom_properties = action.value.custom_properties
      email_subject     = action.value.email_subject
    }
  }

  tags = local.final_tags
}
