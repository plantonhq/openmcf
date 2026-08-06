# Workspace-based Application Insights: telemetry is stored in the
# referenced Log Analytics Workspace (classic mode was retired by Azure
# in February 2024, so the binding is required by the spec).
#
# Notable Azure behaviors this module leans on:
#   - workspace_id can be repointed to another workspace but never
#     removed once set (history stays in the old workspace).
#   - The daily data cap and its notification toggle live on a separate
#     billing API; the provider applies them in a follow-up call after
#     create -- transparent here.
#   - Azure auto-creates a noisy "Failure Anomalies" smart-detector rule
#     with every component; the provider disables it by default
#     (deleting it just resurrects it server-side).
resource "azurerm_application_insights" "main" {
  name                = var.spec.application_insights_name
  location            = var.spec.region
  resource_group_name = var.spec.resource_group
  application_type    = local.application_type
  workspace_id        = var.spec.workspace_id

  retention_in_days    = var.spec.retention_in_days
  daily_data_cap_in_gb = var.spec.daily_data_cap_in_gb
  sampling_percentage  = var.spec.sampling_percentage

  # Privacy, auth, and network posture. All default to Azure's own
  # defaults; explicit false is preserved end to end because the spec
  # models them with presence.
  # PARITY-EXCEPTION: this module uses the provider's v5-era POSITIVE
  # toggles while the Pulumi module can only reach the deprecated
  # negative forms (pulumi-azure v6.38 bridges disableIpMasking,
  # localAuthenticationDisabled, dailyDataCapNotificationsDisabled and
  # inverts the same spec booleans). The wire property is identical for
  # each pair -- behavior and outputs match exactly. Re-align when the
  # bridge ships the positive forms.
  daily_data_cap_notifications_enabled = var.spec.daily_data_cap_notifications_enabled
  ip_masking_enabled                   = var.spec.ip_masking_enabled
  local_authentication_enabled         = var.spec.local_authentication_enabled
  internet_ingestion_enabled           = var.spec.internet_ingestion_enabled
  internet_query_enabled               = var.spec.internet_query_enabled

  force_customer_storage_for_profiler = var.spec.force_customer_storage_for_profiler

  tags = local.final_tags
}
