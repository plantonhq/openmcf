# The workspace is the central Azure Monitor data platform: diagnostic
# settings, Container Insights, Application Insights, and Sentinel all
# store into it, and scheduled query alerts watch what arrives.
#
# Notable Azure behaviors this module leans on:
#   - Switching between PerGB2018 and CapacityReservation updates in
#     place; any other SKU change is ForceNew (the provider's own
#     transition rule).
#   - reservation_capacity_in_gb_per_day is only legal with the
#     CapacityReservation SKU (spec-enforced before the module runs).
#   - daily_quota_gb of -1 means unlimited -- the provider's own default,
#     sent explicitly so both engines carry the same value.
#   - The provider applies data_collection_rule_id and
#     allow_resource_only_permissions via follow-up update calls (ARM
#     rejects them at create) -- transparent here, but it explains why a
#     create can take a few extra seconds when they are set.
resource "azurerm_log_analytics_workspace" "main" {
  name                = var.spec.workspace_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group

  sku                                = local.sku
  reservation_capacity_in_gb_per_day = var.spec.reservation_capacity_in_gb_per_day
  retention_in_days                  = var.spec.retention_in_days
  daily_quota_gb                     = var.spec.daily_quota_gb

  # Security and network posture. All four default to Azure's own
  # defaults (true); explicit false is preserved end to end because the
  # spec models them with presence. The provider expresses internet
  # ingestion/query as an access-type string (Enabled / Disabled /
  # SecuredByPerimeter); the spec's booleans map to the first two, and
  # SecuredByPerimeter (network security perimeter) is intentionally not
  # reachable until the spec models perimeter association itself.
  local_authentication_enabled    = var.spec.local_authentication_enabled
  internet_ingestion_access_type  = var.spec.internet_ingestion_enabled ? "Enabled" : "Disabled"
  internet_query_access_type      = var.spec.internet_query_enabled ? "Enabled" : "Disabled"
  allow_resource_only_permissions = var.spec.allow_resource_only_permissions

  cmk_for_query_forced                    = var.spec.cmk_for_query_forced
  immediate_data_purge_on_30_days_enabled = var.spec.immediate_data_purge_on_30_days_enabled

  # The default DCR is a literal ARM id (no Data Collection Rule kind
  # exists in the catalog); empty means Azure's default handling.
  data_collection_rule_id = (
    var.spec.data_collection_rule_id != null && var.spec.data_collection_rule_id != ""
    ? var.spec.data_collection_rule_id
    : null
  )

  # Managed identity -- used when the workspace itself reads other
  # resources (dedicated-cluster CMK, linked storage).
  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.user_assigned_identity_ids) > 0 ? identity.value.user_assigned_identity_ids : null
    }
  }

  tags = local.final_tags
}
