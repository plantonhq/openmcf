# Create the activity log alert. The criteria block's category selects the
# Activity Log slice; the other fields narrow within it (sent only when set,
# the spec's exclusivity CELs guarantee valid combinations). The plural list
# fields carry the single case as a one-element list. Each action notifies
# an action group.
resource "azurerm_monitor_activity_log_alert" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = local.location
  scopes              = var.spec.scopes
  description         = var.spec.description != "" ? var.spec.description : null
  enabled             = var.spec.enabled

  criteria {
    category       = local.category
    operation_name = var.spec.criteria.operation_name != "" ? var.spec.criteria.operation_name : null
    caller         = var.spec.criteria.caller != "" ? var.spec.criteria.caller : null

    levels             = length(local.levels) > 0 ? local.levels : null
    resource_providers = length(var.spec.criteria.resource_providers) > 0 ? var.spec.criteria.resource_providers : null
    resource_types     = length(var.spec.criteria.resource_types) > 0 ? var.spec.criteria.resource_types : null
    resource_groups    = length(var.spec.criteria.resource_groups) > 0 ? var.spec.criteria.resource_groups : null
    resource_ids       = length(var.spec.criteria.resource_ids) > 0 ? var.spec.criteria.resource_ids : null
    statuses           = length(var.spec.criteria.statuses) > 0 ? var.spec.criteria.statuses : null
    sub_statuses       = length(var.spec.criteria.sub_statuses) > 0 ? var.spec.criteria.sub_statuses : null

    recommendation_category = local.recommendation_category
    recommendation_impact   = local.recommendation_impact
    recommendation_type     = var.spec.criteria.recommendation_type != "" ? var.spec.criteria.recommendation_type : null

    dynamic "resource_health" {
      for_each = local.resource_health != null ? [local.resource_health] : []
      content {
        current  = length(resource_health.value.current) > 0 ? resource_health.value.current : null
        previous = length(resource_health.value.previous) > 0 ? resource_health.value.previous : null
        reason   = length(resource_health.value.reason) > 0 ? resource_health.value.reason : null
      }
    }

    dynamic "service_health" {
      for_each = local.service_health != null ? [local.service_health] : []
      content {
        events    = length(service_health.value.events) > 0 ? service_health.value.events : null
        locations = length(service_health.value.locations) > 0 ? service_health.value.locations : null
        services  = length(service_health.value.services) > 0 ? service_health.value.services : null
      }
    }
  }

  dynamic "action" {
    for_each = var.spec.actions
    content {
      action_group_id    = action.value.action_group_id
      webhook_properties = length(action.value.webhook_properties) > 0 ? action.value.webhook_properties : null
    }
  }

  tags = local.final_tags
}
