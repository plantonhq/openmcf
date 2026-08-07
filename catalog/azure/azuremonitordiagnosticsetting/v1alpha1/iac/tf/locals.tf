locals {
  # Destination-type wire values. The diagnostic setting carries no tags
  # (the ARM extension resource does not support them), so no tag locals
  # exist here.
  log_analytics_destination_type_map = {
    "DEDICATED"         = "Dedicated"
    "AZURE_DIAGNOSTICS" = "AzureDiagnostics"
  }
}
