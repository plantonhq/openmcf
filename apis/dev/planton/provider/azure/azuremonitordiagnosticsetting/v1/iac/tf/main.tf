# The diagnostic setting is an extension resource on the target: it
# selects which platform logs and metrics the target emits and routes
# them to the configured destinations. The spec enforces Azure's real
# contracts up front (at least one category, at least one destination,
# category XOR category_group per log entry) -- the API otherwise
# accepts an empty setting and then 404s on read.
#
# Azure's API is eventually consistent here: the provider polls after
# create and delete until the setting is readable/gone, so applies can
# take a few extra seconds.
resource "azurerm_monitor_diagnostic_setting" "main" {
  name               = var.spec.setting_name
  target_resource_id = var.spec.target_resource_id

  # Destinations -- at least one is set (spec-enforced). The eventhub
  # name only rides along with its authorization rule.
  log_analytics_workspace_id = var.spec.log_analytics_workspace_id
  log_analytics_destination_type = (
    var.spec.log_analytics_destination_type != null && var.spec.log_analytics_destination_type != ""
    ? local.log_analytics_destination_type_map[var.spec.log_analytics_destination_type]
    : null
  )
  storage_account_id             = var.spec.storage_account_id
  eventhub_authorization_rule_id = var.spec.eventhub_authorization_rule_id
  eventhub_name                  = var.spec.eventhub_name
  partner_solution_id            = var.spec.partner_solution_id

  # Log selections: exactly one of category or category_group per entry
  # (spec-enforced XOR); every selected category is enabled.
  dynamic "enabled_log" {
    for_each = var.spec.enabled_logs
    content {
      category       = enabled_log.value.category
      category_group = enabled_log.value.category_group
    }
  }

  # Metric selections.
  dynamic "enabled_metric" {
    for_each = var.spec.enabled_metrics
    content {
      category = enabled_metric.value.category
    }
  }
}
