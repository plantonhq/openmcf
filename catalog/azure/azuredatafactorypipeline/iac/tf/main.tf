# Create the Data Factory pipeline. The activities travel as the raw
# JSON "activities" array (Azure owns that schema -- dozens of
# activity types; the catalog deliberately does not re-model it), and
# the provider normalizes JSON key ordering when diffing, so
# reordering keys never shows as drift. Parameters and variables are
# String-typed on this surface (the provider round-trips only
# string-typed entries); declare other-typed parameters inside the
# activities JSON.
resource "azurerm_data_factory_pipeline" "main" {
  name            = var.spec.name
  data_factory_id = var.spec.data_factory_id

  description     = var.spec.description != "" ? var.spec.description : null
  parameters      = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  variables       = length(var.spec.variables) > 0 ? var.spec.variables : null
  activities_json = var.spec.activities_json != "" ? var.spec.activities_json : null
  annotations     = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  concurrency     = var.spec.concurrency
  folder          = var.spec.folder != "" ? var.spec.folder : null

  # The elapsed-time metric threshold, a Data Factory TimeSpan string
  # (e.g. "0.00:30:00").
  monitor_metrics_after_duration = var.spec.monitor_metrics_after_duration != "" ? var.spec.monitor_metrics_after_duration : null
}
