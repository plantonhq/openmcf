# Create the Network Watcher flow log -- the recorder that writes
# traffic metadata for ONE target (virtual network, subnet, or network
# interface; NSG targets are retired for new creates -- spec-validated)
# into a storage account, optionally enriched by Traffic Analytics.
#
# The flow log is a CHILD of the region's Network Watcher, addressed by
# watcher name + the watcher's resource group (locals resolve the
# auto-created singleton when the spec leaves them unset). Creating a
# flow log also writes a lifecycle-management rule on the target
# storage account that OVERWRITES existing lifecycle rules -- point it
# at an account without hand-managed lifecycle policy.
resource "azurerm_network_watcher_flow_log" "main" {
  name                 = var.spec.name
  network_watcher_name = local.network_watcher_name
  resource_group_name  = local.network_watcher_resource_group
  location             = var.spec.region

  # The recorded scope; retargeting updates in place.
  target_resource_id = var.spec.target_resource_id

  # Where flow-log files land.
  storage_account_id = var.spec.storage_account_id

  # The provider requires an explicit value; the platform default is
  # true (a flow log exists to record).
  enabled = coalesce(var.spec.enabled, true)

  # Schema version 1 is the provider default; 2 adds flow state and
  # byte/packet counters.
  version = coalesce(var.spec.version, 1)

  retention_policy {
    enabled = var.spec.retention_policy.enabled
    days    = var.spec.retention_policy.days
  }

  # Traffic Analytics enrichment into a Log Analytics workspace:
  # workspace_id is the workspace GUID (customer id),
  # workspace_resource_id its ARM id.
  dynamic "traffic_analytics" {
    for_each = var.spec.traffic_analytics != null ? [var.spec.traffic_analytics] : []
    content {
      enabled               = coalesce(traffic_analytics.value.enabled, true)
      workspace_id          = traffic_analytics.value.workspace_id
      workspace_region      = traffic_analytics.value.workspace_region
      workspace_resource_id = traffic_analytics.value.workspace_resource_id
      interval_in_minutes   = coalesce(traffic_analytics.value.interval_in_minutes, 60)
    }
  }

  tags = local.final_tags
}
