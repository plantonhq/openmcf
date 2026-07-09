# The metric alert evaluates platform metrics on a rolling window and
# fires the referenced action groups. Exactly one condition family is
# configured (spec-enforced): static thresholds (AND-combined), one
# dynamic machine-learning threshold, or a web-test availability
# condition. Multi-resource / group / subscription scopes additionally
# require target_resource_type + target_resource_location so Azure can
# resolve the metric definition.
resource "azurerm_monitor_metric_alert" "main" {
  name                = var.spec.alert_name
  resource_group_name = var.spec.resource_group
  scopes              = var.spec.scopes

  description   = var.spec.description
  enabled       = var.spec.enabled
  auto_mitigate = var.spec.auto_mitigate
  severity      = var.spec.severity
  frequency     = var.spec.frequency
  window_size   = var.spec.window_size

  target_resource_type     = var.spec.target_resource_type
  target_resource_location = var.spec.target_resource_location

  # Static thresholds -- multiple criteria AND together.
  dynamic "criteria" {
    for_each = var.spec.static_criteria
    content {
      metric_namespace       = criteria.value.metric_namespace
      metric_name            = criteria.value.metric_name
      aggregation            = local.aggregation_map[criteria.value.aggregation]
      operator               = local.operator_map[criteria.value.operator]
      threshold              = criteria.value.threshold
      skip_metric_validation = criteria.value.skip_metric_validation

      dynamic "dimension" {
        for_each = criteria.value.dimensions
        content {
          name     = dimension.value.name
          operator = local.dimension_operator_map[dimension.value.operator]
          values   = dimension.value.values
        }
      }
    }
  }

  # Dynamic machine-learning threshold -- Azure learns the metric's
  # normal band; sensitivity controls how tightly the band hugs it.
  dynamic "dynamic_criteria" {
    for_each = var.spec.dynamic_criteria != null ? [var.spec.dynamic_criteria] : []
    content {
      metric_namespace         = dynamic_criteria.value.metric_namespace
      metric_name              = dynamic_criteria.value.metric_name
      aggregation              = local.aggregation_map[dynamic_criteria.value.aggregation]
      operator                 = local.operator_map[dynamic_criteria.value.operator]
      alert_sensitivity        = local.sensitivity_map[dynamic_criteria.value.alert_sensitivity]
      evaluation_total_count   = dynamic_criteria.value.evaluation_total_count
      evaluation_failure_count = dynamic_criteria.value.evaluation_failure_count
      ignore_data_before       = dynamic_criteria.value.ignore_data_before
      skip_metric_validation   = dynamic_criteria.value.skip_metric_validation

      dynamic "dimension" {
        for_each = dynamic_criteria.value.dimensions
        content {
          name     = dimension.value.name
          operator = local.dimension_operator_map[dimension.value.operator]
          values   = dimension.value.values
        }
      }
    }
  }

  # Web-test availability -- fires when the referenced Application
  # Insights availability test fails from N locations.
  dynamic "application_insights_web_test_location_availability_criteria" {
    for_each = var.spec.web_test_availability_criteria != null ? [var.spec.web_test_availability_criteria] : []
    content {
      web_test_id           = application_insights_web_test_location_availability_criteria.value.web_test_id
      component_id          = application_insights_web_test_location_availability_criteria.value.component_id
      failed_location_count = application_insights_web_test_location_availability_criteria.value.failed_location_count
    }
  }

  dynamic "action" {
    for_each = var.spec.actions
    content {
      action_group_id    = action.value.action_group_id
      webhook_properties = action.value.webhook_properties
    }
  }

  tags = local.final_tags
}
